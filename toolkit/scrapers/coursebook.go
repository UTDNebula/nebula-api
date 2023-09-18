package scrapers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var chromedpCtx context.Context

func initChromeDp() context.CancelFunc {
	fmt.Printf("Initializing chromedp...\n")
	var cancelFnc context.CancelFunc
	//allocCtx, _ := chromedp.NewExecAllocator(context.Background(), []chromedp.ExecAllocatorOption{}...)
	chromedpCtx, cancelFnc = chromedp.NewContext(context.Background())
	fmt.Printf("Initialized chromedp!\n")
	return cancelFnc
}

// This function generates a fresh auth token and returns the new headers
func refreshToken() map[string][]string {
	netID := os.Getenv("LOGIN_NETID")
	if netID == "" {
		panic(errors.New("LOGIN_NETID is missing from .env!"))
	}
	password := os.Getenv("LOGIN_PASSWORD")
	if password == "" {
		panic(errors.New("LOGIN_PASSWORD is missing from .env!"))
	}

	fmt.Printf("Getting new token...\n")
	_, err := chromedp.RunResponse(chromedpCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := network.ClearBrowserCookies().Do(ctx)
			return err
		}),
		chromedp.Navigate(`https://wat.utdallas.edu/login`),
		chromedp.WaitVisible(`form#login-form`),
		chromedp.SendKeys(`input#netid`, netID),
		chromedp.SendKeys(`input#password`, password),
		chromedp.Click(`input#login-button`),
		chromedp.WaitVisible(`body`),
	)
	if err != nil {
		panic(err)
	}

	var cookie []string
	_, err = chromedp.RunResponse(chromedpCtx,
		chromedp.Navigate(`https://coursebook.utdallas.edu/`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if len(cookies) == 0 {
				return errors.New("Failed to get a new token!")
			}
			for i := 0; i < len(cookies); i++ {
				if strings.Contains(cookies[i].Name, "PTGSESSID") {
					cookie = []string{fmt.Sprintf("%s=%s", cookies[i].Name, cookies[i].Value)}
				}
			}
			fmt.Printf("Got new token: %s\n", cookie)
			return err
		}),
	)
	if err != nil {
		panic(err)
	}

	return map[string][]string{
		"Host":            {"coursebook.utdallas.edu"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"},
		"Accept":          {"text/html"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		"Cookie":          cookie,
		"Connection":      {"keep-alive"},
	}
}

func ScrapeCoursebook(term string, startPrefix string, outDir string) {

	// Load env vars
	if err := godotenv.Load(); err != nil {
		panic(errors.New("Error loading .env file"))
	}

	// Start chromedp
	cancel := initChromeDp()
	defer cancel()

	// Find index of starting prefix, if one has been given
	startPrefixIndex := 0
	if startPrefix != "" && startPrefix != coursePrefixes[0] {
		for i, prefix := range coursePrefixes {
			if prefix == startPrefix {
				startPrefixIndex = i
				break
			}
		}
		if startPrefixIndex == 0 {
			panic(errors.New("Failed to find provided course prefix! Remember, the format is cp_<PREFIX>!"))
		}
	}

	// Init http client
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	cli := &http.Client{Transport: tr}

	// Make the output directory for this term
	termDir := fmt.Sprintf("%s/%s", outDir, term)
	if err := os.MkdirAll(termDir, 0777); err != nil {
		panic(err)
	}

	// Keep track of how many total sections we've scraped
	totalSections := 0

	// Scrape all sections for each course prefix
	for prefixIndex, coursePrefix := range coursePrefixes {

		// Skip to startPrefixIndex
		if prefixIndex < startPrefixIndex {
			continue
		}

		// Make a directory in the output for this course prefix
		courseDir := fmt.Sprintf("%s/%s", termDir, coursePrefix)
		if err := os.MkdirAll(courseDir, 0777); err != nil {
			panic(err)
		}
		// Get a fresh token at the start of each new prefix because we can lol
		coursebookHeaders := refreshToken()
		// Give coursebook some time to recognize the new token
		time.Sleep(500 * time.Millisecond)
		// String builder to store accumulated course HTML data for both class levels
		courseBuilder := strings.Builder{}

		fmt.Printf("Finding sections for course prefix %s...\n", coursePrefix)

		// Get courses for term and prefix, split by grad and undergrad to avoid 300 section cap
		for _, clevel := range []string{"clevel_u", "clevel_g"} {
			queryStr := fmt.Sprintf("action=search&s%%5B%%5D=term_%s&s%%5B%%5D=%s&s%%5B%%5D=%s", term, coursePrefix, clevel)
			req, err := http.NewRequest("POST", "https://coursebook.utdallas.edu/clips/clip-cb11-hat.zog", strings.NewReader(queryStr))
			if err != nil {
				panic(err)
			}
			req.Header = coursebookHeaders
			res, err := cli.Do(req)
			if err != nil {
				panic(err)
			}
			if res.StatusCode != 200 {
				panic(errors.New(fmt.Sprintf("ERROR: Section find failed! Status was: %s\nIf the status is 404, you've likely been IP ratelimited!", res.Status)))
			}
			buf := bytes.Buffer{}
			buf.ReadFrom(res.Body)
			courseBuilder.Write(buf.Bytes())
		}
		// Find all section IDs in returned data
		sectionRegexp := regexp.MustCompile(fmt.Sprintf("View details for section (%s[0-9v]{4}\\.[0-9a-z]{3}\\.[0-9]{2}[suf])", coursePrefix[3:]))
		smatches := sectionRegexp.FindAllStringSubmatch(courseBuilder.String(), -1)
		sectionIDs := make([]string, 0, len(smatches))
		for _, matchSet := range smatches {
			sectionIDs = append(sectionIDs, matchSet[1])
		}
		fmt.Printf("Found %d sections for course prefix %s\n", len(sectionIDs), coursePrefix)

		// Get HTML data for all section IDs
		sectionsInCoursePrefix := 0
		for sectionIndex, id := range sectionIDs {
			queryStr := fmt.Sprintf("id=%s&req=0bd73666091d3d1da057c5eeb6ef20a7df3CTp0iTMYFuu9paDeUptMzLYUiW4BIk9i8LIFcBahX2E2b18WWXkUUJ1Y7Xq6j3WZAKPbREfGX7lZY96lI7btfpVS95YAprdJHX9dc5wM=&action=section&div=r-62childcontent", id)
			req, err := http.NewRequest("POST", "https://coursebook.utdallas.edu/clips/clip-cb11-hat.zog", strings.NewReader(queryStr))
			req.Header = coursebookHeaders
			res, err := cli.Do(req)
			if err != nil {
				panic(err)
			}
			if res.StatusCode != 200 {
				panic(errors.New(fmt.Sprintf("ERROR: Section id lookup for id %s failed! Status was: %s\nIf the status is 404, you've likely been IP ratelimited!", id, res.Status)))
			}
			fptr, err := os.Create(fmt.Sprintf("%s/%s.html", courseDir, id))
			if err != nil {
				panic(err)
			}
			buf := bytes.Buffer{}
			buf.ReadFrom(res.Body)
			if _, err := fptr.Write(buf.Bytes()); err != nil {
				panic(err)
			}
			fptr.Close()
			fmt.Printf("Got section: %s\n", id)
			if sectionIndex%60 == 0 && sectionIndex != 0 {
				// Ratelimit? What ratelimit?
				coursebookHeaders = refreshToken()
				// Give coursebook some time to recognize the new token
				time.Sleep(500 * time.Millisecond)
			}
			sectionsInCoursePrefix++
		}
		fmt.Printf("\nFinished scraping course prefix %s. Got %d sections.\n", coursePrefix, sectionsInCoursePrefix)
		totalSections += sectionsInCoursePrefix
	}
	fmt.Printf("\nDone scraping term! Scraped a total of %d sections.", totalSections)
}

var coursePrefixes = []string{
	"cp_acct",
	"cp_acn",
	"cp_acts",
	"cp_aero",
	"cp_ahst",
	"cp_ams",
	"cp_arab",
	"cp_arhm",
	"cp_arts",
	"cp_atcm",
	"cp_aud",
	"cp_ba",
	"cp_bbsu",
	"cp_bcom",
	"cp_biol",
	"cp_bis",
	"cp_blaw",
	"cp_bmen",
	"cp_bps",
	"cp_buan",
	"cp_ce",
	"cp_cgs",
	"cp_chem",
	"cp_chin",
	"cp_cldp",
	"cp_comd",
	"cp_comm",
	"cp_crim",
	"cp_crwt",
	"cp_cs",
	"cp_danc",
	"cp_econ",
	"cp_ecs",
	"cp_ecsc",
	"cp_ed",
	"cp_ee",
	"cp_eebm",
	"cp_eecs",
	"cp_eect",
	"cp_eedg",
	"cp_eegr",
	"cp_eemf",
	"cp_eeop",
	"cp_eepe",
	"cp_eerf",
	"cp_eesc",
	"cp_engr",
	"cp_engy",
	"cp_entp",
	"cp_envr",
	"cp_epcs",
	"cp_epps",
	"cp_film",
	"cp_fin",
	"cp_fren",
	"cp_ftec",
	"cp_geog",
	"cp_geos",
	"cp_germ",
	"cp_gisc",
	"cp_govt",
	"cp_gst",
	"cp_hcs",
	"cp_hdcd",
	"cp_hist",
	"cp_hlth",
	"cp_hmgt",
	"cp_hons",
	"cp_huas",
	"cp_huhi",
	"cp_huma",
	"cp_idea",
	"cp_ims",
	"cp_ipec",
	"cp_isae",
	"cp_isah",
	"cp_isis",
	"cp_isns",
	"cp_itss",
	"cp_japn",
	"cp_kore",
	"cp_lang",
	"cp_lats",
	"cp_lit",
	"cp_mais",
	"cp_mas",
	"cp_math",
	"cp_mech",
	"cp_meco",
	"cp_mils",
	"cp_mis",
	"cp_mkt",
	"cp_msen",
	"cp_mthe",
	"cp_musi",
	"cp_nats",
	"cp_nsc",
	"cp_ob",
	"cp_obhr",
	"cp_opre",
	"cp_pa",
	"cp_phil",
	"cp_phin",
	"cp_phys",
	"cp_ppol",
	"cp_pppe",
	"cp_psci",
	"cp_psy",
	"cp_psyc",
	"cp_real",
	"cp_rels",
	"cp_rhet",
	"cp_rmis",
	"cp_sci",
	"cp_se",
	"cp_smed",
	"cp_soc",
	"cp_span",
	"cp_spau",
	"cp_stat",
	"cp_syse",
	"cp_sysm",
	"cp_te",
	"cp_thea",
	"cp_univ",
	"cp_utd",
	"cp_utsw",
	"cp_viet",
	"cp_vpas",
}
