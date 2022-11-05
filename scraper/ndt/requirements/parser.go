package requirements

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const DEFAULT_MIN_GRADE = "D-"

// Core parsing
func Parse(path string, deg *Degree) *CollectionRequirement {
	// Open input for reading
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	// Init decoder and make sure we can read into the root array
	decoder := json.NewDecoder(f)
	_, err = decoder.Token()
	if err != nil {
		panic(err)
	}
	// Init shell CollectionRequirement to be populated later
	collectionReq := NewCollectionRequirement("DegreeRequirements", 0, nil)
	// Iterate over and parse the elements of the root array
	for decoder.More() {
		var v interface{}
		err := decoder.Decode(&v)
		if err != nil {
			panic(errors.New("Error decoding JSON. Try validating it here: https://jsonlint.com/"))
		}
		switch reflect.TypeOf(v).Kind() {
		// For JSON strings
		case reflect.String:
			collectionReq.Options = append(collectionReq.Options, parseString(v.(string))...)
		// For JSON arrays
		case reflect.Slice:
			collectionReq.Options = append(collectionReq.Options, parseSlice(v.([]interface{})))
		// For JSON objects
		case reflect.Map:
			parseMap(deg, v.(map[string]interface{}))
		}
	}
	// Set collectionReq to require all options (root collection is always an AND collection)
	collectionReq.Required = len(collectionReq.Options)
	return collectionReq
}

// Parse degree field values values from top-level maps to modify degree
func parseMap(deg *Degree, m map[string]interface{}) {
	degType := reflect.TypeOf(*deg)
	for i := 0; i < degType.NumField(); i++ {
		field := degType.Field(i)
		val, hasVal := m[strings.ToLower(field.Name)]
		if hasVal {
				reflect.ValueOf(deg).Elem().FieldByName(field.Name).Set(reflect.ValueOf(val).Convert(field.Type))
		}
	}
}

// Init a CollectionRequirement's requirements by parsing elements as strings and slices
func parseCollectionReqs(elements []interface{}, collectionReq *CollectionRequirement) *CollectionRequirement {
	for _, v := range elements {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			collectionReq.Options = append(collectionReq.Options, parseString(v.(string))...)
		case reflect.Slice:
			collectionReq.Options = append(collectionReq.Options, parseSlice(v.([]interface{})))
		}
	}
	return collectionReq
}

// Parse simple strings as courses, sets of courses, or ranges of courses with a default grade requirement
func parseString(str string) []interface{} {
	// Split string into course type (i.e. "CS") and number (i.e. "1337"), panic if we cannot
	splitString := strings.SplitN(str, " ", 2)
	if len(splitString) != 2 {
		panic(errors.New(fmt.Sprintf("Invalid course string provided: %s", str)))
	}
	courseType := splitString[0]
	courseNums := splitString[1]
	// Split course number with ", " to get course set
	courseSet := strings.Split(courseNums, ", ")
	courses := make([]interface{}, 0, len(courseSet))
	for _, courseNum := range courseSet {
		// Split course number with ... delim to detect a range
		splitString = strings.SplitN(courseNum, "...", 2)
		splitCount := len(splitString)
		if splitCount == 1 {
			// splitCount is 1, string is not a range
			courses = append(courses, NewCourseRequirement(strings.Join([]string{courseType, splitString[0]}, " "), DEFAULT_MIN_GRADE))
		} else {
			// splitCount is 2, string is a range, parse and return as array of courses
			startNum, err := strconv.Atoi(splitString[0])
			if err != nil {
				panic(fmt.Sprintf("Invalid starting number provided for course range '%s' in course string '%s'", courseNum, str))
			}
			endNum, err := strconv.Atoi(splitString[1])
			if err != nil {
				panic(fmt.Sprintf("Invalid ending number provided for course range '%s' in course string '%s'", courseNum, str))
			}
			if startNum < 0 || endNum < 0 {
				panic(fmt.Sprintf("Negative number(s) provided for course range '%s' in course string '%s'", courseNum, str))
			}
			// Swap numbers if they're not in ascending order
			if endNum < startNum {
				temp := startNum
				startNum = endNum
				endNum = temp
			}
			for cn := startNum; cn <= endNum; cn++ {
				courses = append(courses, NewCourseRequirement(strings.Join([]string{courseType, strconv.Itoa(cn)}, " "), DEFAULT_MIN_GRADE))
			}
		}
	}
	return courses
}

// Parce slices by recursively parsing their elements
func parseSlice(slice []interface{}) interface{} {
	// Slices with no elements, or only a type element, are invalid
	if len(slice) <= 1 {
		panic(errors.New(fmt.Sprintf("Invalid slice provided: %v", slice)))
	}
	// Store elements except for the type element
	elements := slice[1:]
	// Switch based on the type element to handle parsing
	switch strings.ToLower(slice[0].(string)) {
	case "and", "&":
		collectionReq := parseCollectionReqs(elements, NewCollectionRequirement("AND", 0, nil))
		// Set required count after parsing to properly handle course expansion
		collectionReq.Required = len(collectionReq.Options)
		return collectionReq
	case "or", "|":
		return parseCollectionReqs(elements, NewCollectionRequirement("OR", 1, nil))
	case "some":
		// For this type, the first element is the required #
		return parseCollectionReqs(elements[1:], NewCollectionRequirement("SOME", int(elements[0].(float64)), nil))
	case "course":
		return NewCourseRequirement(elements[0].(string), elements[1].(string))
	case "section":
		return NewSectionRequirement(elements[0].(string))
	case "exam":
		return NewExamRequirement(elements[0].(string), elements[1].(float64))
	case "major":
		return NewMajorRequirement(elements[0].(string))
	case "minor":
		return NewMinorRequirement(elements[0].(string))
	case "gpa":
		return NewGPARequirement(elements[0].(float64), elements[1].(string))
	case "consent":
		return NewConsentRequirement(elements[0].(string))
	case "other":
		return NewOtherRequirement(elements[0].(string), elements[1].(string))
	case "hours":
		// For this type, the first element is the required #
		// Create temporary collection requirement for parsing purposes
		options := parseCollectionReqs(elements[1:], NewCollectionRequirement("TEMP", 0, nil)).Options
		// Create new []*CourseRequirement slice to hold parsed CourseRequirements
		courseOptions := make([]*CourseRequirement, len(options), len(options))
		// Make sure all interface{}s in options are pointers to CourseRequirements, panic if not
		for i, op := range options {
			courseRequirement, ok := op.(*CourseRequirement)
			if !ok {
				panic(errors.New(fmt.Sprintf("Non-CourseRequirement provided as option in HoursRequirement: %v", op)))
			}
			courseOptions[i] = courseRequirement
		}
		// Use array of CourseRequirement pointers as options for new HoursRequirement
		return NewHoursRequirement(int(elements[0].(float64)), courseOptions)
	case "choice", "choose", "xor":
		return NewChoiceRequirement(parseCollectionReqs(elements, NewCollectionRequirement("XOR", 1, nil)))
	case "limit", "max", "maximum":
		return NewLimitRequirement(int(elements[0].(float64)))
	case "core":
		return NewCoreRequirement(elements[0].(string), int(elements[1].(float64)))
	default:
		panic(errors.New(fmt.Sprintf("Invalid requirement type provided: %s", slice[0].(string))))
	}
}
