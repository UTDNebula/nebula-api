package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/UTDNebula/nebula-api/toolkit/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Main dictionaries for mapping unique keys to the actual data
var Sections = make(map[schema.IdWrapper]*schema.Section)
var Courses = make(map[string]*schema.Course)
var Professors = make(map[string]*schema.Professor)

// Auxilliary dictionaries for mapping the generated ObjectIDs to the keys used in the above maps, used for validation purposes
var CourseIDMap = make(map[schema.IdWrapper]string)
var ProfessorIDMap = make(map[schema.IdWrapper]string)

// Requisite parser closures associated with courses
var ReqParsers = make(map[schema.IdWrapper]func())

// Grade mappings for section grade distributions, mapping is MAP[SEMESTER] -> MAP[SUBJECT + NUMBER + SECTION] -> GRADE DISTRIBUTION
var GradeMap map[string]map[string][]int

// Time location for dates (uses America/Chicago tz database zone for CDT which accounts for daylight saving)
var timeLocation, timeError = time.LoadLocation("America/Chicago")

// Externally exposed parse function
func Parse(inDir string, outDir string, csvPath string, skipValidation bool) {

	// Panic if timeLocation didn't load properly
	if timeError != nil {
		panic(timeError)
	}

	// Load grade data from csv in advance
	GradeMap = loadGrades(csvPath)
	if len(GradeMap) != 0 {
		fmt.Printf("Loaded grade distributions for %d semesters.\n\n", len(GradeMap))
	}

	// Try to load any existing profile data
	loadProfiles(inDir)

	// Find paths of all scraped data
	paths := getAllSectionFilepaths(inDir)
	fmt.Printf("Parsing and validating %d files...\n", len(paths))

	// Parse all data
	for _, path := range paths {
		parse(path)
	}

	fmt.Printf("\nParsing complete. Created %d courses, %d sections, and %d professors.\n", len(Courses), len(Sections), len(Professors))

	fmt.Printf("\nParsing course requisites...\n")
	// Initialize matchers at runtime for requisite parsing; this is necessary to avoid circular reference errors with compile-time initialization
	initMatchers()
	for _, course := range Courses {
		ReqParsers[course.Id]()
	}
	fmt.Printf("Finished parsing course requisites!\n")

	if !skipValidation {
		fmt.Printf("\nStarting validation stage...\n")
		validate()
		fmt.Printf("\nValidation complete!\n")
	}

	// Make outDir if it doesn't already exist
	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		panic(err)
	}

	// Write validated data to output files
	fptr, err := os.Create(fmt.Sprintf("%s/Courses.json", outDir))
	if err != nil {
		panic(err)
	}
	encoder := json.NewEncoder(fptr)
	encoder.SetIndent("", "\t")
	encoder.Encode(getMapValues(Courses))
	fptr.Close()

	fptr, err = os.Create(fmt.Sprintf("%s/Sections.json", outDir))
	if err != nil {
		panic(err)
	}
	encoder = json.NewEncoder(fptr)
	encoder.SetIndent("", "\t")
	encoder.Encode(getMapValues(Sections))
	fptr.Close()

	fptr, err = os.Create(fmt.Sprintf("%s/Professors.json", outDir))
	if err != nil {
		panic(err)
	}
	encoder = json.NewEncoder(fptr)
	encoder.SetIndent("", "\t")
	encoder.Encode(getMapValues(Professors))
	fptr.Close()
}

func loadProfiles(inDir string) {
	fptr, err := os.Open(fmt.Sprintf("%s/Profiles.json", inDir))
	if err != nil {
		fmt.Printf("Couldn't find/open Profiles.json in the input directory. Skipping profile load.\n")
		return
	}

	decoder := json.NewDecoder(fptr)

	fmt.Printf("Beginning profile load.\n")

	// Read open bracket
	_, err = decoder.Token()
	if err != nil {
		panic(err)
	}

	// While the array contains values
	profileCount := 0
	for ; decoder.More(); profileCount++ {
		// Decode a professor
		var prof schema.Professor
		err := decoder.Decode(&prof)
		if err != nil {
			panic(err)
		}
		professorKey := prof.First_name + prof.Last_name
		Professors[professorKey] = &prof
		ProfessorIDMap[prof.Id] = professorKey
	}

	// Read closing bracket
	_, err = decoder.Token()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d profiles!\n\n", profileCount)
	fptr.Close()
}

// Internal parse function
func parse(path string) {
	fmt.Printf("Parsing %s...\n", path)

	// Open data file for reading
	fptr, err := os.Open(path)
	defer fptr.Close()
	if err != nil {
		panic(err)
	}

	// Create a goquery document for HTML parsing
	doc, err := goquery.NewDocumentFromReader(fptr)
	if err != nil {
		panic(err)
	}

	// Get the rows of the info table
	infoTable := doc.FindMatcher(goquery.Single("table.courseinfo__overviewtable > tbody"))
	infoRows := infoTable.ChildrenFiltered("tr")

	var syllabusURI string

	// Dictionary to hold the row data, keyed by row header
	rowInfo := make(map[string]string, len(infoRows.Nodes))

	// Populate rowInfo
	infoRows.Each(func(_ int, row *goquery.Selection) {
		rowHeader := trimWhitespace(row.FindMatcher(goquery.Single("th")).Text())
		rowData := row.FindMatcher(goquery.Single("td"))
		rowInfo[rowHeader] = trimWhitespace(rowData.Text())
		// Get syllabusURI from syllabus row link
		if rowHeader == "Syllabus:" {
			syllabusURI, _ = rowData.FindMatcher(goquery.Single("a")).Attr("href")
		}
	})

	// Get the rows of the class info subtable
	infoSubTable := infoTable.FindMatcher(goquery.Single("table.courseinfo__classsubtable > tbody"))
	infoRows = infoSubTable.ChildrenFiltered("tr")

	// Dictionary to hold the class info, keyed by data label
	classInfo := make(map[string]string)

	// Populate classInfo
	infoRows.Each(func(_ int, row *goquery.Selection) {
		rowHeaders := row.Find("td.courseinfo__classsubtable__th")
		rowHeaders.Each(func(_ int, header *goquery.Selection) {
			headerText := trimWhitespace(header.Text())
			dataText := trimWhitespace(header.Next().Text())
			classInfo[headerText] = dataText
		})
	})

	// Get the class and course num by splitting classInfo value
	classAndCourseNum := strings.Split(classInfo["Class/Course Number:"], " / ")
	classNum := classAndCourseNum[0]
	courseNum := trimWhitespace(classAndCourseNum[1])

	// Figure out the academic session associated with this specific course/Section
	session := getAcademicSession(rowInfo, classInfo)

	// Try to create the course and section based on collected info
	courseRef := addCourse(courseNum, session, rowInfo, classInfo)
	addSection(courseRef, classNum, syllabusURI, session, rowInfo, classInfo)
	fmt.Printf("Parsed!\n")
}

var coursePrefixRexp *regexp.Regexp = regexp.MustCompile("^([A-Z]{2,4})([0-9V]{4})")
var contactRegexp *regexp.Regexp = regexp.MustCompile("\\(([0-9]+)-([0-9]+)\\)\\s+([SUFY]+)")

func getCatalogYear(session schema.AcademicSession) string {
	sessionYear, err := strconv.Atoi(session.Name[0:2])
	if err != nil {
		panic(err)
	}
	sessionSemester := session.Name[2]
	switch sessionSemester {
	case 'F':
		return strconv.Itoa(sessionYear)
	case 'S':
		return strconv.Itoa(sessionYear - 1)
	case 'U':
		return strconv.Itoa(sessionYear - 1)
	default:
		panic(errors.New(fmt.Sprintf("Encountered invalid session semester '%c!'", sessionSemester)))
	}
}

func addCourse(courseNum string, session schema.AcademicSession, rowInfo map[string]string, classInfo map[string]string) *schema.Course {

	// Courses are internally keyed by their internal course number and the catalog year they're part of
	catalogYear := getCatalogYear(session)
	courseKey := courseNum + catalogYear

	// Don't recreate the course if it already exists
	course, courseExists := Courses[courseKey]
	if courseExists {
		return course
	}

	// Get subject prefix and course number by doing a regexp match on the section id
	sectionId := classInfo["Class Section:"]
	idMatches := coursePrefixRexp.FindStringSubmatch(sectionId)

	course = &schema.Course{}

	course.Id = schema.IdWrapper{Id: primitive.NewObjectID()}
	course.Course_number = idMatches[2]
	course.Subject_prefix = idMatches[1]
	course.Title = rowInfo["Course Title:"]
	course.Description = rowInfo["Description:"]
	course.School = rowInfo["College:"]
	course.Credit_hours = classInfo["Semester Credit Hours:"]
	course.Class_level = classInfo["Class Level:"]
	course.Activity_type = classInfo["Activity Type:"]
	course.Grading = classInfo["Grading:"]
	course.Internal_course_number = courseNum

	// Get closure for parsing course requisites (god help me)
	enrollmentReqs, hasEnrollmentReqs := rowInfo["Enrollment Reqs:"]
	ReqParsers[course.Id] = getReqParser(course, hasEnrollmentReqs, enrollmentReqs)

	// Try to get lecture/lab contact hours and offering frequency from course description
	contactMatches := contactRegexp.FindStringSubmatch(course.Description)
	// Length of contactMatches should be 4 upon successful match
	if len(contactMatches) == 4 {
		course.Lecture_contact_hours = contactMatches[1]
		course.Laboratory_contact_hours = contactMatches[2]
		course.Offering_frequency = contactMatches[3]
	}

	// Set the catalog year
	course.Catalog_year = catalogYear

	Courses[courseKey] = course
	CourseIDMap[course.Id] = courseKey
	return course
}

/*
//	Below is the code for the requisite parser. It is *by far* the most complicated code in this entire project.
//	In summary, it uses a bottom-up "stack"-based parsing technique, building requisites by taking small groups of text, parsing those groups,
//	storing them on the "stack", and then uses those previously parsed groups as dependencies for parsing the larger "higher level" groups.
*/

////////////////////////////////////////////////// BEGIN REQUISITE PARSER CODE //////////////////////////////////////////////////

// Regex matcher object for requisite group parsing
type Matcher struct {
	Regex   *regexp.Regexp
	Handler func(string, []string) interface{}
}

// Regex for group tags
var groupTagRegex = regexp.MustCompile("@(\\d+)")

////////////////////// BEGIN MATCHER FUNCS //////////////////////

var ANDRegex = regexp.MustCompile("(?i)\\s+and\\s+")

func ANDMatcher(group string, subgroups []string) interface{} {
	// Split text along " and " boundaries, then parse subexpressions as groups into an "AND" CollectionRequirement
	subExpressions := ANDRegex.Split(group, -1)
	parsedSubExps := make([]interface{}, 0, len(subExpressions))
	for _, exp := range subExpressions {
		parsedExp := parseGroup(trimWhitespace(exp))
		// Don't include throwaways
		if !reqIsThrowaway(parsedExp) {
			parsedSubExps = append(parsedSubExps, parsedExp)
		}
	}

	parsedSubExps = joinAdjacentOthers(parsedSubExps, " and ")

	if len(parsedSubExps) > 1 {
		return schema.NewCollectionRequirement("AND", len(parsedSubExps), parsedSubExps)
	} else {
		return parsedSubExps[0]
	}
}

// First regex subgroup represents the text to be subgrouped and parsed with parseFnc
// Ex: Text is: "(OPRE 3360 or STAT 3360 or STAT 4351), and JSOM majors and minors only"
// Regex would be: "(JSOM majors and minors only)"
// Resulting substituted text would be: "(OPRE 3360 or STAT 3360 or STAT 4351), and @N", where N is some group number
// When @N is dereferenced from the requisite list, it will have a value equivalent to the result of parseFnc(group, subgroups)

func SubstitutionMatcher(parseFnc func(string, []string) interface{}) func(string, []string) interface{} {
	// Return a closure that uses parseFnc to substitute subgroups[1]
	return func(group string, subgroups []string) interface{} {
		// If there's no text to substitute, just return an OtherRequirement
		if len(subgroups) < 2 {
			return OtherMatcher(group, subgroups)
		}
		// Otherwise, substitute subgroups[1] and parse it with parseFnc
		return parseGroup(makeSubgroup(group, subgroups[1], parseFnc(group, subgroups)))
	}
}

var ORRegex = regexp.MustCompile("(?i)\\s+or\\s+")

func ORMatcher(group string, subgroups []string) interface{} {
	// Split text along " or " boundaries, then parse subexpressions as groups into an "OR" CollectionRequirement
	subExpressions := ORRegex.Split(group, -1)
	parsedSubExps := make([]interface{}, 0, len(subExpressions))
	for _, exp := range subExpressions {
		parsedExp := parseGroup(trimWhitespace(exp))
		// Don't include throwaways
		if !reqIsThrowaway(parsedExp) {
			parsedSubExps = append(parsedSubExps, parsedExp)
		}
	}

	parsedSubExps = joinAdjacentOthers(parsedSubExps, " or ")

	if len(parsedSubExps) > 1 {
		return schema.NewCollectionRequirement("OR", 1, parsedSubExps)
	} else {
		return parsedSubExps[0]
	}
}

func CourseMinGradeMatcher(group string, subgroups []string) interface{} {
	icn, err := findICN(subgroups[1], subgroups[2])
	if err != nil {
		return OtherMatcher(group, subgroups)
	}
	return schema.NewCourseRequirement(icn, subgroups[3])
}

func CourseMatcher(group string, subgroups []string) interface{} {
	icn, err := findICN(subgroups[1], subgroups[2])
	if err != nil {
		return OtherMatcher(group, subgroups)
	}
	return schema.NewCourseRequirement(icn, "F")
}

func ConsentMatcher(group string, subgroups []string) interface{} {
	return schema.NewConsentRequirement(subgroups[1])
}

func LimitMatcher(group string, subgroups []string) interface{} {
	hourLimit, err := strconv.Atoi(subgroups[1])
	if err != nil {
		panic(err)
	}
	return schema.NewLimitRequirement(hourLimit)
}

func MajorMatcher(group string, subgroups []string) interface{} {
	return schema.NewMajorRequirement(subgroups[1])
}

func MinorMatcher(group string, subgroups []string) interface{} {
	return schema.NewMinorRequirement(subgroups[1])
}

func MajorMinorMatcher(group string, subgroups []string) interface{} {
	return schema.NewCollectionRequirement("OR", 1, []interface{}{*schema.NewMajorRequirement(subgroups[1]), *schema.NewMinorRequirement(subgroups[1])})
}

func CoreMatcher(group string, subgroups []string) interface{} {
	hourReq, err := strconv.Atoi(subgroups[1])
	if err != nil {
		panic(err)
	}
	return schema.NewCoreRequirement(subgroups[2], hourReq)
}

func CoreCompletionMatcher(group string, subgroups []string) interface{} {
	return schema.NewCoreRequirement(subgroups[1], -1)
}

func ChoiceMatcher(group string, subgroups []string) interface{} {
	collectionReq, ok := parseGroup(subgroups[1]).(*schema.CollectionRequirement)
	if !ok {
		panic(errors.New(fmt.Sprintf("ChoiceMatcher wasn't able to parse subgroup '%s' into a CollectionRequirement!", subgroups[1])))
	}
	return schema.NewChoiceRequirement(collectionReq)
}

func GPAMatcher(group string, subgroups []string) interface{} {
	GPAFloat, err := strconv.ParseFloat(subgroups[1], 32)
	if err != nil {
		panic(err)
	}
	return schema.NewGPARequirement(GPAFloat, "")
}

func ThrowawayMatcher(group string, subgroups []string) interface{} {
	return schema.Requirement{Type: "throwaway"}
}

func GroupTagMatcher(group string, subgroups []string) interface{} {
	groupIndex, err := strconv.Atoi(subgroups[1])
	if err != nil {
		panic(err)
	}
	// Return a throwaway if index is out of range
	if groupIndex < 0 || groupIndex >= len(requisiteList) {
		return schema.Requirement{Type: "throwaway"}
	}
	// Find referenced group and return it
	parsedGrp := requisiteList[groupIndex]
	return parsedGrp
}

func OtherMatcher(group string, subgroups []string) interface{} {
	return schema.NewOtherRequirement(ungroupText(group), "")
}

/////////////////////// END MATCHER FUNCS ///////////////////////

// Matcher container, matchers must be in order of precedence
// NOTE: PARENTHESES ARE OF HIGHEST PRECEDENCE! (This is due to groupParens() handling grouping of parenthesized text before parsing begins)
var Matchers []Matcher

// Must init matchers via function at runtime to avoid compile-time circular definition error
func initMatchers() {
	Matchers = []Matcher{

		// Throwaways
		Matcher{
			regexp.MustCompile("^(?i)(?:better|\\d-\\d|same as.+)$"),
			ThrowawayMatcher,
		},

		/* TO IMPLEMENT:

		CS or SE Major/Minor

		SUBJECT NUMBER, SUBJECT NUMBER, ..., or SUBJECT NUMBER

		*/

		// "Others", things we need to handle but can't/won't parse
		Matcher{
			regexp.MustCompile("(?i).+(?:freshman|sophomores|juniors|seniors)\\s+only$"),
			OtherMatcher,
		},

		// <SUBJECT> majors and minors only
		Matcher{
			regexp.MustCompile("(?i)(([A-Z]+)\\s+majors\\s+and\\s+minors\\s+only)"),
			SubstitutionMatcher(func(group string, subgroups []string) interface{} {
				return MajorMinorMatcher(subgroups[1], subgroups[1:3])
			}),
		},

		// Core completion
		Matcher{
			regexp.MustCompile("(?i)(Completion\\s+of\\s+(?:an?\\s+)?(\\d{3}).+core(?:\\s+course)?)"),
			SubstitutionMatcher(func(group string, subgroups []string) interface{} {
				return CoreCompletionMatcher(subgroups[1], subgroups[1:3])
			}),
		},

		// Credit cannot be received for both courses, <EXPRESSION>
		Matcher{
			regexp.MustCompile("(?i)(Credit\\s+cannot\\s+be\\s+received\\s+for\\s+both\\s+(?:courses)?,?(.+))"),
			SubstitutionMatcher(func(group string, subgroups []string) interface{} {
				return ChoiceMatcher(subgroups[1], subgroups[1:3])
			}),
		},

		// Logical &
		Matcher{
			ANDRegex,
			ANDMatcher,
		},

		// "<COURSE> with <GRADE> or better"
		Matcher{
			regexp.MustCompile("((?i)([A-Z]{2,4})\\s+([0-9V]{4})\\s+with\\s+a(?:\\s+grade\\s+of)?\\s+([ABCDF][+-]?)\\s+or\\s+better)"), // [name, number, min grade]
			SubstitutionMatcher(func(group string, subgroups []string) interface{} {
				return CourseMinGradeMatcher(subgroups[1], subgroups[1:5])
			}),
		},

		// Logical |
		Matcher{
			ORRegex,
			ORMatcher,
		},

		// <COURSE> with a minimum grade of <GRADE>
		Matcher{
			regexp.MustCompile("^(?i)([A-Z]{2,4})\\s+([0-9V]{4})\\s+with\\s+a\\s+(?:minimum\\s+)?grade\\s+of\\s+(?:at least\\s+)?([ABCDF][+-]?)$"), // [name, number, min grade]
			CourseMinGradeMatcher,
		},
		// A grade of at least <GRADE> in <COURSE>
		Matcher{
			regexp.MustCompile("^(?i)A grade of at least(?: a)? ([ABCDF][+-]?) in ([A-Z]{2,4})\\s+([0-9V]{4})$"), // [min grade, name, number]
			func(group string, subgroups []string) interface{} {
				return CourseMinGradeMatcher(group, []string{subgroups[2], subgroups[3], subgroups[1]})
			},
		},

		// <COURSE>
		Matcher{
			regexp.MustCompile("^([A-Z]{2,4})\\s+([0-9V]{4})"), // [name, number]
			CourseMatcher,
		},

		// <GRANTER> consent required
		Matcher{
			regexp.MustCompile("^(?i)(.+)\\s+consent\\s+required"), // [granter]
			ConsentMatcher,
		},

		// <HOURS> semester credit hours maximum
		Matcher{
			regexp.MustCompile("^(?i)(\\d+)\\s+semester\\s+credit\\s+hours\\s+maximum$"),
			LimitMatcher,
		},
		// This course may only be repeated for <HOURS> credit hours
		Matcher{
			regexp.MustCompile("^(?:[A-Z]{2,4}\\s+[0-9V]{4}\\s+)?Repeat\\s+Limit\\s+-\\s+(?:[A-Z]{2,4}\\s+[0-9V]{4}|This\\s+course)\\s+may\\s+only\\s+be\\s+repeated\\s+for(?:\\s+a\\s+maximum\\s+of)?\\s+(\\d+)\\s+semester\\s+cre?dit\\s+hours(?:\\s+maximum)?$"),
			LimitMatcher,
		},

		// <SUBJECT> majors only
		Matcher{
			regexp.MustCompile("^(?i)(.+)\\s+major(?:s\\s+only)?$"),
			MajorMatcher,
		},

		// <SUBJECT> minors only
		Matcher{
			regexp.MustCompile("^(?i)(.+)\\s+minor(?:s\\s+only)?$"),
			MinorMatcher,
		},

		// Any <HOURS> semester credit hour <CORE> course
		Matcher{
			regexp.MustCompile("^(?i)any\\s+(\\d+)\\s+semester\\s+credit\\s+hour\\s+(\\d{3})(?:\\s+@\\d+)?\\s+core(?:\\s+course)?$"),
			CoreMatcher,
		},

		// Minimum GPA of <GPA>
		Matcher{
			regexp.MustCompile("^(?i)(?:minimum\\s+)?GPA\\s+of\\s+([0-9\\.]+)$"), // [GPA]
			GPAMatcher,
		},
		// <GPA> GPA
		Matcher{
			regexp.MustCompile("^(?i)([0-9\\.]+) GPA$"), // [GPA]
			GPAMatcher,
		},
		// A university grade point average of at least <GPA>
		Matcher{
			regexp.MustCompile("^(?i)a(?:\\s+university)?\\s+grade\\s+point\\s+average\\s+of(?:\\s+at\\s+least)?\\s+([0-9\\.]+)$"), // [GPA]
			GPAMatcher,
		},

		// Group tags (i.e. @1)
		Matcher{
			groupTagRegex, // [group #]
			GroupTagMatcher,
		},
	}
}

var preOrCoreqRegexp *regexp.Regexp = regexp.MustCompile("(?i)((?:Prerequisites? or corequisites?|Corequisites? or prerequisites?):(.*))")
var prereqRegexp *regexp.Regexp = regexp.MustCompile("(?i)(Prerequisites?:(.*))")
var coreqRegexp *regexp.Regexp = regexp.MustCompile("(?i)(Corequisites?:(.*))")

// It is very important that these remain in the same order -- this keeps proper precedence in the below function!
var reqRegexes [3]*regexp.Regexp = [3]*regexp.Regexp{preOrCoreqRegexp, prereqRegexp, coreqRegexp}

// Returns a closure that parses the course's requisites
func getReqParser(course *schema.Course, hasEnrollmentReqs bool, enrollmentReqs string) func() {
	return func() {
		// Pointer array to course requisite properties must be in same order as reqRegexes above
		courseReqs := [3]**schema.CollectionRequirement{&course.Co_or_pre_requisites, &course.Prerequisites, &course.Corequisites}
		// Iterate over and parse each type of requisite, populating the course's relevant requisite property
		for index, reqPtr := range courseReqs {
			// Extract req text from the enrollment req info if it exists, otherwise try using the description
			var reqMatches []string
			if hasEnrollmentReqs {
				course.Enrollment_reqs = enrollmentReqs
				reqMatches = reqRegexes[index].FindStringSubmatch(enrollmentReqs)
			} else {
				reqMatches = reqRegexes[index].FindStringSubmatch(course.Description)
			}
			if reqMatches != nil {
				// Actual useful text is the inner match, index 2
				reqText := reqMatches[2]
				// Erase any sub-matches for other requisite types by matching outer text, index 1
				for _, regex := range reqRegexes {
					matches := regex.FindStringSubmatch(reqText)
					if matches != nil {
						reqText = strings.Replace(reqText, matches[1], "", -1)
					}
				}
				// Split reqText into chunks based on period-space delimiters
				textChunks := strings.Split(trimWhitespace(reqText), ". ")
				parsedChunks := make([]interface{}, 0, len(textChunks))
				// Parse each chunk, then add non-throwaway chunks to parsedChunks
				for _, chunk := range textChunks {
					// Trim any remaining rightmost periods
					chunk = trimWhitespace(strings.TrimRight(chunk, "."))
					parsedChunk := parseChunk(chunk)
					if !reqIsThrowaway(parsedChunk) {
						parsedChunks = append(parsedChunks, parsedChunk)
					}
				}
				// Build CollectionRequirement from parsed chunks and apply to the course property
				if len(parsedChunks) > 0 {
					*reqPtr = schema.NewCollectionRequirement("REQUISITES", len(parsedChunks), parsedChunks)
				}
				fmt.Printf("\n\n")
			}
		}
	}
}

// Function for pulling all requisite references (reqs referenced via group tags) from text
func getReqRefs(text string) []interface{} {
	matches := groupTagRegex.FindAllStringSubmatch(text, -1)
	refs := make([]interface{}, len(matches), len(matches))
	for i, submatches := range matches {
		refs[i] = GroupTagMatcher(submatches[0], submatches)
	}
	return refs
}

// Function for creating a new group by replacing subtext in an existing group, and pushing the new group's info to the req and group list
func makeSubgroup(group string, subtext string, requisite interface{}) string {
	newGroup := strings.Replace(group, subtext, fmt.Sprintf("@%d", len(requisiteList)), -1)
	requisiteList = append(requisiteList, requisite)
	groupList = append(groupList, newGroup)
	return newGroup
}

// Function for joining adjacent OtherRequirements into one OtherRequirement by joining their descriptions with a string
func joinAdjacentOthers(reqs []interface{}, joinString string) []interface{} {
	joinedReqs := make([]interface{}, 0, len(reqs))
	// Temp is a blank OtherRequirement
	temp := *schema.NewOtherRequirement("", "")
	// Iterate over each existing req
	for _, req := range reqs {
		// Determine whether req is an OtherRequirement
		otherReq, isOtherReq := req.(schema.OtherRequirement)
		if !isOtherReq {
			// If temp contains data, append its final result to the joinedReqs
			if temp.Description != "" {
				joinedReqs = append(joinedReqs, temp)
			}
			// Append the non-OtherRequirement to the joinedReqs
			joinedReqs = append(joinedReqs, req)
			// Reset temp's description
			temp.Description = ""
			continue
		}
		// If temp is blank, and req is an otherReq, use otherReq as the initial value of temp
		// Otherwise, join temp's existing description with otherReq's description
		if temp.Description == "" {
			temp = otherReq
		} else {
			temp.Description = strings.Join([]string{temp.Description, otherReq.Description}, joinString)
		}
	}
	// If temp contains data, append its final result to the joinedReqs
	if temp.Description != "" {
		joinedReqs = append(joinedReqs, temp)
	}
	//fmt.Printf("JOINEDREQS ARE: %v\n", joinedReqs)
	return joinedReqs
}

// Function for finding the Internal Course Number associated with the course with the specified subject and course number
func findICN(subject string, number string) (string, error) {
	for _, coursePtr := range Courses {
		if coursePtr.Subject_prefix == subject && coursePtr.Course_number == number {
			return coursePtr.Internal_course_number, nil
		}
	}
	return "ERROR", errors.New(fmt.Sprintf("Couldn't find an ICN for %s %s!", subject, number))
}

// This is the list of produced requisites. Indices coincide with group indices -- aka group @0 will also be the 0th index of the list since it will be processed first.
var requisiteList []interface{}

// This is the list of groups that are to be parsed. They are the raw text chunks associated with the reqs above.
var groupList []string

// Innermost function for parsing individual text groups (used recursively by some Matchers)
func parseGroup(grp string) interface{} {
	// Make sure we trim any mismatched right parentheses
	grp = strings.TrimRight(grp, ")")
	// Find an applicable matcher in Matchers
	for _, matcher := range Matchers {
		matches := matcher.Regex.FindStringSubmatch(grp)
		if matches != nil {
			// If an applicable matcher has been found, return the result of calling its handler
			result := matcher.Handler(grp, matches)
			fmt.Printf("'%s' -> %T\n", grp, result)
			return result
		}
	}
	// Panic if no matcher was able to be found for a given group -- this means we need to add handling for it!!!
	//panic(fmt.Sprintf("NO MATCHER FOUND FOR GROUP '%s'\nSTACK IS: %#v\n", grp, requisiteList))
	//fmt.Printf("NO MATCHER FOR: '%s'\n", grp)
	fmt.Printf("'%s' -> parser.OtherRequirement\n", grp)
	//var temp string
	//fmt.Scanf("%s", temp)
	return *schema.NewOtherRequirement(ungroupText(grp), "")
}

// Outermost function for parsing a chunk of requisite text (potentially containing multiple nested text groups)
func parseChunk(chunk string) interface{} {
	fmt.Printf("\nPARSING CHUNK: '%s'\n", chunk)
	// Extract parenthesized groups from chunk text
	parseText, parseGroups := groupParens(chunk)
	// Initialize the requisite list and group list
	requisiteList = make([]interface{}, 0, len(parseGroups))
	groupList = parseGroups
	// Begin recursive group parsing -- order is bottom-up
	for _, grp := range parseGroups {
		parsedReq := parseGroup(grp)
		// Only append requisite to stack if it isn't marked as throwaway
		if !reqIsThrowaway(parsedReq) {
			requisiteList = append(requisiteList, parsedReq)
		}
	}
	finalGroup := parseGroup(parseText)
	return finalGroup
}

// Check whether a requisite is a throwaway or not by trying a type assertion to Requirement
func reqIsThrowaway(req interface{}) bool {
	baseReq, isBaseReq := req.(schema.Requirement)
	return isBaseReq && baseReq.Type == "throwaway"
}

// Use stack-based parentheses parsing to form text groups and reference them in the original string
func groupParens(text string) (string, []string) {
	var groups []string = make([]string, 0, 5)
	var positionStack []int = make([]int, 0, 5)
	var depth int = 0
	for pos := 0; pos < len(text); pos++ {
		if text[pos] == '(' {
			depth++
			positionStack = append(positionStack, pos)
		} else if text[pos] == ')' && depth > 0 {
			depth--
			lastIndex := len(positionStack) - 1
			// Get last '(' position from stack
			lastPos := positionStack[lastIndex]
			// Pop stack
			positionStack = positionStack[:lastIndex]
			// Make group and replace group text with group index reference
			groupText := text[lastPos+1 : pos]
			groupNum := len(groups)
			groups = append(groups, groupText)
			subText := fmt.Sprintf("@%d", groupNum)
			text = strings.Replace(text, text[lastPos:pos+1], subText, -1)
			// Adjust position to account for replaced text
			pos += len(subText) - len(groupText) - 2
		}
	}
	return text, groups
}

// Function for replacing all group references (groups referenced via group tags) with their actual text
func ungroupText(text string) string {
	text = trimWhitespace(text)
	for groupNum := len(groupList) - 1; groupNum >= 0; groupNum-- {
		subText := fmt.Sprintf("@%d", groupNum)
		replacementText := fmt.Sprintf("(%s)", groupList[groupNum])
		text = strings.Replace(text, subText, replacementText, -1)
	}
	return text
}

/////////////////////////////////////////////////// END REQUISITE PARSER CODE ///////////////////////////////////////////////////

var sectionPrefixRegexp *regexp.Regexp = regexp.MustCompile("^(?i)[A-Z]{2,4}[0-9V]{4}\\.([0-9A-z]+)")
var coreRegexp *regexp.Regexp = regexp.MustCompile("[0-9]{3}")
var personRegexp *regexp.Regexp = regexp.MustCompile("\\s*([\\w ]+)\\s+・\\s+([A-z ]+)\\s+・\\s+([\\w@.]+)")

func addSection(courseRef *schema.Course, classNum string, syllabusURI string, session schema.AcademicSession, rowInfo map[string]string, classInfo map[string]string) {
	// Get subject prefix and course number by doing a regexp match on the section id
	sectionId := classInfo["Class Section:"]
	idMatches := sectionPrefixRegexp.FindStringSubmatch(sectionId)

	section := &schema.Section{}

	section.Id = schema.IdWrapper{Id: primitive.NewObjectID()}
	section.Section_number = idMatches[1]
	section.Course_reference = courseRef.Id

	//TODO: section requisites?

	// Set academic session
	section.Academic_session = session
	// Add professors
	section.Professors = addProfessors(section.Id, rowInfo, classInfo)

	// Get all TA/RA info
	assistantText := rowInfo["TA/RA(s):"]
	assistantMatches := personRegexp.FindAllStringSubmatch(assistantText, -1)
	section.Teaching_assistants = make([]schema.Assistant, 0, len(assistantMatches))
	for _, match := range assistantMatches {
		assistant := schema.Assistant{}
		nameStr := match[1]
		names := strings.Split(nameStr, " ")
		assistant.First_name = names[0]
		assistant.Last_name = names[len(names)-1]
		assistant.Role = match[2]
		assistant.Email = match[3]
		section.Teaching_assistants = append(section.Teaching_assistants, assistant)
	}

	section.Internal_class_number = classNum
	section.Instruction_mode = classInfo["Instruction Mode:"]
	section.Meetings = getMeetings(rowInfo, classInfo)

	// Parse core flags (may or may not exist)
	coreText, hasCore := rowInfo["Core:"]
	if hasCore {
		section.Core_flags = coreRegexp.FindAllString(coreText, -1)
	}

	section.Syllabus_uri = syllabusURI

	semesterGrades, exists := GradeMap[session.Name]
	if exists {
		sectionGrades, exists := semesterGrades[courseRef.Subject_prefix+courseRef.Course_number+section.Section_number]
		if exists {
			section.Grade_distribution = sectionGrades
		}
	}

	// Add new section to section map
	Sections[section.Id] = section

	// Append new section to course's section listing
	courseRef.Sections = append(courseRef.Sections, section.Id)
}

var termRegexp *regexp.Regexp = regexp.MustCompile("Term: ([0-9]+[SUF])")
var datesRegexp *regexp.Regexp = regexp.MustCompile("(?:Start|End)s: ([A-z]+ [0-9]{1,2}, [0-9]{4})")

func getAcademicSession(rowInfo map[string]string, classInfo map[string]string) schema.AcademicSession {
	session := schema.AcademicSession{}
	scheduleText := rowInfo["Schedule:"]

	session.Name = termRegexp.FindStringSubmatch(scheduleText)[1]
	dateMatches := datesRegexp.FindAllStringSubmatch(scheduleText, -1)

	datesFound := len(dateMatches)
	switch {
	case datesFound == 1:
		startDate, err := time.ParseInLocation("January 2, 2006", dateMatches[0][1], timeLocation)
		if err != nil {
			panic(err)
		}
		session.Start_date = startDate
	case datesFound == 2:
		startDate, err := time.ParseInLocation("January 2, 2006", dateMatches[0][1], timeLocation)
		if err != nil {
			panic(err)
		}
		endDate, err := time.ParseInLocation("January 2, 2006", dateMatches[1][1], timeLocation)
		if err != nil {
			panic(err)
		}
		session.Start_date = startDate
		session.End_date = endDate
	}
	return session
}

func addProfessors(sectionId schema.IdWrapper, rowInfo map[string]string, classInfo map[string]string) []schema.IdWrapper {
	professorText := rowInfo["Instructor(s):"]
	professorMatches := personRegexp.FindAllStringSubmatch(professorText, -1)
	var profRefs []schema.IdWrapper = make([]schema.IdWrapper, 0, len(professorMatches))
	for _, match := range professorMatches {

		nameStr := match[1]
		names := strings.Split(nameStr, " ")

		firstName := names[0]
		lastName := names[len(names)-1]

		profKey := firstName + lastName

		prof, profExists := Professors[profKey]
		if profExists {
			prof.Sections = append(prof.Sections, sectionId)
			profRefs = append(profRefs, prof.Id)
			continue
		}

		prof = &schema.Professor{}
		prof.Id = schema.IdWrapper{Id: primitive.NewObjectID()}
		prof.First_name = firstName
		prof.Last_name = lastName
		prof.Titles = []string{match[2]}
		prof.Email = match[3]
		prof.Sections = []schema.IdWrapper{sectionId}
		profRefs = append(profRefs, prof.Id)
		Professors[profKey] = prof
		ProfessorIDMap[prof.Id] = profKey
	}
	return profRefs
}

var meetingsRegexp *regexp.Regexp = regexp.MustCompile(`([A-z]+ [0-9]+, [0-9]{4})-([A-z]+ [0-9]+, [0-9]{4})\W+((?:(?:Mon|Tues|Wednes|Thurs|Fri|Satur|Sun)day(?:, )?)+)\W+([0-9]+:[0-9]+(?:am|pm))-([0-9]+:[0-9]+(?:am|pm))(?:\W+(?:(\S+) (\S+)))`)

func getMeetings(rowInfo map[string]string, classInfo map[string]string) []schema.Meeting {
	scheduleText := rowInfo["Schedule:"]
	meetingMatches := meetingsRegexp.FindAllStringSubmatch(scheduleText, -1)
	var meetings []schema.Meeting = make([]schema.Meeting, 0, len(meetingMatches))
	for _, match := range meetingMatches {
		meeting := schema.Meeting{}

		startDate, err := time.ParseInLocation("January 2, 2006", match[1], timeLocation)
		if err != nil {
			panic(err)
		}
		meeting.Start_date = startDate

		endDate, err := time.ParseInLocation("January 2, 2006", match[2], timeLocation)
		if err != nil {
			panic(err)
		}
		meeting.End_date = endDate

		meeting.Meeting_days = strings.Split(match[3], ", ")

		startTime, err := time.ParseInLocation("3:04pm", match[4], timeLocation)
		if err != nil {
			panic(err)
		}
		meeting.Start_time = startTime

		endTime, err := time.ParseInLocation("3:04pm", match[5], timeLocation)
		if err != nil {
			panic(err)
		}
		meeting.End_time = endTime

		// Only add location data if it's available
		if len(match) > 6 {
			location := schema.Location{}
			location.Building = match[6]
			location.Room = match[7]
			location.Map_uri = fmt.Sprintf("https://locator.utdallas.edu/%s_%s", location.Building, location.Room)
			meeting.Location = location
		}

		meetings = append(meetings, meeting)
	}
	return meetings
}

func validate() {
	// Set up deferred handler for panics to display validation fails
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("VALIDATION FAILED: %s", err)
		}
	}()

	fmt.Printf("\nValidating courses...\n")
	for _, course1 := range Courses {
		// Check for duplicate courses by comparing course_number and subject_prefix as a compound key
		for _, course2 := range Courses {
			// Make sure the course doesn't check itself
			if course1.Internal_course_number == course2.Internal_course_number {
				continue
			}
			if course2.Course_number == course1.Course_number && course2.Subject_prefix == course1.Subject_prefix {
				fmt.Printf("Duplicate course found for %s%s!\n", course1.Subject_prefix, course1.Course_number)
				fmt.Printf("Course 1: %v\n\nCourse 2: %v", course1, course2)
				panic(errors.New("Courses failed to validate!"))
			}
		}
		// Make sure course isn't referencing any nonexistent sections, and that course-section references are consistent both ways
		for _, sectionId := range course1.Sections {
			section, exists := Sections[sectionId]
			if !exists {
				fmt.Printf("Nonexistent section reference found for %s%s!\n", course1.Subject_prefix, course1.Course_number)
				fmt.Printf("Referenced section ID: %s\nCourse ID: %s\n", sectionId, course1.Id)
				panic(errors.New("Courses failed to validate!"))
			}
			if section.Course_reference != course1.Id {
				fmt.Printf("Inconsistent section reference found for %s%s! The course references the section, but not vice-versa!\n", course1.Subject_prefix, course1.Course_number)
				fmt.Printf("Referenced section ID: %s\nCourse ID: %s\nSection course reference: %s\n", sectionId, course1.Id, section.Course_reference)
				panic(errors.New("Courses failed to validate!"))
			}
		}
	}
	fmt.Printf("No invalid courses!\n\n")

	fmt.Printf("Validating sections...\n")
	for _, section1 := range Sections {
		// Check for duplicate sections by comparing section_number, course_reference, and academic_session as a compound key
		for _, section2 := range Sections {
			// Make sure the section doesn't check itself
			if section1.Internal_class_number == section2.Internal_class_number {
				continue
			}
			if section2.Section_number == section1.Section_number &&
				section2.Course_reference == section1.Course_reference &&
				section2.Academic_session == section1.Academic_session {
				fmt.Printf("Duplicate section found!\n")
				fmt.Printf("Section 1: %v\n\nSection 2: %v", section1, section2)
				panic(errors.New("Sections failed to validate!"))
			}
		}
		// Make sure section isn't referencing any nonexistent professors, and that section-professor references are consistent both ways
		for _, profId := range section1.Professors {
			professorKey, exists := ProfessorIDMap[profId]
			if !exists {
				fmt.Printf("Nonexistent professor reference found for section ID %s!\n", section1.Id)
				fmt.Printf("Referenced professor ID: %s\n", profId)
				panic(errors.New("Sections failed to validate!"))
			}
			profRefsSection := false
			for _, profSection := range Professors[professorKey].Sections {
				if profSection == section1.Id {
					profRefsSection = true
					break
				}
			}
			if !profRefsSection {
				fmt.Printf("Inconsistent professor reference found for section ID %s! The section references the professor, but not vice-versa!\n", section1.Id)
				fmt.Printf("Referenced professor ID: %s\n", profId)
				panic(errors.New("Sections failed to validate!"))
			}
		}
		// Make sure section isn't referencing a nonexistant course
		_, exists := CourseIDMap[section1.Course_reference]
		if !exists {
			fmt.Printf("Nonexistent course reference found for section ID %s!\n", section1.Id)
			fmt.Printf("Referenced course ID: %s\n", section1.Course_reference)
			panic(errors.New("Sections failed to validate!"))
		}
	}
	fmt.Printf("No invalid sections!\n\n")

	fmt.Printf("Validating professors...\n")
	// Check for duplicate professors by comparing first_name, last_name, and sections as a compound key
	for _, prof1 := range Professors {
		for _, prof2 := range Professors {
			// Make sure the professor doesn't check itself
			if prof1.Id == prof2.Id {
				continue
			}
			if prof2.First_name == prof1.First_name &&
				prof2.Last_name == prof1.Last_name &&
				prof2.Profile_uri == prof1.Profile_uri {
				fmt.Printf("Duplicate professor found!\n")
				fmt.Printf("Professor 1: %v\n\nProfessor 2: %v", prof1, prof2)
				panic(errors.New("Professors failed to validate!"))
			}
		}
	}
	fmt.Printf("No invalid professors!\n\n")
}

func getAllSectionFilepaths(inDir string) []string {
	var sectionFilePaths []string
	// Try to open inDir
	fptr, err := os.Open(inDir)
	if err != nil {
		panic(err)
	}
	// Try to get term directories in inDir
	termFiles, err := fptr.ReadDir(-1)
	fptr.Close()
	if err != nil {
		panic(err)
	}
	// Iterate over term directories
	for _, file := range termFiles {
		if !file.IsDir() {
			continue
		}
		termPath := fmt.Sprintf("%s/%s", inDir, file.Name())
		fptr, err = os.Open(termPath)
		if err != nil {
			panic(err)
		}
		courseFiles, err := fptr.ReadDir(-1)
		fptr.Close()
		if err != nil {
			panic(err)
		}
		// Iterate over course directories
		for _, file := range courseFiles {
			coursePath := fmt.Sprintf("%s/%s", termPath, file.Name())
			fptr, err = os.Open(coursePath)
			if err != nil {
				panic(err)
			}
			sectionFiles, err := fptr.ReadDir(-1)
			fptr.Close()
			if err != nil {
				panic(err)
			}
			// Get all section file paths from course directory
			for _, file := range sectionFiles {
				sectionFilePaths = append(sectionFilePaths, fmt.Sprintf("%s/%s", coursePath, file.Name()))
			}
		}
	}
	return sectionFilePaths
}

func trimWhitespace(text string) string {
	return strings.Trim(text, " \t\n\r")
}

func getMapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
