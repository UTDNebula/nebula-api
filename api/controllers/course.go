package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var courseCollection *mongo.Collection = configs.GetCollection("courses")

// @Id courseSearch
// @Router /course [get]
// @Description "Returns paginated list of courses matching the query's string-typed key-value pairs. See offset for more details on pagination."
// @Produce json
// @Param offset query number false "The starting position of the current page of courses (e.g. For starting at the 17th course, offset=16)."
// @Param course_number query string false "The course's official number"
// @Param subject_prefix query string false "The course's subject prefix"
// @Param title query string false "The course's title"
// @Param description query string false "The course's description"
// @Param school query string false "The course's school"
// @Param credit_hours query string false "The number of credit hours awarded by successful completion of the course"
// @Param class_level query string false "The level of education that this course course corresponds to"
// @Param activity_type query string false "The type of class this course corresponds to"
// @Param grading query string false "The grading status of this course"
// @Param internal_course_number query string false "The internal (university) number used to reference this course"
// @Param lecture_contact_hours query string false "The weekly contact hours in lecture for a course"
// @Param offering_frequency query string false "The frequency of offering a course"
// @Success 200 {array} schema.Course "A list of courses"
func CourseSearch(c *gin.Context) {
	//name := c.Query("name")            	// value of specific query parameter: string
	//queryParams := c.Request.URL.Query() 	// map of all query params: map[string][]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courses []schema.Course

	// build query key value pairs (only one value per key)
	query, err := schema.FilterQuery[schema.Course](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ErrorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// get cursor for query results
	cursor, err := courseCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &courses); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	// return result
	log.Logger.Print(len(courses))
	c.JSON(http.StatusOK, responses.MultiCourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
}

// @Id courseById
// @Router /course/{id} [get]
// @Description "Returns the course with given ID"
// @Produce json
// @Param id path string true "ID of the course to get"
// @Success 200 {object} schema.Course "A course"
func CourseById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := c.Param("id")

	var course schema.Course

	// parse object id from id parameter
	objId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	// find and parse matching course
	err = courseCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&course)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// return result
	c.JSON(http.StatusOK, responses.SingleCourseResponse{Status: http.StatusOK, Message: "success", Data: course})
}

func CourseAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	var courses []schema.Course

	defer cancel()

	cursor, err := courseCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &courses); err != nil {
		panic(err)
	}

	// return result
	c.JSON(http.StatusOK, responses.MultiCourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
}

// @Id courseSectionSearch
// @Router /course/sections [get]
// @Description "Returns paginated list of sections of all the courses matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce json
// @Param former_offset query number false "The starting position of the current page of courses (e.g. For starting at the 17th course, former_offset=16)."
// @Param latter_offset query number false "The starting position of the current page of sections (e.g. For starting at the 4th section, latter_offset=3)."
// @Param course_number query string false "The course's official number"
// @Param subject_prefix query string false "The course's subject prefix"
// @Param title query string false "The course's title"
// @Param description query string false "The course's description"
// @Param school query string false "The course's school"
// @Param credit_hours query string false "The number of credit hours awarded by successful completion of the course"
// @Param class_level query string false "The level of education that this course course corresponds to"
// @Param activity_type query string false "The type of class this course corresponds to"
// @Param grading query string false "The grading status of this course"
// @Param internal_course_number query string false "The internal (university) number used to reference this course"
// @Param lecture_contact_hours query string false "The weekly contact hours in lecture for a course"
// @Param offering_frequency query string false "The frequency of offering a course"
// @Success 200 {array} schema.Section "A list of sections"
func CourseSectionSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseSection("Search", c)
	}
}

// @Id courseSectionById
// @Router /course/{id}/sections [get]
// @Description "Returns the all of the sections of the course with given ID"
// @Produce json
// @Param id path string true "ID of the course to get"
// @Success 200 {array} schema.Section "A list of sections"
func CourseSectionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseSection("ById", c)
	}
}

// get the sections of the courses, filters depending on the flag
func courseSection(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courseSections []schema.Section // the list of sections of the filtered courses
	var courseQuery bson.M              // query of the courses (or the single course)
	var err error                       // error

	// determine the course query
	if flag == "Search" { // filter courses based on the query parameters
		// build the key-value pair of query parameters
		courseQuery, err = schema.FilterQuery[schema.Course](c)
		if err != nil {
			// return the validation error if there's anything wrong
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
			return
		}
	} else if flag == "ById" { // filter the single course based on it's Id
		// convert the id param with the ObjectID
		courseId := c.Param("id")
		courseObjId, convertIdErr := primitive.ObjectIDFromHex(courseId)
		if convertIdErr != nil {
			// return the id conversion error if there's error
			log.WriteError(convertIdErr)
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error with id", Data: convertIdErr.Error()})
			return
		}
		courseQuery = bson.M{"_id": courseObjId}
	} else {
		err = errors.New("invalid type of filtering courses, either filtering based on available course fields or ID")
		// otherwise, something that messed up the server
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "internal error", Data: err.Error()})
		return
	}

	// determine the offset and limit for pagination stage & delete "offset" fields in professorQuery
	paginateMap, err := configs.GetAggregateLimit(&courseQuery, c)

	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ErrorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// pipeline to query the sections from the filtered courses
	courseSectionPipeline := mongo.Pipeline{
		// filter the courses
		bson.D{{Key: "$match", Value: courseQuery}},

		// paginate the courses before pulling the sections from thoses courses
		bson.D{{Key: "$skip", Value: paginateMap["former_offset"]}}, // skip to the specified offset
		bson.D{{Key: "$limit", Value: paginateMap["limit"]}},        // limit to the specified number of courses

		// lookup the sections of the courses
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},

		// unwind the sections of the courses
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$sections"},
			{Key: "preserveNullAndEmptyArrays", Value: false}, // avoid course documents that can't be replaced
		}}},

		// replace the courses with sections
		bson.D{{Key: "$replaceWith", Value: "$sections"}},

		// keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},

		// paginate the sections
		bson.D{{Key: "$skip", Value: paginateMap["latter_offset"]}},
		bson.D{{Key: "$limit", Value: paginateMap["limit"]}},
	}

	// perform aggregation on the pipeline
	cursor, err := courseCollection.Aggregate(ctx, courseSectionPipeline)
	if err != nil {
		// return error for any aggregation problem
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "aggregation error", Data: err.Error()})
		return
	}
	// parse the array of sections of the course
	if err = cursor.All(ctx, &courseSections); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, responses.MultiSectionResponse{Status: http.StatusOK, Message: "success", Data: courseSections})
}
