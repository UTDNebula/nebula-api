package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var evaluationCollection *mongo.Collection = configs.GetCollection(configs.DB, "evaluations")

func EvalBySectionID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	sectionId := c.Param("id")

	var eval schema.Evaluation
	var section schema.Section
	var course schema.Course

	defer cancel()

	// Parse object id from id parameter
	objId, err := primitive.ObjectIDFromHex(sectionId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	// First, check if we've already parsed an eval for this section before
	err = evaluationCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&eval)

	// If not, perform on-demand scraping
	if err != nil {
		// If err is anything other than the document not existing, it's likely a database issue; notify the user
		if err != mongo.ErrNoDocuments {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.SectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// Find and parse matching section
		err = sectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&section)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.SectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// Find and parse course associated with section
		objId, err = primitive.ObjectIDFromHex(string(section.Course_reference))
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		err = courseCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&course)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.SectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		eval = ScrapeEval(course, section)
	}

	// Return result
	c.JSON(http.StatusOK, responses.SectionResponse{Status: http.StatusOK, Message: "success", Data: eval})
}

var captchaCompleted bool = false

// Performs on-demand scraping for the eval of a given section
func ScrapeEval(course schema.Course, section schema.Section) schema.Evaluation {

	// Make sure chromedp is initialized
	chromedpCtx, cancel := initChromeDp()
	defer cancel()

	refreshToken(chromedpCtx)

	if !captchaCompleted {
		time.Sleep(30 * time.Second)
		captchaCompleted = true
	}

	// Get a token from coursebook, because we need that for the ues-report endpoint to work properly
	//refreshToken(chromedpCtx)

	sectionID := course.Subject_prefix + course.Course_number + "." + section.Section_number + "." + section.Academic_session.Name

	log.WriteDebug(fmt.Sprintf("Finding eval for %s", sectionID))

	// Get eval info
	evalURL := fmt.Sprintf("https://coursebook.utdallas.edu/ues-report/%s", sectionID)
	// Navigate to eval URL and pull all HTML
	var html string
	_, err := chromedp.RunResponse(chromedpCtx, chromedp.Tasks{
		chromedp.Navigate(evalURL),
		chromedp.QueryAfter("table", func(ctx context.Context, eci runtime.ExecutionContextID, n ...*cdp.Node) error {
			if len(n) > 0 {
				// Create and write eval HTML to file
				chromedp.OuterHTML("html", &html).Do(ctx)
				log.WriteDebug("Eval found!")
				return nil
			} else {
				log.WriteDebug("No eval found!")
				return nil
			}
		}, chromedp.AtLeast(0)),
	})
	if err != nil {
		panic(err)
	}

	// Here is where we'd actually parse the HTML and build an evaluation... unfortunately, coursebook captchas have decided that's not happening anytime soon!

	return schema.Evaluation{}
}

// The 2 functions below are copied from API-Tools to support on-demand scraping

var chromeDpMutex sync.Mutex

// Initializes chromedp using the default executable allocator
func initChromeDp() (chromedpCtx context.Context, cancelFnc context.CancelFunc) {
	chromeDpMutex.Lock()
	log.WriteDebug("Initializing chromedp...")
	allocCtx, cancelFnc := chromedp.NewExecAllocator(context.Background())
	chromedpCtx, _ = chromedp.NewContext(allocCtx)
	log.WriteDebug("Initialized chromedp!")
	chromeDpMutex.Unlock()
	return
}

const tokenRateLimit time.Duration = time.Second * 10

var lastTokenTime time.Time = time.Date(2003, time.March, 21, 7, 47, 0, 0, time.Now().Location())
var cachedCookie map[string][]string

// Generates a fresh auth token and returns the new headers
func refreshToken(chromedpCtx context.Context) map[string][]string {

	// Due to how Gin works, multiple request goroutines may try to refresh their token simultaneously; wrap this area in a mutex lock so as to avoid overlapping refreshes
	chromeDpMutex.Lock()

	// Just return the last cached cookie if we're exceeding the ratelimit
	if time.Since(lastTokenTime) < tokenRateLimit {
		return cachedCookie
	}

	netID, password := configs.GetEnvLogin()

	log.WriteDebug("Getting new token...")
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

	var cookieStrs []string
	_, err = chromedp.RunResponse(chromedpCtx,
		chromedp.Navigate(`https://coursebook.utdallas.edu/`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			cookieStrs = make([]string, len(cookies))
			gotToken := false
			for i, cookie := range cookies {
				cookieStrs[i] = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
				if cookie.Name == "PTGSESSID" {
					log.WriteDebug(fmt.Sprintf("Got new token: PTGSESSID = %s", cookie.Value))
					gotToken = true
				}
			}
			if !gotToken {
				return errors.New("failed to get a new token")
			}
			return err
		}),
	)
	if err != nil {
		panic(err)
	}

	cachedCookie = map[string][]string{
		"Host":            {"coursebook.utdallas.edu"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"},
		"Accept":          {"text/html"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		"Cookie":          cookieStrs,
		"Connection":      {"keep-alive"},
	}

	// Unlock the mutex now that we've cached a cookie
	chromeDpMutex.Unlock()

	return cachedCookie
}
