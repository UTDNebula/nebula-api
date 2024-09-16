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

var degreeCollection *mongo.Collection = configs.GetCollection("degrees")

func DegreeSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var degrees []schema.Degree

		defer cancel()

		// build query key value pairs (only one value per key)
		query, err := schema.FilterQuery[schema.Degree](c)
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
		cursor, err := degreeCollection.Find(ctx, query, optionLimit)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &degrees); err != nil {
			log.WritePanic(err)
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.MultiDegreeResponse{Status: http.StatusOK, Message: "success", Data: degrees})
	}
}

func DegreeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		degreeId := c.Param("id")

		var degree schema.Degree

		defer cancel()

		// parse object id from id parameter
		objId, err := primitive.ObjectIDFromHex(degreeId)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		// find and parse matching degree
		err = degreeCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&degree)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// return result
		c.JSON(http.StatusOK, responses.SingleDegreeResponse{Status: http.StatusOK, Message: "success", Data: degree})
	}
}
