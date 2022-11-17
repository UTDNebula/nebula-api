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
			parsedSlice := parseSlice(v.([]interface{}))
			// Ignore nil slices
			if parsedSlice != nil {
				collectionReq.Options = append(collectionReq.Options, parsedSlice)
			}
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
			value := reflect.ValueOf(val)
			if value.CanConvert(field.Type) {
				reflect.ValueOf(deg).Elem().FieldByName(field.Name).Set(value.Convert(field.Type))
			} else {
				panic(fmt.Sprintf("Invalid value \"%v\" of type \"%s\" provided for field \"%s\" of type \"%s\" in map %v", value, value.Type().Name(), field.Name, field.Type.Name(), m))
			}
		}
	}
}

// Parse an array of generic requirements
func parseArray(elements []interface{}) []interface{} {
	parsedReqs := make([]interface{}, 0, len(elements))
	for _, v := range elements {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			parsedReqs = append(parsedReqs, parseString(v.(string))...)
		case reflect.Slice:
			parsedSlice := parseSlice(v.([]interface{}))
			// Ignore nil slices
			if parsedSlice != nil {
				parsedReqs = append(parsedReqs, parsedSlice)
			}
		}
	}
	return parsedReqs
}

// Init a CollectionRequirement's requirements by parsing its elements
func parseCollectionReqs(elements []interface{}, collectionReq *CollectionRequirement) *CollectionRequirement {
	collectionReq.Options = parseArray(elements)
	if len(collectionReq.Options) < collectionReq.Required {
		panic(errors.New(fmt.Sprintf("Insufficient elements provided for collection requirement (# elements < required): %v", collectionReq)))
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
	if len(slice) < 2 {
		panic(errors.New(fmt.Sprintf("Invalid slice provided: %v", slice)))
	}
	// Store elements except for the type element
	elements := slice[1:]
	// Switch based on the type element to handle parsing
	spliceType, ok := slice[0].(string)
	if !ok {
		panic(errors.New(fmt.Sprintf("Invalid slice type provided: %v", slice)))
	}
	switch strings.ToLower(spliceType) {
	case "and", "&":
		collectionReq := parseCollectionReqs(elements, NewCollectionRequirement("AND", 0, nil))
		// Set required count after parsing to properly handle course expansion
		collectionReq.Required = len(collectionReq.Options)
		return collectionReq
	case "or", "|":
		if len(elements) == 0 {
			panic(errors.New(fmt.Sprintf("No elements provided for or requirement: %v", slice)))
		}
		return parseCollectionReqs(elements, NewCollectionRequirement("OR", 1, nil))
	case "some":
		// For this type, the first element is the required #
		required, ok := elements[0].(float64)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid requirement # provided for some requirement: %v", slice)))
		}
		return parseCollectionReqs(elements[1:], NewCollectionRequirement("SOME", int(required), nil))
	case "course":
		if len(elements) < 2 {
			panic(errors.New(fmt.Sprintf("Insufficient elements provided for course requirement: %v", slice)))
		}
		course, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid course name provided for course requirement: %v", slice)))
		}
		minGrade, ok := elements[1].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid minimum grade provided for course requirement: %v", slice)))
		}
		return NewCourseRequirement(course, minGrade)
	case "section":
		sectionRef, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid section reference provided for section requirement: %v", slice)))
		}
		return NewSectionRequirement(sectionRef)
	case "exam":
		if len(elements) < 2 {
			panic(errors.New(fmt.Sprintf("Insufficient elements provided for exam requirement: %v", slice)))
		}
		examRef, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid exam reference provided for exam requirement: %v", slice)))
		}
		minScore, ok := elements[1].(float64)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid minimum score provided for exam requirement: %v", slice)))
		}
		return NewExamRequirement(examRef, minScore)
	case "major":
		major, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid major provided for major requirement: %v", slice)))
		}
		return NewMajorRequirement(major)
	case "minor":
		minor, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid minor provided for minor requirement: %v", slice)))
		}
		return NewMinorRequirement(minor)
	case "gpa":
		if len(elements) < 2 {
			panic(errors.New(fmt.Sprintf("Insufficient elements provided for gpa requirement: %v", slice)))
		}
		minGPA, ok := elements[0].(float64)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid minimum GPA provided for GPA requirement: %v", slice)))
		}
		subset, ok := elements[1].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid subset provided for GPA requirement: %v", slice)))
		}
		return NewGPARequirement(minGPA, subset)
	case "consent":
		granter, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid granter provided for consent requirement: %v", slice)))
		}
		return NewConsentRequirement(granter)
	case "other":
		desc, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid description provided for other requirement: %v", slice)))
		}
		if len(elements) >= 2 {
			condition, ok := elements[1].(string)
			if !ok {
				panic(errors.New(fmt.Sprintf("Invalid condition provided for other requirement: %v", slice)))
			}
			return NewOtherRequirement(desc, condition)
		} else {
			return NewOtherRequirement(desc, "")
		}
	case "hours":
		v, isString := elements[0].(string)
		// Handle hours provided as string (range) i.e. "5-12"
		if isString {
			rangeArr := strings.SplitN(v, "-", 2)
			// Handle string not being a proper range
			if len(rangeArr) < 2 {
				// Try to parse the string as a normal integer if no proper range provided, panic if we can't
				hours, err := strconv.Atoi(rangeArr[0])
				if err != nil {
					panic(errors.New(fmt.Sprintf("Invalid hour range provided for hours requirement %v", slice)))
				}
				return NewHoursRequirement(hours, hours, parseArray(elements[1:]))
			}
			// Validate min and max values of range
			minHours, err := strconv.Atoi(rangeArr[0])
			if err != nil || minHours < 0 {
				panic(errors.New(fmt.Sprintf("Invalid minimum hour requirement provided for hours requirement %v", slice)))
			}
			maxHours, err := strconv.Atoi(rangeArr[1])
			if err != nil || maxHours < 0 || maxHours < minHours {
				panic(errors.New(fmt.Sprintf("Invalid maximum hour requirement provided for hours requirement %v", slice)))
			}
			return NewHoursRequirement(minHours, maxHours, parseArray(elements[1:]))
			// Handle hours provided as number
		} else if v, isFloat := elements[0].(float64); isFloat && v > 0 {
			return NewHoursRequirement(int(v), int(v), parseArray(elements[1:]))
			// Handle error
		} else {
			panic(errors.New(fmt.Sprintf("Invalid hours requirement provided: %v", slice)))
		}
	case "choice", "choose", "xor":
		return NewChoiceRequirement(parseArray(elements))
	case "limit":
		min, ok := elements[0].(float64)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid minimum requirement provided for limit requirement %v", slice)))
		}
		if len(elements) >= 2 {
			max, ok := elements[1].(float64)
			if !ok {
				panic(errors.New(fmt.Sprintf("Invalid maximum requirement provided for limit requirement %v", slice)))
			}
			return NewLimitRequirement(int(min), int(max))
			// If min and max not given, default to setting min to 0 and set max to value
		} else {
			return NewLimitRequirement(0, int(min))
		}
	case "core":
		if len(elements) < 2 {
			panic(errors.New(fmt.Sprintf("Insufficient elements provided for core requirement: %v", slice)))
		}
		flag, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid flag provided for core requirement %v", slice)))
		}
		hours, ok := elements[1].(float64)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid hours provided for core requirement %v", slice)))
		}
		return NewCoreRequirement(flag, int(hours))
	case "not":
		return NewNotRequirement(parseArray(elements[0:1])[0])
	case "elective":
		electiveType, ok := elements[0].(string)
		if !ok {
			panic(errors.New(fmt.Sprintf("Invalid elective type provided for elective requirement %v", slice)))
		}
		if len(elements) >= 2 {
			hours, ok := elements[1].(float64)
			if !ok {
				panic(errors.New(fmt.Sprintf("Invalid hours provided for elective requirement %v", slice)))
			}
			// Handle optional "level" param given
			if len(elements) >= 3 {
				level, ok := elements[2].(string)
				if !ok {
					panic(errors.New(fmt.Sprintf("Invalid level provided for elective requirement %v", slice)))
				}
				return NewElectiveRequirement(electiveType, int(hours), level)
				// Handle optional "level" param not given
			} else {
				return NewElectiveRequirement(electiveType, int(hours), "")
			}
			// Handle optional "hours" param not given, set hours to -1 to specify checking the outer hours requirement
		} else {
			return NewElectiveRequirement(electiveType, -1, "")
		}
	case "comment":
		// Return nil for comments to ignore them
		return nil
	default:
		panic(errors.New(fmt.Sprintf("Invalid requirement type provided: %s", slice[0].(string))))
	}
}
