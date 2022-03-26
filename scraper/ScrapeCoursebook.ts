import { CoursebookScraper } from './CoursebookScraper';
import firefox from 'selenium-webdriver/firefox';

const args = process.argv.slice(2);

// Load Selenium config
const options = new firefox.Options();
const service = new firefox.ServiceBuilder(process.env.SELENIUM_DRIVER);

const coursebook_scraper = new CoursebookScraper(options, service);
coursebook_scraper.Scrape(
    args[0] ? new RegExp(args[0], 'i') : null,
    args[1] ? new RegExp(args[1], 'i') : null,
    args[2] ? new RegExp(args[2], 'i') : null,
    args[3] ? new RegExp(args[3], 'i') : null
).then(() => {
    coursebook_scraper.Kill();
});