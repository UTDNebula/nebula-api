package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"

	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var courseCollection *mongo.Collection = configs.GetCollection("courses")

// @Id				courseSearch
// @Router			/course [get]
// @Tags			Courses
// @Description	"Returns paginated list of courses matching the query's string-typed key-value pairs. See offset for more details on pagination."
// @Produce		json
// @Param			offset					query		number								false	"The starting position of the current page of courses (e.g. For starting at the 17th course, offset=16)."
// @Param			course_number			query		string								false	"The course's official number"
// @Param			subject_prefix			query		string								false	"The course's subject prefix"
// @Param			title					query		string								false	"The course's title"
// @Param			description				query		string								false	"The course's description"
// @Param			school					query		string								false	"The course's school"
// @Param			credit_hours			query		string								false	"The number of credit hours awarded by successful completion of the course"
// @Param			class_level				query		string								false	"The level of education that this course course corresponds to"
// @Param			activity_type			query		string								false	"The type of class this course corresponds to"
// @Param			grading					query		string								false	"The grading status of this course"
// @Param			internal_course_number	query		string								false	"The internal (university) number used to reference this course"
// @Param			lecture_contact_hours	query		string								false	"The weekly contact hours in lecture for a course"
// @Param			offering_frequency		query		string								false	"The frequency of offering a course"
// @Success		200						{object}	schema.APIResponse[[]schema.Course]	"A list of courses"
// @Failure		500						{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400						{object}	schema.APIResponse[string]			"A string describing the error"
func CourseSearch(c *gin.Context) {
	//name := c.Query("name")            	// value of specific query parameter: string
	//queryParams := c.Request.URL.Query() 	// map of all query params: map[string][]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courses []schema.Course

	// Build query key value pairs (only one value per key)
	query, err := getQuery[schema.Course]("Search", c)
	if err != nil {
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	// Get cursor for query results
	cursor, err := courseCollection.Find(ctx, query, optionLimit)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	// Retrieve and parse all valid documents
	if err = cursor.All(ctx, &courses); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// return result
	respond(c, http.StatusOK, "success", courses)
}

// @Id				courseById
// @Router			/course/{id} [get]
// @Tags			Courses
// @Description	"Returns the course with given ID"
// @Produce		json
// @Param			id	path		string								true	"ID of the course to get"
// @Success		200	{object}	schema.APIResponse[schema.Course]	"A course"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
func CourseById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var course schema.Course

	// parse object id from id parameter
	query, err := getQuery[schema.Course]("ById", c)
	if err != nil {
		return
	}

	// find and parse matching course
	err = courseCollection.FindOne(ctx, query).Decode(&course)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No courses with given ID")
		} else {
			respondWithInternalError(c, err)
		}
		return
	}

	// return result
	respond(c, http.StatusOK, "success", course)
}

// @Id				courseAll
// @Router			/course/all [get]
// @Tags			Courses
// @Description	"Returns all courses"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]schema.Course]	"All courses"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
func CourseAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var courses []schema.Course

	cursor, err := courseCollection.Find(ctx, bson.M{})
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &courses); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", courses)
}

// @Id				courseSectionSearch
// @Router			/course/sections [get]
// @Tags			Courses
// @Description	"Returns paginated list of sections of all the courses matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset			query		number									false	"The starting position of the current page of courses (e.g. For starting at the 17th course, former_offset=16)."
// @Param			latter_offset			query		number									false	"The starting position of the current page of sections (e.g. For starting at the 4th section, latter_offset=3)."
// @Param			course_number			query		string									false	"The course's official number"
// @Param			subject_prefix			query		string									false	"The course's subject prefix"
// @Param			title					query		string									false	"The course's title"
// @Param			description				query		string									false	"The course's description"
// @Param			school					query		string									false	"The course's school"
// @Param			credit_hours			query		string									false	"The number of credit hours awarded by successful completion of the course"
// @Param			class_level				query		string									false	"The level of education that this course course corresponds to"
// @Param			activity_type			query		string									false	"The type of class this course corresponds to"
// @Param			grading					query		string									false	"The grading status of this course"
// @Param			internal_course_number	query		string									false	"The internal (university) number used to reference this course"
// @Param			lecture_contact_hours	query		string									false	"The weekly contact hours in lecture for a course"
// @Param			offering_frequency		query		string									false	"The frequency of offering a course"
// @Success		200						{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500						{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400						{object}	schema.APIResponse[string]				"A string describing the error"
func CourseSectionSearch(c *gin.Context) {
	courseSection("Search", c)
}

// @Id				courseSectionById
// @Router			/course/{id}/sections [get]
// @Tags			Courses
// @Description	"Returns the all of the sections of the course with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the course to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func CourseSectionById(c *gin.Context) {
	courseSection("ById", c)
}

// courseSection gets the sections of the courses, filters depending on the flag
func courseSection(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courseSections []schema.Section
	var courseQuery bson.M

	// Determine the course query
	courseQuery, err := getQuery[schema.Course](flag, c)
	if err != nil {
		return
	}

	// Determine the offset and limit for pagination & delete offset fields
	paginateMap, err := configs.GetAggregateLimit(&courseQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	// Pipeline to query the sections from the filtered courses
	courseSectionPipeline := buildCoursePipeline("sections", courseQuery, paginateMap)

	// perform aggregation on the pipeline
	cursor, err := courseCollection.Aggregate(ctx, courseSectionPipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &courseSections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", courseSections)
}

// @Id				courseProfessorSearch
// @Router			/course/professors [get]
// @Tags			Courses
// @Description	"Returns paginated list of professors of all the courses matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset			query		number									false	"The starting position of the current page of courses (e.g. For starting at the 17th course, former_offset=16)."
// @Param			latter_offset			query		number									false	"The starting position of the current page of professors (e.g. For starting at the 4th professor, latter_offset=3)."
// @Param			course_number			query		string									false	"The course's official number"
// @Param			subject_prefix			query		string									false	"The course's subject prefix"
// @Param			title					query		string									false	"The course's title"
// @Param			description				query		string									false	"The course's description"
// @Param			school					query		string									false	"The course's school"
// @Param			credit_hours			query		string									false	"The number of credit hours awarded by successful completion of the course"
// @Param			class_level				query		string									false	"The level of education that this course course corresponds to"
// @Param			activity_type			query		string									false	"The type of class this course corresponds to"
// @Param			grading					query		string									false	"The grading status of this course"
// @Param			internal_course_number	query		string									false	"The internal (university) number used to reference this course"
// @Param			lecture_contact_hours	query		string									false	"The weekly contact hours in lecture for a course"
// @Param			offering_frequency		query		string									false	"The frequency of offering a course"
// @Success		200						{object}	schema.APIResponse[[]schema.Professor]	"A list of professors"
// @Failure		500						{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400						{object}	schema.APIResponse[string]				"A string describing the error"
func CourseProfessorSearch(c *gin.Context) {
	courseProfessor("Search", c)
}

// @Id				courseProfessorById
// @Router			/course/{id}/professors [get]
// @Tags			Courses
// @Description	"Returns the all of the professors of the course with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the course to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Professor]	"A list of professors"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func CourseProfessorById(c *gin.Context) {
	courseProfessor("ById", c)
}

// courseProfessor gets the professors of the courses, filters depending on the flag
func courseProfessor(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courseProfessors []schema.Professor
	var courseQuery bson.M

	courseQuery, err := getQuery[schema.Course](flag, c)
	if err != nil {
		return
	}

	// Determine the offset and limit for pagination and delete offset field
	paginateMap, err := configs.GetAggregateLimit(&courseQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	// Pipeline to query the professors from the filtered courses
	courseProfPipeline := buildCoursePipeline("professors", courseQuery, paginateMap)

	// perform aggregation on the pipeline
	cursor, err := courseCollection.Aggregate(ctx, courseProfPipeline)
	if err != nil {
		// return error for any aggregation problem
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	// parse the array of professors of the course
	if err = cursor.All(ctx, &courseProfessors); err != nil {
		panic(err)
	}

	respond(c, http.StatusOK, "success", courseProfessors)
}

// buildCoursePipeline builds the pipeline to aggregate the list of specified objects from list of courses
func buildCoursePipeline(targetObj string, query bson.M, paginateMap map[string]bson.D) mongo.Pipeline {
	coursePipeline := mongo.Pipeline{}

	// Before looking up the target collection: Filter, then paginage the courses
	preLookupStages := mongo.Pipeline{
		bson.D{{Key: "$match", Value: query}},

		// Skip to the offset, then limit to the number of courses
		paginateMap["former_offset"], paginateMap["limit"],
	}
	coursePipeline = append(coursePipeline, preLookupStages...)

	// Looking up: Aggregate list of target objects from list of courses
	lookupStages := mongo.Pipeline{
		// Lookup the list of sections from the courses
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sections"},
				{Key: "localField", Value: "sections"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections"},
			}},
		},
	}
	if targetObj == "professors" {
		lookupStages = append(lookupStages,
			// Lookup the list of professors from the list of sections
			bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "professors"},
					{Key: "localField", Value: "sections.professors"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "professors"},
				}},
			},
		)
	}
	coursePipeline = append(coursePipeline, lookupStages...)

	// After looking up target collection: Replace with, order, and paginate the looked-up objects
	postLookupStages := mongo.Pipeline{
		// Unwind the target object of the sections
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: fmt.Sprintf("$%s", targetObj)},
				{Key: "preserveNullAndEmptyArrays", Value: false},
			}},
		},

		// Replace the courses with the target objects
		bson.D{{Key: "$replaceWith", Value: fmt.Sprintf("$%s", targetObj)}},

		// Keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},

		// Paginate the target objects
		paginateMap["latter_offset"], paginateMap["limit"],
	}

	coursePipeline = append(coursePipeline, postLookupStages...)
	return coursePipeline
}
