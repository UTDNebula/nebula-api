# Scrapers

## Installation
To use the scrapers:
- Download geckodriver for your OS from [here](https://github.com/mozilla/geckodriver/releases)
- Unzip the driver somewhere
- Set the path to the driver in .env as `SELENIUM_DRIVER=`
	- i.e., on Windows, this could be `.\geckodriver.exe`
- Set your NetID in .env as `NETID=`
- Set your NetID password in .env as `PASSWORD=`

## Usage
- Run ScrapeProfiles.ts to obtain base professor profile information
- Run ScrapeCoursebook.ts to obtain course/section information as well as further professor info (i.e. sections taught)
- Run Parse.ts to parse requisites after all relevant sections/courses have been obtained
- All data is output as JSON files in ./data