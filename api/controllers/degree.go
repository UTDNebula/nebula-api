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

var degreeCollection *mongo.Collection = configs.GetCollection(configs.DB, "degrees")

func DegreeSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var degrees []schema.Degree

	query, err := schema.FilterQuery[schema.Degree](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.DegreeResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.DegreeResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	cursor, err := degreeCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.DegreeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err = cursor.All(ctx, &degrees); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: degrees})
}

func DegreeById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	degreeId := c.Param("id")

	var degree schema.Degree

	objId, err := primitive.ObjectIDFromHex(degreeId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	err = degreeCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&degree)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.DegreeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: degree})
}
