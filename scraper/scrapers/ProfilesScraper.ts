/////////////////////////////////
//	Scrapes profiles.utdallas.edu for professor profile data. Should run before CoursebookScraper to ensure full professor data is available.
////////////////////////////////

import { writeFileSync } from 'fs';
import { Builder, By, until, WebElement, NoSuchElementError } from 'selenium-webdriver';
import { FirefoxScraper, ParsingUtils } from '../Utils';
import schemas from '../../api/ts/schemas';
import mongoose from 'mongoose';

export class ProfilesScraper extends FirefoxScraper {
    private BASE_URL = 'https://profiles.utdallas.edu/browse?page=';

    // Cache for the scraped professor data
    private professors: Array<schemas.Professor> = [];

    private parseLocation(text: string): schemas.Location {
        let building = '';
        let room = '';

        const tempSplit = text.split(' '); // TODO format `.` (perhaps map buildings to string format)
        if (tempSplit.length == 2) {
            building = tempSplit[0];
            room = tempSplit[1];
        } else {
            const firstDigitIndex = text.search(/\d/);
            building = text.substring(0, firstDigitIndex);
            room = text.substring(firstDigitIndex);
        }

        return {
            building: building,
            room: room,
            map_uri: 'https://locator.utdallas.edu/' + building + '_' + room,
        };
    }

    private parseList(list: Array<string>): [string, schemas.Location] {
        const result: [string, schemas.Location] = ['', { building: '', room: '', map_uri: '' }];

        for (const element of list) {
            if (element.includes('-')) {
                result[0] = element;
            } else {
                result[1] = this.parseLocation(element);
                return result;
            }
        }

        return result;
    }

    private parseName(fullName: string): [string, string] {
        const commaIndex = fullName.indexOf(',');
        if (commaIndex != -1) {
            fullName = fullName.substring(0, commaIndex);
        }
        const names = fullName.split(' ');
        const ultimateName = names[names.length - 1].toLowerCase();
        if (
            names.length > 2 &&
            (ultimateName === 'jr' ||
                ultimateName === 'sr' ||
                ultimateName === 'I' ||
                ultimateName === 'II' ||
                ultimateName === 'III')
        ) {
            names.pop();
        }
        return [names[0], names[names.length - 1]];
    }

    private async scrapeProfessorLinks(): Promise<Array<string>> {
        await this.Driver.get(this.BASE_URL + '1');
        const pageLinks: Array<WebElement> = await this.Driver.findElements(By.className('page-link'));
        const numPages = parseInt(await pageLinks[pageLinks.length - 2].getText());

        const professorLinks: Array<string> = [];
        for (let curPage = 1; curPage <= numPages; curPage++) {
            await this.Driver.get(this.BASE_URL + curPage);

            const linkElements: Array<WebElement> = await this.Driver.findElements(
                By.xpath("//h5[@class='card-title profile-name']//a"),
            );
            for (const element of linkElements) {
                professorLinks.push(await element.getAttribute('href'));
            }
        }
        return professorLinks;
    }

    async Scrape(): Promise<void> {
        const professorLinks: Array<string> = await this.scrapeProfessorLinks();
        //const professorLinks: Array<string> = ['https://profiles.utdallas.edu/herve.abdi'];

        for (const link of professorLinks) {
            await this.Driver.get(link);

            const fullName: string = await (await this.Driver.findElement(By.xpath('//h2'))).getText();
            const [firstName, lastName]: [string, string] = this.parseName(fullName);

            let imageUri: string;
            try {
                imageUri = await (
                    await this.Driver.findElement(By.xpath("//img[@class='profile_photo']"))
                ).getAttribute('src');
            } catch (error: NoSuchElementError) {
                imageUri = await (
                    await this.Driver.findElement(By.xpath("//div[@class='profile-header  fancy_header ']"))
                ).getAttribute('style');
                imageUri = imageUri.substring(23, imageUri.length - 3);
            }

            const titles: Array<string> = [];
            try {
                const titleHeaders = await this.Driver.findElements(By.xpath('//h6'));
                for (const element of titleHeaders) {
                    let tempText = await element.getText();
                    if (!tempText.includes('$')) {
                        titles.push(tempText);
                    }
                }
            } catch (error: NoSuchElementError) {
                continue;
            }

            let email = null;
            try {
                email = await (
                    await this.Driver.findElement(By.xpath("//a[contains(@id,'☄️')]"))
                ).getText();
            } catch (error: NoSuchElementError) {
                continue;
            }

            const tempDiv: WebElement = await this.Driver.findElement(By.xpath('//div[not(@class)]'));
            let texts: Array<string> = (await tempDiv.getText()).split('\n');
            const tempDivChildren: Array<WebElement> = await tempDiv.findElements(By.xpath('.//*'));
            const toRemoveTexts: Array<string> = [];
            for (const child of tempDivChildren) {
                toRemoveTexts.push(await child.getText());
            }
            texts = texts.filter((text) => !toRemoveTexts.includes(text));
            const [phoneNumber, office]: [string, schemas.Location] = this.parseList(texts);

            this.professors.push({
                _id: new mongoose.Types.ObjectId(),
                first_name: firstName,
                last_name: lastName,
                titles: titles,
                email: email,
                phone_number: phoneNumber,
                office: office,
                profile_uri: link,
                image_uri: imageUri,
                office_hours: [],
                sections: []
            });
        }

        writeFileSync(
            './data/Professors.json',
            JSON.stringify(this.GetProfessors(), null, '\t'),
            { flag: 'w' },
        );
        // Clear the buffer
        this.Clear();
    }

    Clear(): void {
        this.professors = [];
    }

    GetProfessors(): Array<schemas.Professor> {
        return this.professors;
    }
}
