package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CourseSectionV2(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List of sections of the filtered courses
	var courseSections []schema.Section

	// Build the course query
	courseQuery, err := getQuery[schema.Course](flag, c)
	if err != nil {
		return
	}

	// Get the offset and limit for pagination
	paginate, err := configs.GetAggregateLimit(&courseQuery, c)
	if err != nil {
		message := "Error offset is not of type integer"
		respond(c, http.StatusBadRequest, message, err.Error())
		return
	}

	// Get the list of course IDs
	cursor, err := courseCollection.Find(
		ctx,
		courseQuery,
		options.Find().SetSkip(paginate["former_offset"]).SetLimit(paginate["limit"]),
	)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	var courses []map[string]any
	if err = cursor.All(ctx, &courses); err != nil {
		respondWithInternalError(c, err)
		return
	}
	courseIDs := make([]primitive.ObjectID, 0, len(courses))
	for _, course := range courses {
		courseIDs = append(courseIDs, course["_id"].(primitive.ObjectID))
	}

	// Get the list of sections from the given list of course IDs
	cursor, err = sectionCollection.Find(
		ctx,
		bson.M{"course_reference": bson.M{"$in": courseIDs}},
		options.Find().SetSkip(paginate["latter_offset"]).SetLimit(paginate["limit"]),
	)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if err = cursor.All(ctx, &courseSections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "sucess", courseSections)
}
