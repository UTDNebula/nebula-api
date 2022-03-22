import { ProfilesScraper } from './ProfilesScraper';
import { CoursebookScraper } from './CoursebookScraper';
import { ParsingUtils } from './Utils';
import firefox from 'selenium-webdriver/firefox';
import { existsSync, readFileSync, writeFileSync } from 'fs';

const args = process.argv.slice(2);

// Load Selenium config
const options = new firefox.Options();
const service = new firefox.ServiceBuilder(process.env.SELENIUM_DRIVER);

const coursebook_scraper = new CoursebookScraper(options, service);

// Only run the profile scraper if we don't already have professor data
if (!existsSync("./data/Professors.json")) {
    const profiles_scraper = new ProfilesScraper(options, service);
    profiles_scraper.Scrape().then(() => {
        profiles_scraper.Kill();
    });
}

coursebook_scraper.Scrape(/2022/).then(() => {
    coursebook_scraper.Kill();
});

//let Courses = JSON.parse(readFileSync("./data/courses.json", { encoding: 'utf-8' }).replace("][", ","));
//let Sections = JSON.parse(readFileSync("./data/sections.json", { encoding: 'utf-8' }).replace("][", ","));
//console.log(JSON.stringify(ParsingUtils.ParseReq("A grade of at least a C- in either MATH 2415 or in MATH 2419 or equivalent and a grade of at least a C- in MATH 2418 or equivalent.", Courses, Sections), null, '\t'));