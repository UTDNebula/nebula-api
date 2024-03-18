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

var professorCollection *mongo.Collection

func init() {
	professorCollection = configs.GetCollection(configs.DB, "professors")
}

func ProfessorSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var professors []schema.Professor

		query, err := schema.FilterQuery[schema.Professor](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.ProfessorResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
			return
		}

		optionLimit, err := configs.GetOptionLimit(&query, c)
		if err != nil {
			log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
			c.JSON(http.StatusConflict, responses.ProfessorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
			return
		}

		cursor, err := professorCollection.Find(ctx, query, optionLimit)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if err = cursor.All(ctx, &professors); err != nil {
			log.WritePanic(err)
			panic(err)
		}

		c.JSON(http.StatusOK, responses.ProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})
	}
}

func ProfessorById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		professorId := c.Param("id")
		var professor schema.Professor

		objId, err := primitive.ObjectIDFromHex(professorId)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		err = professorCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&professor)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.ProfessorResponse{Status: http.StatusOK, Message: "success", Data: professor})
	}
}

func ProfessorAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var professors []schema.Professor

		cursor, err := professorCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if err = cursor.All(ctx, &professors); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, responses.ProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})
	}
}
