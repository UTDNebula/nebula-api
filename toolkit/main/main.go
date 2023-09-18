package main

import (
	"errors"
	"flag"
	"github.com/UTDNebula/nebula-api/toolkit/parser"
	"github.com/UTDNebula/nebula-api/toolkit/scrapers"
	"github.com/UTDNebula/nebula-api/toolkit/uploader"
)

func main() {

	// I/O Flags
	inDir := flag.String("i", "./data", "The directory to read data from.")
	outDir := flag.String("o", "./data", "The directory to write resulting data to.")

	// Flags for all scraping
	scrape := flag.Bool("scrape", false, "Puts the tool into scraping mode. Use flags -coursebook and -profiles for specifying what to scrape.")

	// Flags for coursebook scraping
	scrapeCoursebook := flag.Bool("coursebook", false, "Alongside -scrape, signifies that coursebook should be scraped. Use with -term and -startprefix.")
	term := flag.String("term", "", "For coursebook scraping, the term to scrape, i.e. 23S")
	startPrefix := flag.String("startprefix", "", "For coursebook scraping, the course prefix to start scraping from, i.e. cp_span")

	// Flag for profile scraping
	scrapeProfiles := flag.Bool("profiles", false, "Alongside -scrape, signifies that professor profiles should be scraped.")

	// Flags for parsing
	parse := flag.Bool("parse", false, "Puts the tool into parsing mode. Use the -i flag to specify the input directory for scraped data.")
	csvDir := flag.String("csv", "", "The path to the directory of CSV files containing grade data for the parser to use. No grade distributions will be included if this flag is exluded.")
	skipValidation := flag.Bool("skipv", false, "Signifies that the post-parsing validation should be skipped. Be careful with this!")

	// Flags for uploading data
	upload := flag.Bool("upload", false, "Puts the tool into upload mode. Used alongside the -i and -replace flags.")
	replace := flag.Bool("replace", false, "Specifies that data should be uploaded to the database by replacing any existing data, used alongside -upload.")

	flag.Parse()

	switch {
	case *scrape:
		switch {
		case *scrapeProfiles:
			scrapers.ScrapeProfiles(*outDir)
		case *scrapeCoursebook:
			if *term == "" {
				panic(errors.New("No term specified for coursebook scraping! Use -term to specify."))
			}
			scrapers.ScrapeCoursebook(*term, *startPrefix, *outDir)
		default:
			panic(errors.New("One of the -coursebook or -profiles flags must be set for scraping!"))
		}
	case *parse:
		parser.Parse(*inDir, *outDir, *csvDir, *skipValidation)
	case *upload:
		uploader.Upload(*inDir, *replace)
	default:
		flag.PrintDefaults()
		return
	}
}
