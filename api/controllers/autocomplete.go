package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/responses"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Modified code
func AutocompleteDAG(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var autocompleteDAG []map[string]interface{}

	autocompletePipeline := mongo.Pipeline{
		{{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "sections"},
				{Key: "localField", Value: "sections"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "section"},
			},
		}},
		{{Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$section"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		}},
		{{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "section.professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "professor"},
			},
		}},
		{{Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$professor"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		}},
		{{Key: "$project",
			Value: bson.D{
				{Key: "subject_prefix", Value: "$subject_prefix"},
				{Key: "course_number", Value: "$course_number"},
				{Key: "academic_session.name", Value: "$section.academic_session.name"},
				{Key: "section_number", Value: "$section.section_number"},
				{Key: "professor",
					Value: bson.D{
						{Key: "first_name", Value: "$professor.first_name"},
						{Key: "last_name", Value: "$professor.last_name"},
					},
				},
			},
		}},
		{{Key: "$group",
			Value: bson.D{
				{Key: "_id",
					Value: bson.D{
						{Key: "subject_prefix", Value: "$subject_prefix"},
						{Key: "course_number", Value: "$course_number"},
						{Key: "academic_session", Value: "$academic_session"},
						{Key: "section_number", Value: "$section_number"},
					},
				},
				{Key: "professor",
					Value: bson.D{
						{Key: "$push", Value: "$professor"},
					},
				},
			},
		}},
		{{Key: "$group",
			Value: bson.D{
				{Key: "_id",
					Value: bson.D{
						{Key: "subject_prefix", Value: "$_id.subject_prefix"},
						{Key: "course_number", Value: "$_id.course_number"},
						{Key: "academic_session", Value: "$_id.academic_session"},
					},
				},
				{Key: "sections",
					Value: bson.D{
						{Key: "$push",
							Value: bson.D{
								{Key: "section_number", Value: "$_id.section_number"},
								{Key: "professors", Value: "$professor"},
							},
						},
					},
				},
			},
		}},
		{{Key: "$group",
			Value: bson.D{
				{Key: "_id",
					Value: bson.D{
						{Key: "subject_prefix", Value: "$_id.subject_prefix"},
						{Key: "course_number", Value: "$_id.course_number"},
					},
				},
				{Key: "academic_sessions",
					Value: bson.D{
						{Key: "$push",
							Value: bson.D{
								{Key: "academic_session", Value: "$_id.academic_session"},
								{Key: "sections", Value: "$sections"},
							},
						},
					},
				},
			},
		}},
		{{Key: "$group",
			Value: bson.D{
				{Key: "_id",
					Value: bson.D{
						{Key: "subject_prefix", Value: "$_id.subject_prefix"},
					},
				},
				{Key: "course_numbers",
					Value: bson.D{
						{Key: "$push",
							Value: bson.D{
								{Key: "course_number", Value: "$_id.course_number"},
								{Key: "academic_sessions", Value: "$academic_sessions"},
							},
						},
					},
				},
			},
		}},
		{{Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "subject_prefix", Value: "$_id.subject_prefix"},
				{Key: "course_numbers", Value: "$course_numbers"},
			},
		}},
	}

	cursor, err := courseCollection.Aggregate(ctx, autocompletePipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.AutocompleteResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err = cursor.All(ctx, &autocompleteDAG); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, responses.AutocompleteResponse{Status: http.StatusOK, Message: "success", Data: autocompleteDAG})
}
