package controllers

import (
	"context"
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
// @Description "Returns all courses matching the query's string-typed key-value pairs"
// @Produce json
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
func CourseSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

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
		c.JSON(http.StatusOK, responses.MultiCourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
	}
}

// @Id courseById
// @Router /course/{id} [get]
// @Description "Returns the course with given ID"
// @Produce json
// @Param id path string true "ID of the course to get"
// @Success 200 {object} schema.Course "A course"
func CourseById() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func CourseAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

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
}
