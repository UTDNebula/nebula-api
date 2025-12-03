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

// CourseSectionV2 does the exact same thing as CourseSection, but it doesn't
// do use any pipline. This could potentially enhance performance for large DB like ours
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
		respond(c, http.StatusBadRequest, "Error offset is not of type integer", err.Error())
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

// CourseProfessorV2 also does the exact same thing as CourseProfessor, but doesn't
// do use any pipline
func CourseProfessorV2(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List of professors of the filtered courses
	var courseProfessors []schema.Professor

	// Build the course query
	courseQuery, err := getQuery[schema.Course](flag, c)
	if err != nil {
		return
	}
	paginate, err := configs.GetAggregateLimit(&courseQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not of type integer", err.Error())
		return
	}

	results := make([]map[string]any, 0) // Auxiliary to parse the intermediary

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
	if err = cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}
	// The auxiliary array of IDs used to for fetching sections and profs
	auxIDs := make([]primitive.ObjectID, 0, len(results))
	for _, course := range results {
		auxIDs = append(auxIDs, course["_id"].(primitive.ObjectID))
	}

	// Get the list of section IDs
	cursor, err = sectionCollection.Find(ctx, bson.M{"course_reference": bson.M{"$in": auxIDs}})
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if err = cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}
	auxIDs = make([]primitive.ObjectID, 0, len(results))
	for _, section := range results {
		auxIDs = append(auxIDs, section["_id"].(primitive.ObjectID))
	}

	// Get the list of professors from the given list of section IDs
	cursor, err = professorCollection.Find(
		ctx,
		bson.M{"sections": bson.M{"$in": auxIDs}},
		options.Find().SetSkip(paginate["latter_offset"]).SetLimit(paginate["limit"]),
	)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if err = cursor.All(ctx, &courseProfessors); err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "sucess", courseProfessors)
}
