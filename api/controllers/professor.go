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

var professorCollection *mongo.Collection = configs.GetCollection(configs.DB, "professors")

func ProfessorSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var professors []schema.Professor

		defer cancel()

		// build query key value pairs (only one value per key)
		query := bson.M{}
		for key := range queryParams {
			query[key] = c.Query(key)
		}

		optionLimit, err := configs.GetOptionLimit(&query, c)
		if err != nil {
			log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
			c.JSON(http.StatusConflict, responses.ProfessorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
			return
		}

		// get cursor for query results
		cursor, err := professorCollection.Find(ctx, query, optionLimit)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &professors); err != nil {
			log.WritePanic(err)
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.ProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})
	}
}

func ProfessorById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		professorId := c.Param("id")

		var professor schema.Professor

		defer cancel()

		// parse object id from id parameter
		objId, err := primitive.ObjectIDFromHex(professorId)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		// find and parse matching professor
		err = professorCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&professor)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// return result
		c.JSON(http.StatusOK, responses.ProfessorResponse{Status: http.StatusOK, Message: "success", Data: professor})
	}
}

func ProfessorAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var professors []schema.Professor

		defer cancel()

		cursor, err := professorCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &professors); err != nil {
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.ProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})

	}
}
