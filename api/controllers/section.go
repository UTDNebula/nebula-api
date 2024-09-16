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

// @Id sectionSearch
// @Router /section [get]
// @Description "Returns all courses matching the query's string-typed key-value pairs"
// @Produce json
// @Param section_number query string false "The section's official number"
// @Param course_reference query string false "An id that points to the course in MongoDB that this section is an instantiation of"
// @Param academic_session.name query string false "The name of the academic session of the section"
// @Param academic_session.start_date query string false "The date of classes starting for the section"
// @Param academic_session.end_date query string false "The date of classes ending for the section"
// @Param professors query string false "One of the professors teaching the section"
// @Param teaching_assistants.first_name query string false "The first name of one of the teaching assistants of the section"
// @Param teaching_assistants.last_name query string false "The last name of one of the teaching assistants of the section"
// @Param teaching_assistants.role query string false "The role of one of the teaching assistants of the section"
// @Param teaching_assistants.email query string false "The email of one of the teaching assistants of the section"
// @Param internal_class_number query string false "The internal (university) number used to reference this section"
// @Param instruction_mode query string false "The instruction modality for this section"
// @Param meetings.start_date query string false "The start date of one of the section's meetings"
// @Param meetings.end_date query string false "The end date of one of the section's meetings"
// @Param meetings.meeting_days query string false "One of the days that one of the section's meetings"
// @Param meetings.start_time query string false "The time one of the section's meetings starts"
// @Param meetings.end_time query string false "The time one of the section's meetings ends"
// @Param meetings.modality query string false "The modality of one of the section's meetings"
// @Param meetings.location.building query string false "The building of one of the section's meetings"
// @Param meetings.location.room query string false "The room of one of the section's meetings"
// @Param meetings.location.map_uri query string false "A hyperlink to the UTD room locator of one of the section's meetings"
// @Param core_flags query string false "One of core requirement codes this section fulfills"
// @Param syllabus_uri query string false "A link to the syllabus on the web"
// @Success 200 {array} schema.Section "A list of sections"
func SectionSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

// @Id sectionById
// @Router /section/{id} [get]
// @Description "Returns the section with given ID"
// @Produce json
// @Param id path string true "ID of the section to get"
// @Success 200 {object} schema.Section "A section"
func SectionById() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}
