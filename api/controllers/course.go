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

var courseCollection *mongo.Collection = configs.GetCollection(configs.DB, "courses")

func CourseSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courses []schema.Course

	query, err := schema.FilterQuery[schema.Course](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.CourseResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	cursor, err := courseCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err = cursor.All(ctx, &courses); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
}

func CourseById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := c.Param("id")

	var course schema.Course

	objId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	err = courseCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&course)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: course})
}

func CourseAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courses []schema.Course

	cursor, err := courseCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err = cursor.All(ctx, &courses); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
}
