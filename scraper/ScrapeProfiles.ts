/////////////////////////////////
//	Runs the ProfilesScraper.
////////////////////////////////

import { ProfilesScraper } from './scrapers/ProfilesScraper';
import firefox from 'selenium-webdriver/firefox';

// Load Selenium config
const options = new firefox.Options();
const service = new firefox.ServiceBuilder(process.env.SELENIUM_DRIVER);

const profiles_scraper = new ProfilesScraper(options, service);
profiles_scraper.Scrape().then(() => {
    profiles_scraper.Kill();
});