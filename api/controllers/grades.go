package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

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

// @Id				gradeAggregationBySemester
// @Router			/grades/semester [get]
// @Tags			Grades
// @Description	"Returns grade distributions aggregated by semester"
// @Produce		json
// @Param			prefix			query		string									false	"The course's subject prefix"
// @Param			number			query		string									false	"The course's official number"
// @Param			first_name		query		string									false	"The professor's first name"
// @Param			last_name		query		string									false	"The professors's last name"
// @Param			section_number	query		string									false	"The number of the section"
// @Success		200				{object}	schema.APIResponse[[]schema.GradeData]	"An array of grade distributions for each semester included"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400				{object}	schema.APIResponse[string]				"A string describing the error"
func GradeAggregationSemester(c *gin.Context) {
	gradesAggregation("semester", c)
}

// @Id				gradeAggregationSectionType
// @Router			/grades/semester/sectionType [get]
// @Tags			Grades
// @Description	"Returns the grade distributions aggregated by semester and broken down into section type"
// @Produce		json
// @Param			prefix			query		string										false	"The course's subject prefix"
// @Param			number			query		string										false	"The course's official number"
// @Param			first_name		query		string										false	"The professor's first name"
// @Param			last_name		query		string										false	"The professors's last name"
// @Param			section_number	query		string										false	"The number of the section"
// @Success		200				{object}	schema.APIResponse[[]schema.TypedGradeData]	"An array of grade distributions for each section type for each semester included"
// @Failure		500				{object}	schema.APIResponse[string]					"A string describing the error"
// @Failure		400				{object}	schema.APIResponse[string]					"A string describing the error"
func GradesAggregationSectionType(c *gin.Context) {
	gradesAggregation("section_type", c)
}

// @Id				gradeAggregationOverall
// @Router			/grades/overall [get]
// @Tags			Grades
// @Description	"Returns the overall grade distribution"
// @Produce		json
// @Param			prefix			query		string						false	"The course's subject prefix"
// @Param			number			query		string						false	"The course's official number"
// @Param			first_name		query		string						false	"The professor's first name"
// @Param			last_name		query		string						false	"The professors's last name"
// @Param			section_number	query		string						false	"The number of the section"
// @Success		200				{object}	schema.APIResponse[[]int]	"A grade distribution array"
// @Failure		500				{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400				{object}	schema.APIResponse[string]	"A string describing the error"
func GradesAggregationOverall(c *gin.Context) {
	gradesAggregation("overall", c)
}

// @Id				GradesByCourseID
// @Router			/course/{id}/grades [get]
// @Tags			Courses
// @Description	"Returns the overall grade distribution for a course"
// @Produce		json
// @Param			id	path		string						true	"ID of course to get grades for"
// @Success		200	{object}	schema.APIResponse[[]int]	"A grade distribution array for the course"
// @Failure		500	{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]	"A string describing the error"
func GradesByCourseID(c *gin.Context) {
	gradesAggregation("course_endpoint", c)
}

// @Id				GradesBySectionID
// @Router			/section/{id}/grades [get]
// @Tags			Sections
// @Description	"Returns the overall grade distribution for a section"
// @Produce		json
// @Param			id	path		string						true	"ID of section to get grades for"
// @Success		200	{object}	schema.APIResponse[[]int]	"A grade distribution array for the section"
// @Failure		500	{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]	"A string describing the error"
func GradesBySectionID(c *gin.Context) {
	gradesAggregation("section_endpoint", c)
}

// @Id				GradesByProfessorID
// @Router			/professor/{id}/grades [get]
// @Tags			Professors
// @Description	"Returns the overall grade distribution for a professor"
// @Produce		json
// @Param			id	path		string						true	"ID of professor to get grades for"
// @Success		200	{object}	schema.APIResponse[[]int]	"A grade distribution array for the professor"
// @Failure		500	{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]	"A string describing the error"
func GradesByProfessorID(c *gin.Context) {
	gradesAggregation("professor_endpoint", c)
}

// gradesAggregation is base function, returns the grade distribution depending on type of flag
func gradesAggregation(flag string, c *gin.Context) {
	var grades []schema.GradeData
	var results []map[string]any

	var sectionTypeGrades []schema.TypedGradeData // used to parse the response to section-type endpoints

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

	var objId *primitive.ObjectID

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//	@TODO:	Recommend forcing using first_name and last_name to ensure single professors per query.
	// All professors sharing the name will be aggregated together in the current implementation
	prefix := c.Query("prefix")
	number := c.Query("number")
	section_number := c.Query("section_number")
	first_name := c.Query("first_name")
	last_name := c.Query("last_name")

	if flag == "course_endpoint" || flag == "section_endpoint" || flag == "professor_endpoint" {
		// parse object id from id parameter
		objId, err = objectIDFromParam(c, "id")
		if err != nil {
			return
		}
	}

	professor := (first_name != "" || last_name != "")

	// Find internal_course_number associated with subject_prefix and course_number, which will be used later on
	sampleCourseFind = bson.D{
		{Key: "subject_prefix", Value: prefix},
		{Key: "course_number", Value: number},
	}

	err = courseCollection.FindOne(ctx, sampleCourseFind).Decode(&sampleCourse)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		// If the error is not that there is no matching documents,
		// throw an internal server error
		respondWithInternalError(c, err)
		return
	}
	internalCourseNumber := sampleCourse.Internal_course_number

	// Arrays of regexes and section types
	typeRegexes := [14]string{"0[0-9][0-9]", "0W[0-9]", "0H[0-9]", "0L[0-9]", "5H[0-9]", "1[0-9][0-9]", "2[0-9][0-9]", "3[0-9][0-9]", "5[0-9][0-9]", "6[0-9][0-9]", "7[0-9][0-9]", "HN[0-9]", "HON", "[0-9]U[0-9]"}
	typeStrings := [14]string{"0xx", "0Wx", "0Hx", "0Lx", "5Hx", "1xx", "2xx", "3xx", "5xx", "6xx", "7xx", "HNx", "HON", "xUx"}

	var branches [14]bson.D            // for without-section pipeline
	var withSectionBranches [14]bson.D // for with-section pipeline

	for i, typeRegex := range typeRegexes {
		branches[i] = bson.D{
			{Key: "case", Value: bson.D{
				{Key: "$regexMatch", Value: bson.D{
					{Key: "input", Value: "$sections.section_number"},
					{Key: "regex", Value: typeRegex},
				}},
			}},
			{Key: "then", Value: typeStrings[i]},
		}

		withSectionBranches[i] = bson.D{
			{Key: "case", Value: bson.D{
				{Key: "$regexMatch", Value: bson.D{
					{Key: "input", Value: "$section_number"},
					{Key: "regex", Value: typeRegex},
				}},
			}},
			{Key: "then", Value: typeStrings[i]},
		}
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
	unwindSectionsStage := bson.D{
		{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$sections"}}},
	}

	// Stage to project grade distribution stage
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
	if flag == "section_type" {
		// Add section_type to _id to group grades by both academic_session and section_type
		groupID = append(groupID, bson.E{Key: "section_type", Value: "$section_type"})
	}
	groupGradesStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupID},
			{Key: "grades", Value: bson.D{
				{Key: "$push", Value: "$grade_distribution"},
			}},
		}},
	}

	// Stage to sort grades
	sort := bson.D{
		{Key: "_id.ix", Value: 1},
		{Key: "_id", Value: 1},
	}
	if flag == "section_type" {
		// insert the section-type criteria to index 1
		sort = append(sort[:1], append(bson.D{{Key: "_id.section_type", Value: 1}}, sort[1:]...)...)
	}
	sortGradesStage := bson.D{{Key: "$sort", Value: sort}}

	// Stage to sum grades
	sumGradesStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "grades", Value: bson.D{{Key: "$sum", Value: "$grades"}}}}}}

	// Stage to group grade distribution
	var groupDistributionID any = "$_id.academic_session"
	if flag == "section_type" {
		// Add the section-type criteria
		groupDistributionID = bson.D{
			{Key: "academic_section", Value: "$_id.academic_session"},
			{Key: "section_type", Value: "$_id.section_type"},
		}
	}
	groupGradeDistributionStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupDistributionID},
			{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grades"}}},
		}},
	}

	// Additional stages for "section-type" pipeline
	// Stage to sort the section-type-specific grade distributions before grouping
	sortGradeDistributionsStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "_id.section_type", Value: 1},
			{Key: "_id", Value: 1},
		}},
	}

	// Stage to group section-type-specific grade distributions together based on semester
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
	case flag == "course_endpoint":
		// Filter on course ID, from the course endpoint
		collection = courseCollection

		courseMatch := bson.D{{Key: "$match", Value: bson.M{"_id": objId}}}
		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	case flag == "section_endpoint":
		// Filter on section ID, from section endpoint
		collection = sectionCollection

		sectionMatch := bson.D{{Key: "$match", Value: bson.M{"_id": objId}}}
		pipeline = mongo.Pipeline{sectionMatch, projectGradeDistributionWithSectionsStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

	case flag == "professor_endpoint":
		// Filter on Professor from professor endpoint
		collection = professorCollection

		professorMatch := bson.D{{Key: "$match", Value: bson.M{"_id": objId}}}
		pipeline = mongo.Pipeline{professorMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

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
			respondWithInternalError(c, err)
			return
		}
		if err = cursor.All(ctx, &results); err != nil {
			respondWithInternalError(c, err)
			return
		}

		profIDs := make([]primitive.ObjectID, 0, len(results))
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
			respondWithInternalError(c, err)
			return
		}
		if err = cursor.All(ctx, &results); err != nil {
			respondWithInternalError(c, err)
			return
		}

		courseIDs := make([]primitive.ObjectID, 0, len(results))
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
		respond(c, http.StatusBadRequest, "error", "Invalid query parameters.")
		return
	}

	// if this is for section type, add the 2 additional stages to the pipeline
	if flag == "section_type" {
		pipeline = append(pipeline, sortGradeDistributionsStage, groupSemesterGradeDistributionsStage)
	}

	// peform aggregation
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// retrieve and parse all valid documents to appropriate type
	if flag != "section_type" {
		if err = cursor.All(ctx, &grades); err != nil {
			respondWithInternalError(c, err)
			return
		}
	} else {
		if err = cursor.All(ctx, &sectionTypeGrades); err != nil {
			respondWithInternalError(c, err)
			return
		}
	}

	switch flag {
	case "overall", "course_endpoint", "section_endpoint", "professor_endpoint":
		// combine all semester grade_distributions
		overallResponse := [14]int{}
		for _, sem := range grades {
			for i, grade := range sem.GradeDistribution {
				overallResponse[i] += grade
			}
		}
		respond(c, http.StatusOK, "success", overallResponse)
	case "semester":
		respond(c, http.StatusOK, "success", grades)
	case "section_type":
		respond(c, http.StatusOK, "success", sectionTypeGrades)
	}
}
