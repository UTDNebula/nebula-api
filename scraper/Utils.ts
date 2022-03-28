import 'dotenv/config';
import { readFileSync, writeFileSync } from 'fs';
import { Builder, By, until, WebElement, NoSuchElementError } from 'selenium-webdriver';
import schemas from '../api/schemas';

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
    static GetTextBetween(Text: string, Start: string, End: string, TrimWhitespace: Boolean = true): string {
        return TrimWhitespace ? Text.split(Start, 2)[1].split(End, 2)[0].trim() : Text.split(Start, 2)[1].split(End, 2)[0];
    }

    // any 3 semester credit hour 040 core course.
    // a 060 (American History) core course
    // (?:([0-9]) semester credit hour ))?([0-9]{3}) .* core course

    // Maps requisite text patterns to appropriate parser functions
    // DO NOT CHANGE THE ORDER OF THESE PATTERNS. They are ordered in a meaningful way so as to ensure the most beneficial parsing order.
    // Note that the defined order of operations here has ANDs processed before ORs. This is intentionally done so as to enhance some specific parsing cases.
    // This should have no noticeable effect on the final outcome assuming coursebook properly wraps all sub-expressions in parentheses, however it does not always do so.
    // For example, coursebook can sometimes have requisites like "A or B and C", which should really be either "(A or B) and C" or "A or (B and C)".
    // In this case, because of the aforementioned order of operations, the parser will default to interpreting it as "(A or B) and C" since ANDs are processed first. There is not much that can be done about this.
    private static ReqPatternMap: Array<[RegExp, (ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]) => schemas.Requirement]> = [
        [/cannot be received for/i, ParsingUtils.ParseChoice],
        [/grade (?:of )?(?:at least )?(?:a )?([A-f][+-]?) in(?: either)?((?: [A-z]+ [V0-9]{4}(?: or| and)?(?: in| equivalent)?)+)/i, ParsingUtils.ParseGradeChoice],
        // TODO: Revise pattern matching to only evaluate and/or if other patterns are found
        [/ and /i, ParsingUtils.ParseAnd],
        [/ or (?!better|higher)/i, ParsingUtils.ParseOr],
        [/[A-z]+ majors only/i, ParsingUtils.ParseMajor],
        [/[A-z]+ minors only/i, ParsingUtils.ParseMinor],
        [/(?:([0-9]+) semester credit hour )?([0-9]{3}).* core/i, ParsingUtils.ParseCore],
        [/repeated for a maximum of ([0-9]+) semester credit hours/i, ParsingUtils.ParseLimit],
        [/(.+) with a (?:minimum )?grade (?:of )?(?:at least )?(?:a )?([A-f][+-]?)/i, ParsingUtils.ParseGradeList],
        [/^\W*[A-z]+ [V0-9]{4}\W*$/i, ParsingUtils.ParseCourse],
        [/GPA of|grade point average/i, ParsingUtils.ParseGPA],
        [/consent (?=required|of)/i, ParsingUtils.ParseConsent]
    ];

    private static ParsePattern(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.Requirement {
        // Check if the pattern is a nested group
        let Matches: RegExpMatchArray = ReqText.match(/^\W*@([0-9]+)\W*$/);
        if (Matches) // If it is, just return the specified group as a requirement
            return Groups[Matches[1]] as schemas.Requirement;
        // Else, look for an applicable parsing pattern for this text, parse and return if a match is found
        for (let PatternPair of this.ReqPatternMap)
            if (ReqText.match(PatternPair[0]))
                return PatternPair[1](ReqText, Courses, Sections, Groups);
        // Else, if we've exhausted all other options, simply return an OtherRequirement
        let Requirement = new schemas.OtherRequirement();
        Requirement.condition = null;
        Requirement.description = ReqText;
        return Requirement;
    }

    // Replaces part of a group with a subgroup reference for further parsing. Basically the same as inserting parentheses around a chunk of text.
    private static MakeSubgroup(ReqText: string, Groups: any[], ToSubgroup: RegExp | string): void {
        // Find current group's position
        let CurrentGroup = Groups.find((group) => { return group == ReqText });
        let GroupPos = Groups.indexOf(CurrentGroup);
        // Insert group after current group, replacing part of the current group with a subgroup reference
        Groups.splice(GroupPos + 1, 0, ReqText.replace(ToSubgroup, `@${GroupPos}`));
    }

    private static ParseLimit(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.LimitRequirement {
        let Requirement: schemas.LimitRequirement = new schemas.LimitRequirement();
        let LimitMatches: RegExpMatchArray = ReqText.match(/repeated for a maximum of ([0-9]+) semester credit hours/i);
        Requirement.max_hours = Number.parseInt(LimitMatches[1]);
        return Requirement;
    }

    private static ParseCore(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.CoreRequirement {
        let Requirement: schemas.CoreRequirement = new schemas.CoreRequirement();
        let CoreMatches: RegExpMatchArray = ReqText.match(/(?:([0-9]+) semester credit hour )?([0-9]{3}).* core/i);
        if (CoreMatches.length < 3) {// No credit hour requirement, just grab core flag
            Requirement.hours = null;
            Requirement.core_flag = CoreMatches[1];
        }
        else { // Credit hour requirement, grab flag and hours
            Requirement.hours = Number.parseInt(CoreMatches[1]);
            Requirement.core_flag = CoreMatches[2];
        }
        return Requirement;
    }

    private static ParseMinor(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.MinorRequirement {
        let Requirement: schemas.MinorRequirement = new schemas.MinorRequirement();
        let MinorMatches: RegExpMatchArray = ReqText.match(/([A-z]+) minors only/i);
        Requirement.minor = MinorMatches[1];
        return Requirement;
    }

    private static ParseMajor(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.MajorRequirement {
        let Requirement: schemas.MajorRequirement = new schemas.MajorRequirement();
        let MajorMatches: RegExpMatchArray = ReqText.match(/([A-z]+) majors only/i);
        Requirement.major = MajorMatches[1];
        return Requirement;
    }

    private static ParseConsent(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.ConsentRequirement {
        let Requirement: schemas.ConsentRequirement = new schemas.ConsentRequirement();
        let ConsentMatches: RegExpMatchArray = ReqText.trim().match(/(.+)?consent (required|of the)(.+)?/i);
        if (ConsentMatches[2] == "required")
            Requirement.granter = ConsentMatches[1].trim();
        else
            Requirement.granter = ConsentMatches[3].trim();
        return Requirement;
    }

    private static ParseGPA(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.GPARequirement {
        let Requirement: schemas.GPARequirement = new schemas.GPARequirement();
        if (ReqText.includes("or better in")) {
            let GPAMatches: RegExpMatchArray = ReqText.match(/GPA of ([0-9]\.[0-9]+) or better in (.+)/i);
            Requirement.minimum = Number.parseFloat(GPAMatches[1]);
            Requirement.subset = GPAMatches[2];
        } else { // Handle university GPA
                Requirement.minimum = Number.parseFloat(ReqText.match(/[0-9]\.[0-9]+/)[0]);
                Requirement.subset = "university";
        }
        return Requirement;
    }

    private static ParseGradeChoice(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.Requirement {
        let GradeMatches: RegExpMatchArray = ReqText.match(/(?:a )?grade (?:of )?(?:at least )?(?:a )?([A-f][+-]?) in(?: either)?((?: [A-z]+ [V0-9]{4}(?: or| and)?(?: in| equivalent)?)+)/i);
        let Requirement: schemas.Requirement = ParsingUtils.ParsePattern(GradeMatches[2], Courses, Sections, Groups);
        if (Requirement instanceof schemas.CollectionRequirement) {
            for (let Option of Requirement.options)
                if (Option.type == "course")
                    (Option as schemas.CourseRequirement).minimum_grade = GradeMatches[1];
        } else if (Requirement instanceof schemas.CourseRequirement) {
            Requirement.minimum_grade = GradeMatches[1];
        }
        // Make a subgroup to ensure and/or is handled correctly
        ParsingUtils.MakeSubgroup(ReqText, Groups, GradeMatches[0]);
        return Requirement;
    }

    private static ParseGradeList(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.Requirement {
        let GradeMatches: RegExpMatchArray = ReqText.match(/(.+) with a (?:minimum )?grade (?:of )?(?:at least )?(?:a )?([A-f][+-]?)/i);
        let Requirement: schemas.Requirement = ParsingUtils.ParsePattern(GradeMatches[1], Courses, Sections, Groups);
        if (Requirement instanceof schemas.CollectionRequirement) {
            for (let Option of Requirement.options)
                if (Option.type == "course")
                    (Option as schemas.CourseRequirement).minimum_grade = GradeMatches[2];
        } else if (Requirement instanceof schemas.CourseRequirement) {
            Requirement.minimum_grade = GradeMatches[2];
        }
        return Requirement;
    }

    private static ParseCourse(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.CourseRequirement {
        let Requirement: schemas.CourseRequirement = new schemas.CourseRequirement();
        let Matches: RegExpMatchArray = ReqText.match(/([A-z]+) ([V0-9]{4})/); // Find the subject prefix and number
        let MatchingCourse: schemas.PsuedoCourse = Courses.find((Course: schemas.PsuedoCourse) => {
            return (Course.subject_prefix == Matches[1] && Course.course_number == Matches[2]);
        });
        Requirement.class_reference = MatchingCourse?._id ?? null; // Set reference to null if no matching course found
        return Requirement;
    }

    private static ParseOr(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.CollectionRequirement {
        let Requirement: schemas.CollectionRequirement = new schemas.CollectionRequirement();
        let Options: string[] = ReqText.split(/ or /i); // Split text into options
        for (let Option of Options)
            if (!Option.match(/better|higher/i)) // Filter out erroneous text
                Requirement.options.push(ParsingUtils.ParsePattern(Option, Courses, Sections, Groups));
        Requirement.required = 1;
        return Requirement;
    }

    private static ParseAnd(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.CollectionRequirement {
        let Requirement: schemas.CollectionRequirement = new schemas.CollectionRequirement();
        let Options: string[] = ReqText.split(/ and /i); // Split text into options
        for (let Option of Options) {
            Requirement.options.push(ParsingUtils.ParsePattern(Option, Courses, Sections, Groups));
        }
        Requirement.required = Requirement.options.length;
        return Requirement;
    }

    private static ParseChoice(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[], Groups: any[]): schemas.ChoiceRequirement {
        let RelevantText;
        if (ReqText.includes(','))
            RelevantText = ReqText.split(', ')[1]; // Get relevant text after the comma
        else if (ReqText.includes(':'))
            RelevantText = ReqText.split(':')[1]; // Get relevant text after the colon
        else if (ReqText.includes('both'))
            RelevantText = ReqText.split('both')[1]; // Get relevant text after "both"
        else
            RelevantText = ReqText.split("for")[1]; // Get relevant text after "for"
        let Requirement: schemas.ChoiceRequirement = new schemas.ChoiceRequirement();
        // Parse pattern of relevant text, should always be a CollectionRequirement
        Requirement.choices = ParsingUtils.ParsePattern(RelevantText, Courses, Sections, Groups) as schemas.CollectionRequirement;
        Requirement.choices.required = 1;
        // Return created ChoiceRequirement
        return Requirement;
    }

    private static ParseGroups(ReqText: string): string[] {
        if (!ReqText.match(/\(.+\)/)) // If there's no parentheses to parse, return
            return;
        let startPoses: number[] = []; // Array for storing ( positions
        let groups: string[] = []; // Array for storing captured groups
        for (let pos = 0; pos < ReqText.length; ++pos) { // Iterate through chars
            if (ReqText.charAt(pos) == '(') // Increase depth upon finding (
                startPoses.push(pos + 1);
            else if (ReqText.charAt(pos) == ')') { // Decrease depth and push capture group upon finding )
                let startPos = startPoses.pop();
                if (!startPos) // Break in case of broken parenthesis parity
                    break;
                groups.push(ReqText.substring(startPos, pos));
            }
        }
        groups.push(ReqText); // Append full text as last group
        // Replace nested groups with index markers (i.e. @1) in reverse order
        for (let outerPos = groups.length - 1; outerPos > 0; --outerPos)
            for (let nestedPos = outerPos - 1; nestedPos > -1; --nestedPos)
                groups[outerPos] = groups[outerPos].replace(`(${groups[nestedPos]})`, `@${nestedPos}`);
        return groups;
    }

    static ParseReq(ReqText: string, Courses: schemas.PsuedoCourse[], Sections: schemas.Section[]): schemas.Requirement {
        console.log(ReqText);
        let Groups: any[] = ParsingUtils.ParseGroups(ReqText); // Split text into parse-able sub-texts
        if (!Groups) // Set entire text as only group if no groups found
            Groups = [ReqText];
        for (let pos = 0; pos < Groups.length; ++pos) // Parse groups individually, then combine
            Groups[pos] = ParsingUtils.ParsePattern(Groups[pos], Courses, Sections, Groups);
        return Groups[Groups.length - 1]; // Return last group, as it represents the overall parsed req
    }

    static ParseAllReqs() {
        let PsuedoCourses: schemas.PsuedoCourse[] = JSON.parse(readFileSync("./data/CombinedPsuedoCourses.json", { encoding: 'utf-8' }));
        let Sections: schemas.Section[] = JSON.parse(readFileSync("./data/CombinedSections.json", { encoding: 'utf-8' }));
        let ParsedCourses: schemas.Course[] = [];
        for (let Course of PsuedoCourses) {
            let ParsedCourse: object = Course;
            let PreReqs: schemas.CollectionRequirement = new schemas.CollectionRequirement();
            PreReqs.options = Course.prerequisites.map((Req: string) => { return ParsingUtils.ParseReq(Req.trim(), PsuedoCourses, Sections) });
            PreReqs.required = PreReqs.options.length;
            let CoReqs: schemas.CollectionRequirement = new schemas.CollectionRequirement();
            CoReqs.options = Course.corequisites.map((Req: string) => { return ParsingUtils.ParseReq(Req.trim(), PsuedoCourses, Sections) });
            CoReqs.required = CoReqs.options.length;
            let CoOrPreReqs: schemas.CollectionRequirement = new schemas.CollectionRequirement();
            CoOrPreReqs.options = Course.co_or_pre_requisites.map((Req: string) => { return ParsingUtils.ParseReq(Req.trim(), PsuedoCourses, Sections) });
            CoOrPreReqs.required = CoOrPreReqs.options.length;
            ParsedCourse["prerequisites"] = PreReqs;
            ParsedCourse["corequisites"] = CoReqs;
            ParsedCourse["co_or_pre_requisites"] = CoOrPreReqs;
            ParsedCourses.push(ParsedCourse as schemas.Course);
        }
        writeFileSync("./data/CombinedCourses.json", JSON.stringify(ParsedCourses, null, '\t'), { flag: 'w' });
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