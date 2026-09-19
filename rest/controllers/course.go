package controllers

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/UTDNebula/nebula-api/rest/configs"

	"github.com/UTDNebula/nebula-api/rest/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var courses []schema.Course

	// Build query key value pairs (only one value per key)
	query, err := getQuery[schema.Course]("Search", c)
	if err != nil {
		return
	}

	offset, limit, err := configs.GetLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}
	opts := options.Find().SetSort(getSort("Course")).SetSkip(offset).SetLimit(limit)

	// Get cursor for query results
	cursor, err := courseCollection.Find(ctx, query, opts)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
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
	courseAggregate[schema.Section]("Search", c)
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
	courseAggregate[schema.Section]("ById", c)
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
	courseAggregate[schema.Professor]("Search", c)
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
	courseAggregate[schema.Professor]("ById", c)
}

// courseAggregate returns the list of aggregated objects from list of filtered courses
func courseAggregate[T any](flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var queryResults []T
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

	// Pipeline to query the field from the filtered courses
	schemaType := strings.Split(reflect.TypeFor[[]T]().String(), ".")[1]
	courseQueryPipeline := buildCoursePipeline(schemaType, courseQuery, paginateMap)

	// perform aggregation on the pipeline
	cursor, err := courseCollection.Aggregate(ctx, courseQueryPipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &queryResults); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", queryResults)
}

// buildCoursePipeline builds the pipeline to aggregate the list of object from list of courses
func buildCoursePipeline(schemaType string, courseQuery bson.M, paginateMap map[string]bson.D) mongo.Pipeline {
	field := typeToField[schemaType]
	filterCourse := mongo.Pipeline{
		bson.D{{Key: "$match", Value: courseQuery}},

		bson.D{{Key: "$sort", Value: getSort("Course")}},
		paginateMap["former_offset"],
		paginateMap["limit"],
	}

	var lookup = mongo.Pipeline{
		// Lookup the list of sections from the courses
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},
	}
	var dedup mongo.Pipeline
	switch schemaType {
	case "Section":
		// No extra stages

	case "Professor":
		// Lookup the list of professors from the list of sections
		lookup = append(lookup,
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "sections.professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "professors"},
			}}})

		// Remove the duplicate professors
		dedup = mongo.Pipeline{
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "professors", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
			}}},

			bson.D{{Key: "$replaceWith", Value: "$professors"}},
		}

	default:
		panic("invalid schema for coursePipeline: " + schemaType)
	}

	extract := mongo.Pipeline{
		// Unwind the target objects
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + field},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		// Replace the courses with the target objects
		bson.D{{Key: "$replaceWith", Value: "$" + field}},
	}

	paginate := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: getSort(schemaType)}},
		paginateMap["latter_offset"],
		paginateMap["limit"],
	}

	pipeline := filterCourse
	pipeline = append(pipeline, lookup...)
	pipeline = append(pipeline, extract...)
	pipeline = append(pipeline, dedup...)
	pipeline = append(pipeline, paginate...)
	return pipeline
}
