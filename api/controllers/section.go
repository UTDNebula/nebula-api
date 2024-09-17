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

var sectionCollection *mongo.Collection = configs.GetCollection("sections")

func SectionSearch(c *gin.Context) {
	//name := c.Query("name")            // value of specific query parameter: string
	//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var sections []schema.Section

	defer cancel()

	// build query key value pairs (only one value per key)
	query, err := schema.FilterQuery[schema.Section](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	if v, ok := query["course_reference"]; ok {
		objId, err := primitive.ObjectIDFromHex(v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		} else {
			query["course_reference"] = objId
		}
	}

	if v, ok := query["professor"]; ok {
		objId, err := primitive.ObjectIDFromHex(v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		} else {
			query["professor"] = objId
		}
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ErrorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// get cursor for query results
	cursor, err := sectionCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &sections); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	// return result
	c.JSON(http.StatusOK, responses.MultiSectionResponse{Status: http.StatusOK, Message: "success", Data: sections})
}

func SectionById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	sectionId := c.Param("id")

	var section schema.Section

	defer cancel()

	// parse object id from id parameter
	objId, err := primitive.ObjectIDFromHex(sectionId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	// find and parse matching section
	err = sectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&section)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// return result
	c.JSON(http.StatusOK, responses.SingleSectionResponse{Status: http.StatusOK, Message: "success", Data: section})
}
