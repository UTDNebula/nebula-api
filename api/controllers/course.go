package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var courseCollection *mongo.Collection = configs.GetCollection(configs.DB, "courses")

func CourseSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// @TODO: Fix with model - There is NO typechecking!
		// var courses []models.Course
		var courses []map[string]interface{}

		// build query key value pairs (only one value per key)
		query := bson.M{}
		for key, _ := range queryParams {
			query[key] = c.Query(key)
		}

		optionLimit, err := configs.GetOptionLimit(&query, c); if err != nil {
			c.JSON(http.StatusConflict, responses.CourseResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
			return
		}

		// get cursor for query results
		cursor, err := courseCollection.Find(ctx, query, optionLimit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &courses); err != nil {
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
	}
}

func CourseById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		courseId := c.Param("id")

		// @TODO: Fix with model - There is NO typechecking!
		// var course models.Course
		var course map[string]interface{}

		// parse object id from id parameter
		objId, err := primitive.ObjectIDFromHex(courseId)
		if err != nil{
			c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		// find and parse matching course
		err = courseCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&course)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// return result
		c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: course})
	}
}
