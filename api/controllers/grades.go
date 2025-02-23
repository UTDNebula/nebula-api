package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// We want to Filter (Match) ASAP

// --------------------------------------------------------

// Aggregate By Course -> Section:
// ---- Prefix
// ---- Prefix, Number
// ---- Prefix, Number, SectionNumber

// Aggregate By Professor -> Section
// ---- Professor

// Aggregate By Find Course, Find Professor: then Match Section
// ---- Prefix, Professor
// ---- Prefix, Number, Professor
// ---- Prefix, Number, Professor, SectionNumber

// --------------------------------------------------------

// Filter on Course
// ---- Prefix
// ---- Prefix, Number

// Filter on Course then Section
// ---- Prefix, Number, SectionNumber

// Filter on Professor
// ---- Professor

// Filter on Section by Matching Course and Professor IDs
// ---- Prefix, Professor
// ---- Prefix, Number, Professor
// ---- Prefix, Number, Professor, SectionNumber

// 5 Functions

// @Id gradeAggregationBySemester
// @Router /grades/semester [get]
// @Description "Returns grade distributions aggregated by semester"
// @Produce json
// @Param prefix query string false "The course's subject prefix"
// @Param number query string false "The course's official number"
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professors's last name"
// @Param section_number query string false "The number of the section"
// @Success 200 {array} responses.GradeResponse "An array of grade distributions for each semester included"
func GradeAggregationSemester() gin.HandlerFunc {
	return func(c *gin.Context) {
		gradesAggregation("semester", c)
	}
}

// @Id gradeAggregationSectionType
// @Router /grades/semester/sectionType [get]
// @Description "Returns the grade distributions aggregated by semester and broken down into section type"
// @Produce json
// @Param prefix query string false "The course's subject prefix"
// @Param number query string false "The course's official number"
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professors's last name"
// @Param section_number query string false "The number of the section"
// @Success 200 {array} responses.SectionGradeResponse "An array of grade distributions for each section type for each semester included"
func GradesAggregationSectionType() gin.HandlerFunc {
	return func(c *gin.Context) {
		gradesAggregation("section_type", c)
	}
}

// @Id gradeAggregationOverall
// @Router /grades/overall [get]
// @Description "Returns the overall grade distribution"
// @Produce json
// @Param prefix query string false "The course's subject prefix"
// @Param number query string false "The course's official number"
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professors's last name"
// @Param section_number query string false "The number of the section"
// @Success 200 {array} integer "A grade distribution array"
func GradesAggregationOverall() gin.HandlerFunc {
	return func(c *gin.Context) {
		gradesAggregation("overall", c)
	}
}

// base function, returns the grade distribution depending on type of flag
func gradesAggregation(flag string, c *gin.Context) {
	var grades []map[string]interface{}
	var results []map[string]interface{}

	var sectionTypeGrades []responses.GradeData // used to parse the response to section-type endpoints

	var cursor *mongo.Cursor
	var collection *mongo.Collection
	var pipeline mongo.Pipeline

	var sectionMatch bson.D
	var courseMatch bson.D
	var courseFind bson.D
	var professorMatch bson.D
	var professorFind bson.D

	var sampleCourse schema.Course // the sample course with the given prefix and course number parameter
	var sampleCourseFind bson.D    // the filter using prefix and course number to get sample course

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// @TODO: Recommend forcing using first_name and last_name to ensure single professors per query.
	// All professors sharing the name will be aggregated together in the current implementation
	prefix := c.Query("prefix")
	number := c.Query("number")
	section_number := c.Query("section_number")
	first_name := c.Query("first_name")
	last_name := c.Query("last_name")

	professor := (first_name != "" || last_name != "")

	// Find internal_course_number associated with subject_prefix and course_number, which will be used later on
	sampleCourseFind = bson.D{
		{Key: "subject_prefix", Value: prefix},
		{Key: "course_number", Value: number},
	}
	// Parse the queried document into the sample course
	err = courseCollection.FindOne(ctx, sampleCourseFind).Decode(&sampleCourse)
	// If the error is not that there is no matching documents, panic the error
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		panic(err)
	}
	internalCourseNumber := sampleCourse.Internal_course_number

	// arrays of regexes and section types
	typeRegexes := [14]string{"0[0-9][0-9]", "0W[0-9]", "0H[0-9]", "0L[0-9]", "5H[0-9]", "1[0-9][0-9]", "2[0-9][0-9]", "3[0-9][0-9]", "5[0-9][0-9]", "6[0-9][0-9]", "7[0-9][0-9]", "HN[0-9]", "HON", "[0-9]U[0-9]"}
	typeStrings := [14]string{"0xx", "0Wx", "0Hx", "0Lx", "5Hx", "1xx", "2xx", "3xx", "5xx", "6xx", "7xx", "HNx", "HON", "xUx"}

	var branches []bson.D            // for without section pipeline
	var withSectionBranches []bson.D // for with section pipeline
	for i := 0; i < len(typeRegexes); i++ {
		branches = append(branches, bson.D{
			{Key: "case", Value: bson.D{{Key: "$regexMatch", Value: bson.D{
				{Key: "input", Value: "$sections.section_number"},
				{Key: "regex", Value: typeRegexes[i]},
			}}}},
			{Key: "then", Value: typeStrings[i]},
		})

		withSectionBranches = append(withSectionBranches, bson.D{
			{Key: "case", Value: bson.D{{Key: "$regexMatch", Value: bson.D{
				{Key: "input", Value: "$section_number"},
				{Key: "regex", Value: typeRegexes[i]},
			}}}},
			{Key: "then", Value: typeStrings[i]},
		})
	}

	// Stage to look up sections
	lookupSectionsStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}},
	}

	// Stage to unwind sections
	unwindSectionsStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$sections"}}}}

	// Project grade distribution stage
	project := bson.D{
		{Key: "_id", Value: "$sections.academic_session.name"},
		{Key: "grade_distribution", Value: "$sections.grade_distribution"},
	}
	if flag == "section_type" { // add the section_type for each section
		project = append(project, bson.E{Key: "section_type", Value: bson.D{
			{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: branches},
				{Key: "default", Value: "OTHERS"}, // might be cases where section doesn't have type listed
			}},
		}})
	}
	projectGradeDistributionStage := bson.D{{Key: "$project", Value: project}}

	// Stage to project grade distribution with section
	project = bson.D{
		{Key: "_id", Value: "$academic_session.name"},
		{Key: "grade_distribution", Value: "$grade_distribution"},
	}
	if flag == "section_type" { // add the section_type for each section
		project = append(project, bson.E{Key: "section_type", Value: bson.D{
			{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: withSectionBranches},
				{Key: "default", Value: "OTHERS"},
			}},
		}})
	}
	projectGradeDistributionWithSectionsStage := bson.D{{Key: "$project", Value: project}}

	// Stage to unwind grade distribution
	unwindGradeDistributionStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$grade_distribution"},
			{Key: "includeArrayIndex", Value: "ix"},
		}},
	}

	// Stage to group grades
	groupID := bson.D{
		{Key: "academic_session", Value: "$_id"},
		{Key: "ix", Value: "$ix"},
	}
	// add section_type to _id so as to group grades by both academic_session and section_type
	if flag == "section_type" {
		groupID = append(groupID, bson.E{Key: "section_type", Value: "$section_type"})
	}
	groupGradesStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupID},
			{Key: "grades", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
		}},
	}

	// Stage to sort grades
	sort := bson.D{
		{Key: "_id.ix", Value: 1},
		{Key: "_id", Value: 1},
	}
	if flag == "section_type" {
		sort = bson.D{ // add section_type to id
			{Key: "_id.ix", Value: 1},
			{Key: "_id.section_type", Value: 1},
			{Key: "_id", Value: 1},
		}
	}
	sortGradesStage := bson.D{{Key: "$sort", Value: sort}}

	// Stage to sum grades
	sumGradesStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "grades", Value: bson.D{{Key: "$sum", Value: "$grades"}}}}}}

	// Stage to group grade distribution
	groupDistributionID := bson.E{Key: "_id", Value: "$_id.academic_session"}
	if flag == "section_type" {
		groupDistributionID = bson.E{Key: "_id", Value: bson.D{
			{Key: "academic_section", Value: "$_id.academic_session"},
			{Key: "section_type", Value: "$_id.section_type"},
		}}
	}
	groupGradeDistributionStage := bson.D{
		{Key: "$group", Value: bson.D{
			groupDistributionID,
			{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grades"}}},
		}},
	}

	// Additional stages for section-type pipeline
	// Stage to sort the section-type-specific grade distributions before grouping
	sortGradeDistributionsStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "_id.section_type", Value: 1},
			{Key: "_id", Value: 1},
		}},
	}

	// Stage group section-type-specific grade distributions together based on semester
	groupSemesterGradeDistributionsStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.academic_section"},
			{Key: "data", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "type", Value: "$_id.section_type"},
				{Key: "grade_distribution", Value: "$grade_distribution"},
			}}}},
		}},
	}

	switch {
	case prefix != "" && number == "" && section_number == "" && !professor:
		// Filter on Course
		collection = courseCollection

		courseMatch = bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix}}}
		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	case prefix != "" && number != "" && section_number == "" && !professor:
		// Filter on Course
		collection = courseCollection

		// Query using internal_course_number of the documents
		courseMatch := bson.D{{Key: "$match", Value: bson.M{"internal_course_number": internalCourseNumber}}}
		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	case prefix != "" && number != "" && section_number != "" && !professor:
		// Filter on Course then Section
		collection = courseCollection

		// Here we query all the courses with the given internal_couse_number,
		// and then filter on the section_number of those courses
		courseMatch := bson.D{{Key: "$match", Value: bson.M{"internal_course_number": internalCourseNumber}}}
		sectionMatch := bson.D{{Key: "$match", Value: bson.M{"sections.section_number": section_number}}}

		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, sectionMatch, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	case prefix == "" && number == "" && section_number == "" && professor:
		// Filter on Professor
		collection = professorCollection

		// Build professorMatch
		if last_name == "" {
			professorMatch = bson.D{{Key: "$match", Value: bson.M{"first_name": first_name}}}
		} else if first_name == "" {
			professorMatch = bson.D{{Key: "$match", Value: bson.M{"last_name": last_name}}}
		} else {
			professorMatch = bson.D{{Key: "$match", Value: bson.M{"first_name": first_name, "last_name": last_name}}}
		}

		pipeline = mongo.Pipeline{professorMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	case prefix != "" && professor:
		// Filter on Section by Matching Course and Professor IDs

		// Here we get the valid course ids and professor ids
		// and then we perform the grades aggregation against the sections collection,
		// matching on the course_reference and professor

		var profIDs []primitive.ObjectID
		var courseIDs []primitive.ObjectID

		collection = sectionCollection

		// Find valid professor ids
		if last_name == "" {
			professorFind = bson.D{{Key: "first_name", Value: first_name}}
		} else if first_name == "" {
			professorFind = bson.D{{Key: "last_name", Value: last_name}}
		} else {
			professorFind = bson.D{{Key: "first_name", Value: first_name}, {Key: "last_name", Value: last_name}}
		}

		cursor, err = professorCollection.Find(ctx, professorFind)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GradeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &results); err != nil {
			panic(err)
		}

		for _, prof := range results {
			profID := prof["_id"].(primitive.ObjectID)
			profIDs = append(profIDs, profID)
		}

		// Get valid course ids
		if number == "" {
			// if only the prefix is provided, filter only on the prefix
			courseFind = bson.D{{Key: "subject_prefix", Value: prefix}}
		} else {
			// if both prefix and course_number are provided, filter on internal_course_number
			courseFind = bson.D{{Key: "internal_course_number", Value: internalCourseNumber}}
		}

		cursor, err = courseCollection.Find(ctx, courseFind)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GradeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &results); err != nil {
			panic(err)
		}

		for _, course := range results {
			courseID := course["_id"].(primitive.ObjectID)
			courseIDs = append(courseIDs, courseID)
		}

		// Build sectionMatch
		if section_number == "" {
			sectionMatch =
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "course_reference", Value: bson.D{{Key: "$in", Value: courseIDs}}},
					{Key: "professors", Value: bson.D{{Key: "$in", Value: profIDs}}},
				}}}
		} else {
			sectionMatch =
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "course_reference", Value: bson.D{{Key: "$in", Value: courseIDs}}},
					{Key: "professors", Value: bson.D{{Key: "$in", Value: profIDs}}},
					{Key: "section_number", Value: section_number},
				}}}
		}

		pipeline = mongo.Pipeline{sectionMatch, projectGradeDistributionWithSectionsStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	default:
		c.JSON(http.StatusBadRequest, responses.GradeResponse{Status: http.StatusBadRequest, Message: "error", Data: "Invalid query parameters."})
		return
	}

	// if this is for section type, add the 2 additional stages to the pipeline
	if flag == "section_type" {
		pipeline = append(pipeline, sortGradeDistributionsStage)
		pipeline = append(pipeline, groupSemesterGradeDistributionsStage)
	}

	// peform aggregation
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents to appropriate type
	if flag != "section_type" {
		if err = cursor.All(ctx, &grades); err != nil {
			panic(err)
		}
	} else {
		if err = cursor.All(ctx, &sectionTypeGrades); err != nil {
			panic(err)
		}
	}

	if flag == "overall" {
		// combine all semester grade_distributions
		overallResponse := make([]int32, 14)
		for _, sem := range grades {
			if len(sem["grade_distribution"].(primitive.A)) != 14 {
				print("Length of Array: ")
				println(len(sem["grade_distribution"].(primitive.A)))
			}
			for i, grade := range sem["grade_distribution"].(primitive.A) {
				overallResponse[i] += grade.(int32)
			}
		}
		c.JSON(http.StatusOK, responses.GradeResponse{Status: http.StatusOK, Message: "success", Data: overallResponse})
	} else if flag == "semester" {
		c.JSON(http.StatusOK, responses.GradeResponse{Status: http.StatusOK, Message: "success", Data: grades})
	} else if flag == "section_type" {
		c.JSON(http.StatusOK, responses.SectionGradeResponse{Status: http.StatusOK, Message: "success", GradeData: sectionTypeGrades})
	} else {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Endpoint broken"})
	}
}
