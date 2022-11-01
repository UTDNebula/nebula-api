# Scrapers

## Installation
To use the scrapers:
- Download geckodriver for your OS from [here](https://github.com/mozilla/geckodriver/releases)
- Unzip the driver somewhere
- Set the path to the driver in /scrapers/.env as `SELENIUM_DRIVER=`
	- i.e., on Windows, this could be `.\geckodriver.exe`
- Set your NETID and password in /scrapers/.env as `NETID= and PASSWORD=`
- Set your MongoDB login in /tools/.env as `MONGO_USERNAME= and MONGO_PASSWORD=`

## Usage
- Run ScrapeProfiles.ts to obtain base professor profile information
- Run ScrapeCoursebook.ts to obtain course/section information as well as further professor info (i.e. sections taught)
- Organize new data in /data/ into a new directory (usually representing a semester, i.e. 22S for 2022 Spring)
- Once all necessary data is scraped, combine data directories by runing /tools/Combine.ts
- Run /tools/Parse.ts to parse CombinedPsuedoCourses into CombinedCourses (parses course requisites)
- Upload using /tools/DataToMongo.ts