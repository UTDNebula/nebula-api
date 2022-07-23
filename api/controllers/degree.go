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

var degreeCollection *mongo.Collection = configs.GetCollection(configs.DB, "degrees")

func DegreeSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// @TODO: Fix with model - There is NO typechecking!
		// var degrees []models.Degree
		var degrees []map[string]interface{}

		defer cancel()

		// build query key value pairs (only one value per key)
		query := bson.M{}
		for key, _ := range queryParams {
			query[key] = c.Query(key)
		}

		// get cursor for query results
		cursor, err := degreeCollection.Find(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DegreeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &degrees); err != nil {
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: degrees})
	}
}

func DegreeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		degreeId := c.Param("id")

		// @TODO: Fix with model - There is NO typechecking!
		// var degree models.Degree
		var degree map[string]interface{}

		defer cancel()

		// parse object id from id parameter
		objId, _ := primitive.ObjectIDFromHex(degreeId)

		// find and parse matching degree
		err := degreeCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&degree)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DegreeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// return result
		c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: degree})
	}
}
