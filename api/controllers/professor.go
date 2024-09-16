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

// @Id professorSearch
// @Router /professor [get]
// @Description "Returns all professors matching the query's string-typed key-value pairs"
// @Produce json
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professor's last name"
// @Param titles query string false "One of the professor's title"
// @Param email query string false "The professor's email address"
// @Param phone_number query string false "The professor's phone number"
// @Param office.building query string false "The building of the location of the professor's office"
// @Param office.room query string false "The room of the location of the professor's office"
// @Param office.map_uri query string false "A hyperlink to the UTD room locator of the professor's office"
// @Param profile_uri query string false "A hyperlink pointing to the professor's official university profile"
// @Param image_uri query string false "A link to the image used for the professor on the professor's official university profile"
// @Param office_hours.start_date query string false "The start date of one of the office hours meetings of the professor"
// @Param office_hours.end_date query string false "The end date of one of the office hours meetings of the professor"
// @Param office_hours.meeting_days query string false "One of the days that one of the office hours meetings of the professor"
// @Param office_hours.start_time query string false "The time one of the office hours meetings of the professor starts"
// @Param office_hours.end_time query string false "The time one of the office hours meetings of the professor ends"
// @Param office_hours.modality query string false "The modality of one of the office hours meetings of the professor"
// @Param office_hours.location.building query string false "The building of one of the office hours meetings of the professor"
// @Param office_hours.location.room query string false "The room of one of the office hours meetings of the professor"
// @Param office_hours.location.map_uri query string false "A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Param sections query string false "The _id of one of the sections the professor teaches"
// @Success 200 {array} schema.Professor "A list of professors"
func ProfessorSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.Query("name")            // value of specific query parameter: string
		//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var professors []schema.Professor

		defer cancel()

		// build query key value pairs (only one value per key)
		query, err := schema.FilterQuery[schema.Professor](c)
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
		cursor, err := professorCollection.Find(ctx, query, optionLimit)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &professors); err != nil {
			log.WritePanic(err)
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.MultiProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})
	}
}

// @Id professorById
// @Router /professor/{id} [get]
// @Description "Returns the professor with given ID"
// @Produce json
// @Param id path string true "ID of the professor to get"
// @Success 200 {object} schema.Professor "A professor"
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
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		// find and parse matching professor
		err = professorCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&professor)
		if err != nil {
			log.WriteError(err)
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// return result
		c.JSON(http.StatusOK, responses.SingleProfessorResponse{Status: http.StatusOK, Message: "success", Data: professor})
	}
}

func ProfessorAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var professors []schema.Professor

		defer cancel()

		cursor, err := professorCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &professors); err != nil {
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.MultiProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})

	}
}
