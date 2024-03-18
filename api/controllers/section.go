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

var sectionCollection *mongo.Collection = configs.GetCollection(configs.DB, "sections")

func SectionSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sections []schema.Section

	query, err := schema.FilterQuery[schema.Section](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.SectionResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	if v, ok := query["course_reference"]; ok {
		objId, err := primitive.ObjectIDFromHex(v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		} else {
			query["course_reference"] = objId
		}
	}

	if v, ok := query["professor"]; ok {
		objId, err := primitive.ObjectIDFromHex(v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		} else {
			query["professor"] = objId
		}
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.SectionResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	cursor, err := sectionCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.SectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err = cursor.All(ctx, &sections); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	c.JSON(http.StatusOK, responses.SectionResponse{Status: http.StatusOK, Message: "success", Data: sections})
}

func SectionById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sectionId := c.Param("id")

	var section schema.Section

	objId, err := primitive.ObjectIDFromHex(sectionId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	err = sectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&section)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.SectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.SectionResponse{Status: http.StatusOK, Message: "success", Data: section})
}
