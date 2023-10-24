package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var examCollection *mongo.Collection = configs.GetCollection(configs.DB, "exams")

func ExamSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// @TODO: Fix with model - There is NO typechecking!
		// var exams []models.Exam
		var exams []map[string]interface{}

		defer cancel()

		// build query key value pairs (only one value per key)
		query := bson.M{}
		for key := range queryParams {
			query[key] = c.Query(key)
		}

		optionLimit, err := configs.GetOptionLimit(&query, c)
		if err != nil {
			log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
			c.JSON(http.StatusConflict, responses.ExamResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
			return
		}

		// get cursor for query results
		cursor, err := examCollection.Find(ctx, query, optionLimit)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ExamResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &exams); err != nil {
			log.WritePanic(err)
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.ExamResponse{Status: http.StatusOK, Message: "success", Data: exams})
	}
}

func ExamById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		examId := c.Param("id")

		// @TODO: Fix with model - There is NO typechecking!
		// var exam models.Exam
		var exam map[string]interface{}

		defer cancel()

		// parse object id from id parameter
		objId, err := primitive.ObjectIDFromHex(examId)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusBadRequest, responses.ExamResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		// find and parse matching exam
		err = examCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&exam)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ExamResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// return result
		c.JSON(http.StatusOK, responses.ExamResponse{Status: http.StatusOK, Message: "success", Data: exam})
	}
}

func ExamAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// @TODO: Fix with model - There is NO typechecking!
		var exams []map[string]interface{}

		defer cancel()

		// get cursor for all exams in the collection
		cursor, err := examCollection.Find(ctx, bson.M{})
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ExamResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		err = cursor.All(ctx, &exams)
		if err != nil {
			log.WritePanic(err)
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.ExamResponse{Status: http.StatusOK, Message: "success", Data: exams})
	}
}
