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

var examCollection *mongo.Collection = configs.GetCollection(configs.DB, "exams")

type examFilter struct {
	Type  string `schema:"type"`
	Name  string `schema:"name"`
	Level string `schema:"level"`
}

func ExamSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exams []map[string]interface{}

	query, err := schema.FilterQuery[examFilter](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ExamResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ExamResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	cursor, err := examCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ExamResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err = cursor.All(ctx, &exams); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	c.JSON(http.StatusOK, responses.ExamResponse{Status: http.StatusOK, Message: "success", Data: exams})
}

func ExamById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	examId := c.Param("id")

	var exam map[string]interface{}

	objId, err := primitive.ObjectIDFromHex(examId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.ExamResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	err = examCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&exam)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ExamResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.ExamResponse{Status: http.StatusOK, Message: "success", Data: exam})
}

func ExamAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exams []map[string]interface{}

	cursor, err := examCollection.Find(ctx, bson.M{})
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ExamResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	err = cursor.All(ctx, &exams)
	if err != nil {
		log.WritePanic(err)
		panic(err)
	}

	c.JSON(http.StatusOK, responses.ExamResponse{Status: http.StatusOK, Message: "success", Data: exams})
}
