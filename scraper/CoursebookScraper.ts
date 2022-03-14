import 'dotenv/config';
import { readFileSync, writeFileSync } from 'fs';
import { Builder, By, until, WebElement, NoSuchElementError } from 'selenium-webdriver';
import firefox from 'selenium-webdriver/firefox';
import schemas from '../api/schemas';
import mongoose from 'mongoose';

type Credentials = {
    NetID: string,
    Password: string
}

export abstract class ParsingUtils {

    // Convert a list of elements to their string representations
    static async GetElementStrings(Elements: WebElement[]): Promise<string[]> {
        let ConversionPromises: Promise<string>[] = Elements.map(async (Element: WebElement) => { return await Element.getText() });
        return await Promise.all(ConversionPromises);
    }

    // Find an element out of a list of elements by matching a preceding element's text
    static FindLabeledElement(ToSearch: WebElement[], ToSearchStrings: string[], Label: string): WebElement | null {
        let LabelIndex = ToSearchStrings.findIndex((value: string) => { return value == Label });
        if (LabelIndex >= 0)
            return ToSearch[LabelIndex + 1];
        else
            return null;
    };

    // Implementation of FindLabeledElement for string data
    static FindLabeledText(ToSearch: WebElement[], ToSearchStrings: string[], Label: string): string | null {
        let LabelIndex = ToSearchStrings.findIndex((value: string) => { return value == Label });
        if (LabelIndex >= 0)
            return ToSearchStrings[LabelIndex + 1];
        else
            return null;
    }

    // Hacky solution for getting inconsistently formatted text between other text, trims surrounding whitespace by default
    static GetTextBetween(Text: string, Back: string, Front: string, TrimWhitespace: Boolean = true): string {
        return TrimWhitespace ? Text.split(Back, 2)[1].split(Front, 2)[0].trim() : Text.split(Back, 2)[1].split(Front, 2)[0];
    }
}

export abstract class FirefoxScraper {

    protected Driver;

    constructor(options, service) {
      this.Driver = new Builder()
          .forBrowser('firefox')
          .setFirefoxService(service)
          .setFirefoxOptions(options)
          .build();
    };

    // End scraper process
    async Kill() {
        await this.Driver.quit();
    };

    abstract Scrape(): Promise<void>;
    abstract Clear(): void;
};

class CoursebookScraper extends FirefoxScraper {

    // Enum of IDs for the combo boxes
    static DropdownIDs = {
        "TERM": "combobox_term",
        "PREFIX": "combobox_cp"
    };

    // Grab professor data obtained from profile scraper
    private ScrapedProfessors: schemas.Professor[] = JSON.parse(readFileSync("./data/Professors.json", { encoding: "utf-8" }));

    // Caches for the scraped course/section data
    private Courses: Map<string, schemas.Course> = new Map<string, schemas.Course>();
    private Sections: Map<string, schemas.Section> = new Map<string, schemas.Section>();

    // Find the buttons corresponding to a dropdown element
    private async FindDropdownButtons(DropdownID: string, StartIndex: number | RegExp = 0, Filter: RegExp = null): Promise<WebElement[]> {
        // Find the correct dropdown box
        await this.Driver.wait(until.elementLocated(By.id(DropdownID)));
        let Dropdown: WebElement = await this.Driver.findElement(By.id(DropdownID));
        // Find all of the buttons in the dropdown
        let Buttons: WebElement[] = await Dropdown.findElements(By.css("option"));
        // Get the text of all of the buttons in the dropdown
        let ButtonsText: string[] = await ParsingUtils.GetElementStrings(Buttons);
        // Filter out the divider buttons (the ones containing "-----") and blanks
        for (let i: number = 0; i < Buttons.length; ++i) {
            // Ignore dropdown button if it's blank or empty
            if (ButtonsText[i] == "" || ButtonsText[i].match(/---+/g))
                Buttons.splice(i, 1);
            // Ignore dropdown button if it doesn't match a provided filter
            if (Filter && !ButtonsText[i].match(Filter))
                Buttons.splice(i, 1);
        };
        // Refresh ButtonsText
        ButtonsText = await ParsingUtils.GetElementStrings(Buttons);
        // If StartIndex is a RegExp object, convert it to an integer index by finding the first button matching the pattern
        if (StartIndex instanceof RegExp) {
            let Pattern: RegExp = StartIndex;
            StartIndex = ButtonsText.findIndex((ButtonText: string) => { return ButtonText.match(Pattern) })
        }
        // Splice out the skipped indexes
        Buttons.splice(0, StartIndex);
        return Buttons;
    };

    // Searches for sections, returns the section elements once they have loaded
    private async FindSections(): Promise<WebElement[]> {

        // Check if a reCaptcha has been triggered
        let CaptchaIFrame: WebElement = null;
        try {
            let CaptchaBox = await this.Driver.findElement(By.id("recaptcha_v2_here"));
            CaptchaIFrame = await CaptchaBox.findElement(By.css("iframe"));
        } catch (error: NoSuchElementError) {
            // Continue without action
        }

        // If captcha box found, wait for the user to solve it before continuing
        if (CaptchaIFrame) {
            // Switch active frame to the captcha's iframe
            this.Driver.switchTo().frame(CaptchaIFrame);
            let CheckedBox: WebElement = null;
            while (!CheckedBox) {
                try {
                    CheckedBox = await this.Driver.findElement(By.css("span.recaptcha-checkbox-checked"));
                }
                catch (error: NoSuchElementError) {
                    await this.Driver.sleep(1000);
                }
            }
            // Switch active frame back to the topmost frame
            this.Driver.switchTo().defaultContent();
        }
        // Find the class search button, then click it
        let SearchButton: WebElement =await this.Driver.findElement(By.linkText("Search Classes"));
        await SearchButton.click();
        // Wait for loading spinner to go away
        let Selector = By.css("svg.uil-ring-alt");
        await this.Driver.wait(until.stalenessOf(this.Driver.findElement(Selector)));
        // Try to get the section list (may not exist)
        Selector = By.css("div.section-list");
        let SectionList;
        try {
             SectionList = await this.Driver.findElement(Selector);
        }
        // Return empty array if no section list found
        catch (error: NoSuchElementError) {
            return [];
        };
        // Get all of the detail box buttons
        let DetailButtons: WebElement[] = await SectionList.findElements(By.css("div[data-action=info]"));
        // Click every other detail button with a 6 second delay to somewhat diminish Cousebook's stupid ratelimits
        for (let i: number = 0; i < DetailButtons.length; i += 2) {
            await DetailButtons[i].click();
            await this.Driver.sleep(6000);
        }
        // Return the individual section elements
        return await SectionList.findElements(By.className("expandedrow"));
    };

    // NOTE: Removed indefinitely
    //async ParseTextbooks(SectionData: schemas.Section, Section: WebElement) {
    //    // Grab the textbooks tab
    //    let TextbooksTab: WebElement = await Section.findElement(By.css("#tab-textbooks"));
    //    // Click on textbooks tab
    //    await TextbooksTab.click();
    //    // Find the textbook info area, wait for page to load more if we can't find it yet
    //    let TextbookNote: WebElement = null;
    //    while (!TextbookNote) {
    //        try {
    //            TextbookNote = await Section.findElement(By.css("div.textbook-note"));
    //        }
    //        catch (error: NoSuchElementError) {
    //            await this.Driver.sleep(500);
    //        }
    //    }
    //    // Grab the textbook table, may not exist
    //    try {
    //        let TextbookTable: WebElement = await Section.findElement(By.css("div.textbook-table"));
    //        // Get the individual textbook data blocks
    //        let TextbookBlocks: WebElement[] = await TextbookTable.findElements(By.css("td.textbook"));
    //        // Iterate through the textbook blocks
    //        for (let BookBlock of TextbookBlocks) {
    //            let Title: string = await(await BookBlock.findElement(By.css("div.booktitle"))).getText();
    //            let ImageURI: string = await(await BookBlock.findElement(By.css("img"))).getAttribute("src");
    //            // Author is used for more than just the author of the book, so we'll get a list
    //            let AuthorElements: WebElement[] = await BookBlock.findElements(By.css("div.author"));
    //            // The actual author info will be in the first author element
    //            let Author: string = await(await AuthorElements[0].findElement(By.css("a"))).getText();
    //            // Price data is stored in the third author element
    //            let Price: string = await(await AuthorElements[2].findElement(By.css("b"))).getText();
    //            // The version of the book is stored in the 4th author element
    //            let Version: string = await AuthorElements[3].getText();
    //            let ISBNElement: WebElement = await BookBlock.findElement(By.css("div.isbn"));
    //            let ISBN: string = await(await ISBNElement.findElement(By.css("a"))).getText();

    //            SectionData.TextBooks.push({
    //                Title: Title,
    //                Author: Author,
    //                ImageURI: ImageURI,
    //                Price: Price,
    //                Version: Version,
    //                ISBN: ISBN
    //            });
    //        }
    //    } catch (error: NoSuchElementError) { };
    //}

    async ParseSyllabus(SectionData: schemas.Section, TableData: WebElement[], TableDataStrings: string[]) {
        let SyllabusBlock = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Syllabus:");
        SectionData.syllabus_uri = await (SyllabusBlock.findElement(By.css("a")).then(
            // Get href of syllabus link on successful find
            async (SyllabusLink: WebElement) => {
                return await SyllabusLink.getAttribute("href");
            },
            // Return null on fail
            () => { return null; }
        ));
    }

    async ParseMeetings(SectionData: schemas.Section, SectionTable: WebElement) {
        let MeetingBlocks: WebElement[] = await SectionTable.findElements(By.css("p.courseinfo__meeting-time"));
        // Iterate over potentially multiple meetings
        for (let MeetingBlock of MeetingBlocks) {

            let MeetingData: schemas.Meeting = {
                start_date: null,
                end_date: null,
                start_time: null,
                end_time: null,
                location: null,
                meeting_days: null,
                modality: null
            }

            let MeetingText: string = await MeetingBlock.getText();
            let MeetingChunks: string[] = MeetingText.split('\n');
            let Dates: string[] = MeetingChunks[0].split('-');
            let Times: string[] = MeetingChunks[2].split('-');
            MeetingData.start_date = Dates[0] ?? null;
            MeetingData.end_date = Dates[1] ?? null;
            MeetingData.start_time = Times[0] ?? null;
            MeetingData.end_time = Times[1] ?? null;
            MeetingData.meeting_days = MeetingChunks[1].split(', ');
            // Nullify meeting_days if incorrectly formatted
            if (MeetingData.meeting_days[0] == "")
                MeetingData.meeting_days = null;
            let LocationElement: WebElement = MeetingBlock.findElement(By.css("a"));
            let LocationChunks: string[] = (await LocationElement.getText()).split(' ');
            // Nullify location data if incorrectly formatted
            if (LocationChunks.length != 2)
                LocationChunks = [null, null];
            MeetingData.location = {
                building: LocationChunks[0],
                room: LocationChunks[1],
                map_uri: await LocationElement.getAttribute("href")
            };
            // Push meeting information
            SectionData.meetings.push(MeetingData);
        };
    }

    async ParseInstructors(SectionData: schemas.Section, TableData: WebElement[], TableDataStrings: string[]) {
        let InstructorInstance: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Instructor(s):");
        try {
            InstructorInstance = await InstructorInstance.findElement(By.css("div"));
        } catch (error: NoSuchElementError) { return };
        // Iterate by updating the InstructorElement with recursive searches until search fails
        while (true) {
            try {
                let NestedDivs: WebElement[] = await InstructorInstance.findElements(By.css("div"));
                let InstructorElement: WebElement = NestedDivs[0];
                // Unfortunately, string splitting is pretty much the best way to do this part
                let InstructorText: string[] = (await InstructorElement.getText()).split(" ・ ");
                // Derive some instructor data from split string
                let Names: string[] = InstructorText[0].split(' ');
                //let Role: string = InstructorText[1];
                let Email: string = InstructorText[2];
                // Find the matching professor in the professor data
                let Professor: schemas.Professor = this.ScrapedProfessors.find((professor: schemas.Professor) => {
                    return (
                        professor.first_name == Names[0]
                        &&
                        professor.last_name == Names[1]
                    )
                });
                // If we found the matching professor, add its id to this section's professor list
                if (Professor) {
                    // Set professor's email since it's not currently found properly by the professor scraper
                    Professor.email = Email;
                    // Add a reference to this section
                    Professor.sections.push(SectionData._id);
                    SectionData.professors.push(Professor._id);
                }
                InstructorInstance = NestedDivs[1];
            } catch { break }
        }
    }

    async ParseAssistants(SectionData: schemas.Section, TableData: WebElement[], TableDataStrings: string[]) {
        let AssistantInstance: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "TA/RA(s):");
        try {
            AssistantInstance = await AssistantInstance.findElement(By.css("div"));
        } catch (error: NoSuchElementError) { return };
        // Iterate by updating the AssistantElement with recursive searches until search fails
        while (true) {
            try {
                let NestedDivs: WebElement[] = await AssistantInstance.findElements(By.css("div"));
                let AssistantElement: WebElement = NestedDivs[0];
                // Unfortunately, string splitting is pretty much the best way to do this part
                let AssistantText: string[] = (await AssistantElement.getText()).split(" ・ ");
                let Names: string[] = AssistantText[0].split(' ');
                // Push assistant data
                SectionData.teaching_assistants.push({
                    first_name: Names[0],
                    last_name: Names[1],
                    role: AssistantText[1],
                    email: AssistantText[2]
                });
                AssistantInstance = NestedDivs[1];
            } catch { break }
        }
    }

    //ParseRequisiteString(ReqString: string, CourseData: schemas.Course, SectionData: schemas.Section): void {
    //    let ChunkedString: string[] = ReqString.split('.');
    //}

    async ParseRequisites(CourseData: schemas.Course, SectionData: schemas.Section, TableData: WebElement[], TableDataStrings: string[]) {
        //let RequisitesBlock: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Enrollment Reqs:");
        //// Class requisites may not always be available (or not in the requisites block)
        //if (RequisitesBlock != null) {
        //    let RequisiteElements = await RequisitesBlock.findElements(By.css("li"));
        //    for (let Element of RequisiteElements) {
        //        //SectionData.Requisites.push(await Element.getText());
        //        let ReqString: string = await Element.getText();
        //        this.ParseRequisiteString(ReqString, CourseData, SectionData);
        //    };
        //};
        // Add the section's consent requirement
        let ConsentReq: schemas.ConsentRequirement = new schemas.ConsentRequirement();
        ConsentReq.granter = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Add Consent:");
        SectionData.section_corequisites.options.push(ConsentReq);
    }

    async ParseAttributes(SectionData: schemas.Section, TableData: WebElement[], TableDataStrings: string[]) {
        let AttributesBlock: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Class Attributes:");
        // Class attributes may not always be available
        if (AttributesBlock != null) {
            let AttributeElements = await AttributesBlock.findElements(By.css("li"));
            for (let Element of AttributeElements)
                SectionData.attributes["raw_attributes"].push(await Element.getText());
        };
    }

    async ParseCourse(CourseNum: string, TableData: WebElement[], TableDataStrings: string[]): Promise<schemas.Course> {
        let CourseData: schemas.Course = {
            _id: new mongoose.Types.ObjectId(parseInt(CourseNum)),
            course_number: null,
            subject_prefix: null,
            title: null,
            description: null,
            school: null,
            credit_hours: null,
            class_level: null,
            activity_type: null,
            grading: null,
            internal_course_number: null,
            prerequisites: new schemas.CollectionRequirement(),
            corequisites: new schemas.CollectionRequirement(),
            sections: [],
            lecture_contact_hours: null,
            laboratory_contact_hours: null,
            offering_frequency: null,
            attributes: {}
        }
        CourseData.title = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Course Title:");
        CourseData.description = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Description:");
        // Split the section's text to obtain the subject_prefix and course_number
        let SplitSectionText: string[] = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Class Section:").split('.');
        let TextMatches: RegExpMatchArray = SplitSectionText[0].match(/([A-z]+)([0-9]V?[0-9]+)/);
        CourseData.subject_prefix = TextMatches[1];
        CourseData.course_number = TextMatches[2];
        CourseData.school = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "College:");

        let Credits: string = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Semester Credit Hours:");
        // TODO: Handle variable-credit courses that may list a non-numerical # of credits?
        CourseData.credit_hours = Credits;

        CourseData.class_level = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Class Level:");
        CourseData.activity_type = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Activity Type:");
        CourseData.grading = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Grading:");
        CourseData.internal_course_number = CourseNum.trim();

        // Split the course's description into words
        let SplitDescription: string[] = CourseData.description.split(' ');
        // Get just the last two elements of SplitDescription for contact hours and frequency
        SplitDescription = SplitDescription.slice(-2);
        // Grab the contact hours via regex match
        let ContactHours: RegExpMatchArray = SplitDescription[0].match(/([0-9])-([0-9])/);
        // Handle case where offering_freqency is missing
        if (!ContactHours || ContactHours.length < 3)
            ContactHours = SplitDescription[1].match(/([0-9])-([0-9])/);
        else
            CourseData.offering_frequency = SplitDescription[1];
        if (ContactHours) {
            CourseData.lecture_contact_hours = ContactHours[1];
            CourseData.laboratory_contact_hours = ContactHours[2];
        };
        // Store course data in the buffer
        this.Courses[CourseNum] = CourseData;
        // Return the collected course data
        return CourseData;
    }

    // Parse the "section"'s data into course and section data
    async ParseSection(Section: WebElement) {

        // Scroll the section element into view (seems to help page loading)
        await this.Driver.executeScript("arguments[0].scrollIntoView();", Section);

        // Grab the section's full table, wait for page to load more if we can't find it yet
        let SectionTable: WebElement = null;
        while (!SectionTable) {
            try {
                SectionTable = await Section.findElement(By.css("table.courseinfo__overviewtable"));
            }
            catch (error: NoSuchElementError) {
                await this.Driver.sleep(500);
            }
        }

        // Get all of the useful table data elements
        let TableData: WebElement[] = await SectionTable.findElements(By.css("th, td"));
        // Get string versions of all of the table data elements (for FindLabeledText and FindLabeledElement)
        let TableDataStrings: string[] = await ParsingUtils.GetElementStrings(TableData);

        // Find, split, and parse the class/course numbers
        let Nums: string = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Class/Course Number:");
        let SplitNums: string[] = Nums.split('/');
        let ClassNum: string = SplitNums[0];
        let CourseNum: string = SplitNums[1];

        // Init the section data objects
        let SectionData: schemas.Section = {
            _id: new mongoose.Types.ObjectId(parseInt(ClassNum)),
            section_number: null,
            course_reference: null,
            section_corequisites: new schemas.CollectionRequirement(),
            academic_session: null,
            professors: [],
            teaching_assistants: [],
            internal_class_number: null,
            instruction_mode: null,
            meetings: [],
            core_flags: [],
            syllabus_uri: null,
            grade_distribution: [],
            attributes: {
                raw_attributes: []
            }
        }

        let CourseData: schemas.Course = this.Courses[CourseNum];

        // Find and set course data, if not already found
        if (!CourseData)
            CourseData = await this.ParseCourse(CourseNum, TableData, TableDataStrings);

        // Find and set section data
        SectionData.internal_class_number = ClassNum.trim();
        // Split the section's text to obtain the section_number
        let SplitSectionText: string[] = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Class Section:").split('.');
        SectionData.section_number = SplitSectionText[1];
        // Reference the course associated with this section
        SectionData.course_reference = CourseData._id;
        // Parse the section's academic session
        let TermText: string = await (await SectionTable.findElement(By.css("p.courseinfo__sectionterm"))).getText();
        SectionData.academic_session = {
            name: ParsingUtils.GetTextBetween(TermText, "Term: ", "\n"),
            start_date: ParsingUtils.GetTextBetween(TermText, "Starts: ", "\n"),
            end_date: ParsingUtils.GetTextBetween(TermText, "Ends: ", "\n")
        };
        // Get the section's instructors
        await this.ParseInstructors(SectionData, TableData, TableDataStrings);
        // Get the section's TAs/RAs
        await this.ParseAssistants(SectionData, TableData, TableDataStrings);
        // Parse the section's instruction mode
        SectionData.instruction_mode = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Instruction Mode:");
        // Get section's meeting times/dates and location data
        await this.ParseMeetings(SectionData, SectionTable);
        // Parse the section's core flags (may not exist)
        SectionData.core_flags = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Core:")?.match(/[0-9]{3}/g) ?? [];
        // Get the URI to the section's syllabus
        await this.ParseSyllabus(SectionData, TableData, TableDataStrings);
        // Get the section's attributes
        await this.ParseAttributes(SectionData, TableData, TableDataStrings);

        await this.ParseRequisites(CourseData, SectionData, TableData, TableDataStrings);

        // Get the section's textbooks (do this last because switching to the textbook tab makes all previous elements stale)
        //await this.ParseTextbooks(SectionData, Section);

        // Add collected section data to section data cache
        this.Sections[ClassNum] = SectionData;

        // Add this section's ID to the course's sections list
        this.Courses[CourseNum].sections.push(SectionData._id);

        // Print the collected course data
        console.log(JSON.stringify(CourseData, null, '\t'));
        // Print the collected section data
        console.log(JSON.stringify(SectionData, null, '\t'));
        console.log('\n');
    };

    async Login(credentials: Credentials): Promise<void> {
        // Navigate to login page
        await this.Driver.get("https://coursebook.utdallas.edu/login/coursebook");
        // Find netID input box
        let NetIDBox: WebElement = await this.Driver.findElement(By.id("netid"));
        // Enter the netID
        await NetIDBox.sendKeys(credentials.NetID);
        // Find the password input box
        let PasswordBox: WebElement = await this.Driver.findElement(By.id("password"));
        // Enter the password
        await PasswordBox.sendKeys(credentials.Password);
        // Submit and wait for page load
        await PasswordBox.submit();
        await this.Driver.wait(until.elementLocated(By.css("div.search-panel-form-div")));
    }

    // Scrape everything
    async Scrape(TermIndex: number | RegExp = 0, TermFilter: RegExp = null, PrefixIndex: number | RegExp = 0, PrefixFilter: RegExp = null): Promise<void> {
        // Log in with COURSEBOOK_AUTH credentials
        await this.Login({
            NetID: process.env.NETID,
            Password: process.env.Password
        });
        // Find the term buttons
        let TermButtons = await this.FindDropdownButtons(CoursebookScraper.DropdownIDs.TERM, TermIndex, TermFilter);
        // Find the prefix buttons
        let PrefixButtons = await this.FindDropdownButtons(CoursebookScraper.DropdownIDs.PREFIX, PrefixIndex, PrefixFilter);
        // Iterate over sections from every desired class prefix for every desired term
        for (let TermButton of TermButtons) {
            // Click the desired term dropdown button
            await TermButton.click();
            // Iterate over desired prefix buttons
            for (let PrefixButton of PrefixButtons) {
                // Click the desired prefix dropdown button
                await PrefixButton.click();
                // Search for sections and parse them
                let SectionList: WebElement[] = await this.FindSections();
                for (let Section of SectionList) {
                    await this.ParseSection(Section);
                }
                // Write section and course data to data output after all sections under the given prefix are parsed
                writeFileSync("./data/Sections.json", JSON.stringify(this.GetSections(), null, '\t'), { flag: 'a' });
                writeFileSync("./data/Courses.json", JSON.stringify(this.GetCourses(), null, '\t'), { flag: 'a' });
                writeFileSync("./data/Professors.json", JSON.stringify(this.GetProfs(), null, '\t'), { flag: 'w' });
                this.Clear();
            };
        };
    };

    Clear(): void {
        this.Courses.clear();
        this.Sections.clear();
    }

    GetCourses(): schemas.Course[] { return Object.values(this.Courses) };
    GetSections(): schemas.Section[] { return Object.values(this.Sections) };
    GetProfs(): schemas.Professor[] { return this.ScrapedProfessors };
};

// Load Selenium config

const options = new firefox.Options();
const service = new firefox.ServiceBuilder(process.env.SELENIUM_DRIVER);
let CBScraper = new CoursebookScraper(options, service);

CBScraper.Scrape(/2022 Spring/g).then(() => {
    CBScraper.Kill();
});
