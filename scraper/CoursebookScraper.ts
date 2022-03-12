import 'dotenv/config';
import 'fs';
import { writeFileSync } from 'node:fs';
import { Builder, By, until, WebElement, NoSuchElementError } from 'selenium-webdriver';
import firefox from 'selenium-webdriver/firefox';

type ClassLevel = "Undergraduate" | "Graduate";

type Location = {
    Room: string,
    MapURI: string
}

type Exam = {
    Date: string,
    Time: string,
    Location: Location
};

type Person = {
    Name: string,
    Role: string,
    EMail: string
};

type Textbook = {
    Title: string,
    ImageURI: string,
    Author: string,
    Price: string,
    Version: string,
    ISBN: string
};

class CourseData {
    Title: string = "";
    College: string = "";
};

class SectionData {
    Section: string = "";
    Term: string = "";
    Level: ClassLevel = "Undergraduate";
    Credits: number = 0;
    Grading: string = "";
    Consent: string = "";
    Method: string = "";
    ActivityType: string = "";
    SessionType: string = "";
    EnrollmentStatus: string = "";
    AvailableSeats: number = 0;
    TotalEnrolled: string = "";
    Waitlisted: number = 0;
    Description: string = "";
    Attributes: string[] = [];
    Requisites: string[] = [];
    CoreInfo: string = "";
    Instructors: Person[] = [];
    Assistants: Person[] = [];
    StartDate: string = "";
    EndDate: string = "";
    MeetingDays: string = "";
    Times: string = "";
    Location: Location = {Room: "", MapURI: ""};
    Exams: Exam[] = [];
    SyllabusURI: string = "";
    TextBooks: Textbook[] = []
};

type Courses = {
    [CourseNum: number]: CourseData;
}

type Sections = {
    [ClassNum: number]: SectionData;
}

type Credentials = {
    NetID: string,
    Password: string
}

abstract class ParsingUtils {

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

    // Caches for the scraped course/section data
    private Courses: Courses = {};
    private Sections: Sections = {};

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
        // Find the class search button, then click it
        let SearchButton: WebElement = this.Driver.findElement(By.linkText("Search Classes"));
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
        // Click every other detail button in rapid succession to bypass ratelimit and expose data
        for (let i: number = 0; i < DetailButtons.length; i += 2)
            await DetailButtons[i].click();
        // Return the individual section elements
        return await SectionList.findElements(By.className("expandedrow"));
    };

    async ParseTextbooks(sectionData: SectionData, Section: WebElement) {
        // Grab the textbooks tab
        let TextbooksTab: WebElement = await Section.findElement(By.css("#tab-textbooks"));
        // Click on textbooks tab
        await TextbooksTab.click();
        // Find the textbook info area, wait for page to load more if we can't find it yet
        let TextbookNote: WebElement = null;
        while (!TextbookNote) {
            try {
                TextbookNote = await Section.findElement(By.css("div.textbook-note"));
            }
            catch (error: NoSuchElementError) {
                await this.Driver.sleep(500);
            }
        }
        // Grab the textbook table, may not exist
        try {
            let TextbookTable: WebElement = await Section.findElement(By.css("div.textbook-table"));
            // Get the individual textbook data blocks
            let TextbookBlocks: WebElement[] = await TextbookTable.findElements(By.css("td.textbook"));
            // Iterate through the textbook blocks
            for (let BookBlock of TextbookBlocks) {
                let Title: string = await(await BookBlock.findElement(By.css("div.booktitle"))).getText();
                let ImageURI: string = await(await BookBlock.findElement(By.css("img"))).getAttribute("src");
                // Author is used for more than just the author of the book, so we'll get a list
                let AuthorElements: WebElement[] = await BookBlock.findElements(By.css("div.author"));
                // The actual author info will be in the first author element
                let Author: string = await(await AuthorElements[0].findElement(By.css("a"))).getText();
                // Price data is stored in the third author element
                let Price: string = await(await AuthorElements[2].findElement(By.css("b"))).getText();
                // The version of the book is stored in the 4th author element
                let Version: string = await AuthorElements[3].getText();
                let ISBNElement: WebElement = await BookBlock.findElement(By.css("div.isbn"));
                let ISBN: string = await(await ISBNElement.findElement(By.css("a"))).getText();

                sectionData.TextBooks.push({
                    Title: Title,
                    Author: Author,
                    ImageURI: ImageURI,
                    Price: Price,
                    Version: Version,
                    ISBN: ISBN
                });
            }
        } catch (error: NoSuchElementError) { };
    }

    async ParseSyllabus(sectionData: SectionData, TableData: WebElement[], TableDataStrings: string[]) {
        let SyllabusBlock = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Syllabus:");
        sectionData.SyllabusURI = await (SyllabusBlock.findElement(By.css("a")).then(
            // Get href of syllabus link on successful find
            async (SyllabusLink: WebElement) => {
                return await SyllabusLink.getAttribute("href");
            },
            // Return null on fail
            () => { return null; }
        ));
    }

    async ParseExams(sectionData: SectionData, SectionTable: WebElement) {
        try {
            let ExamBlock: WebElement = await SectionTable.findElement(By.css("#class_exams"));
            let ExamElements: WebElement[] = await ExamBlock.findElements(By.css("li"));
            // Iterate over all exams
            for (let Element of ExamElements) {
                // Unfortunately, string splitting is pretty much the best way to do this part
                let ExamText: string = await Element.getText();
                let LocationElement: WebElement = Element.findElement(By.css("a"));
                // Derive exam data from split string
                sectionData.Exams.push({
                    Date: ParsingUtils.GetTextBetween(ExamText, "Date:", "Time:"),
                    Time: ParsingUtils.GetTextBetween(ExamText, "Time:", "Location:"),
                    Location: {
                        Room: await LocationElement.getText(),
                        MapURI: await LocationElement.getAttribute("href")
                    }
                });
            };
        } catch (error: NoSuchElementError) { }
    }

    async ParseMeetingTimes(sectionData: SectionData, SectionTable: WebElement) {
        let MeetingBlock: WebElement = await SectionTable.findElement(By.css("p.courseinfo__meeting-time"));
        let MeetingText: string = await MeetingBlock.getText();
        let MeetingData: string[] = MeetingText.split('\n');
        let Dates: string[] = MeetingData[0].split('-');
        sectionData.StartDate = Dates[0];
        sectionData.EndDate = Dates[1];
        sectionData.MeetingDays = MeetingData[1];
        sectionData.Times = MeetingData[2];
        let LocationElement: WebElement = MeetingBlock.findElement(By.css("a"));
        sectionData.Location = {
            Room: await LocationElement.getText(),
            MapURI: await LocationElement.getAttribute("href")
        };
    }

    async ParseAssistants(sectionData: SectionData, TableData: WebElement[], TableDataStrings: string[]) {
        let AssistantBlock: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "TA/RA(s):");
        let AssistantElements = await AssistantBlock.findElements(By.css("div > div"));
        // Iterate over all TAs/RAs
        for (let Element of AssistantElements) {
            // Unfortunately, string splitting is pretty much the best way to do this part
            let AssistantData: string[] = (await Element.getText()).split(" ・ ");
            // Derive instructor data from split string
            sectionData.Assistants.push({
                Name: AssistantData[0],
                Role: AssistantData[1],
                EMail: AssistantData[2]
            });
        };
    }

    async ParseInstructors(sectionData: SectionData, TableData: WebElement[], TableDataStrings: string[]) {
        let InstructorBlock: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Instructor(s):");
        let InstructorElements = await InstructorBlock.findElements(By.css("div > div"));
        // Iterate over potentially multiple instructors
        for (let Element of InstructorElements) {
            // Unfortunately, string splitting is pretty much the best way to do this part
            let InstructorData = (await Element.getText()).split(" ・ ");
            // Derive instructor data from split string
            sectionData.Instructors.push({
                Name: InstructorData[0],
                Role: InstructorData[1],
                EMail: InstructorData[2]
            });
        }
    }

    async ParseRequisites(sectionData: SectionData, TableData: WebElement[], TableDataStrings: string[]) {
        let RequisitesBlock: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Enrollment Reqs:");
        // Class requisites may not always be available (or not in the requisites block)
        if (RequisitesBlock != null) {
            let RequisiteElements = await RequisitesBlock.findElements(By.css("li"));
            for (let Element of RequisiteElements) {
                sectionData.Requisites.push(await Element.getText());
            };
        };
    }

    async ParseAttributes(sectionData: SectionData, TableData: WebElement[], TableDataStrings: string[]) {
        let AttributesBlock: WebElement = ParsingUtils.FindLabeledElement(TableData, TableDataStrings, "Class Attributes:");
        // Class attributes may not always be available
        if (AttributesBlock != null) {
            let AttributeElements = await AttributesBlock.findElements(By.css("li"));
            for (let Element of AttributeElements) {
                sectionData.Attributes.push(await Element.getText());
            };
        };
    }

    // Parse the "section"'s data into course and section data
    async ParseSection(Section: WebElement) {
        // Init the data objects
        let courseData: CourseData = new CourseData();
        let sectionData: SectionData = new SectionData();

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
        let ClassNum = Number.parseInt(SplitNums[0]);
        let CourseNum = Number.parseInt(SplitNums[1]);

        // Find and set course data
        courseData.Title = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Course Title:");
        courseData.College = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "College:");
        this.Courses[CourseNum] = courseData;

        // Find and set section data
        sectionData.Section = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Class Section:");
        sectionData.Level = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Class Level:") as ClassLevel;
        sectionData.Grading = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Grading:");
        sectionData.Consent = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Add Consent:");
        sectionData.Method = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Instruction Mode:");
        sectionData.ActivityType = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Activity Type:");
        sectionData.SessionType = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Session Type:");
        sectionData.Description = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Description:");
        sectionData.CoreInfo = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Core:");

        let TermText: string = await (await SectionTable.findElement(By.css("p.courseinfo__sectionterm"))).getText();
        sectionData.Term = ParsingUtils.GetTextBetween(TermText, "Term: ", "\n");

        let Credits: string = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Semester Credit Hours:");
        // Handle variable-credit courses that may list a non-numerical # of credits
        sectionData.Credits = Number.parseInt(Credits);
        if (isNaN(sectionData.Credits))
            sectionData.Credits = -1;
        // There doesn't seem to be a consistent pattern for the status, so we use GetTextBetween() here
        let StatusText: string = ParsingUtils.FindLabeledText(TableData, TableDataStrings, "Status:");
        sectionData.EnrollmentStatus = ParsingUtils.GetTextBetween(StatusText, "Enrollment Status:", "Available");
        sectionData.AvailableSeats = Number.parseInt(ParsingUtils.GetTextBetween(StatusText, "Seats:", "Enrolled"));
        sectionData.TotalEnrolled = ParsingUtils.GetTextBetween(StatusText, "Total:", "Waitlist");
        sectionData.Waitlisted = Number.parseInt(StatusText.split("Waitlist:")[1]);

        // Get the section's attributes
        await this.ParseAttributes(sectionData, TableData, TableDataStrings);

        // Get the section's requisites
        await this.ParseRequisites(sectionData, TableData, TableDataStrings);

        // Get the section's instructors
        await this.ParseInstructors(sectionData, TableData, TableDataStrings);

        // Get the section's TAs/RAs
        await this.ParseAssistants(sectionData, TableData, TableDataStrings);

        // Get section's meeting times/dates and location data
        await this.ParseMeetingTimes(sectionData, SectionTable);

        // Get exams (may not exist)
        await this.ParseExams(sectionData, SectionTable);

        // Get the URI to the section's syllabus
        await this.ParseSyllabus(sectionData, TableData, TableDataStrings);

        // Get the section's textbooks (do this last because switching to the textbook tab makes all previous elements stale)
        //await this.ParseTextbooks(sectionData, Section);

        // Add collected section data to section data cache
        this.Sections[ClassNum] = sectionData;
        // Print the collected section data
        console.log(JSON.stringify(sectionData, null, '\t'));
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
                for (let Section of SectionList)
                    await this.ParseSection(Section);
                // Write section and course data to data output after all sections under the given prefix are parsed
                writeFileSync("./data/Sections.json", JSON.stringify(this.GetSections(), null, '\t'), { flag: 'a+' });
                writeFileSync("./data/Courses.json", JSON.stringify(this.GetCourses(), null, '\t'), { flag: 'a+' });
                // Clear the buffer
                this.Clear();
            };
        };
    };

    Clear(): void {
        this.Courses = {};
        this.Sections = {};
    }

    GetCourses(): Courses { return this.Courses; };
    GetSections(): Sections { return this.Sections; };
};

// Load Selenium config
/*
let options = new firefox.Options();
const service = new firefox.ServiceBuilder(process.env.SELENIUM_DRIVER);
let CBScraper = new CoursebookScraper(options, service);

CBScraper.Scrape(/2022 Spring/g).then(() => {
    CBScraper.Kill();
});*/
