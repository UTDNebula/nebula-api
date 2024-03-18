package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GradesAggregation handles the aggregation of grades based on various filters.
func GradesAggregation(c *gin.Context) {
	var grades []map[string]interface{}
	var results []map[string]interface{}

	var cursor *mongo.Cursor
	var collection *mongo.Collection
	var pipeline mongo.Pipeline

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract query parameters
	prefix := c.Query("prefix")
	number := c.Query("number")
	sectionNumber := c.Query("section_number")
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")

	professor := (firstName != "" || lastName != "")

	// MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.GradeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Database and collection
	database := client.Database("your_database")
	courseCollection := database.Collection("courses")
	sectionCollection := database.Collection("sections")
	professorCollection := database.Collection("professors")

	// Stages for MongoDB aggregation pipeline
	lookupSectionsStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}},
	}

	unwindSectionsStage := bson.D{{Key: "$unwind", Value: "$sections"}}

	projectGradeDistributionStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: "$sections.academic_session.name"},
			{Key: "grade_distribution", Value: "$sections.grade_distribution"},
		}},
	}

	unwindGradeDistributionStage := bson.D{{Key: "$unwind", Value: "$grade_distribution"}}

	groupGradesStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "academic_session", Value: "$_id"},
				{Key: "ix", Value: "$ix"},
			}},
			{Key: "grades", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
		}},
	}

	sortGradesStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "_id.ix", Value: 1},
			{Key: "_id", Value: 1},
		}},
	}

	sumGradesStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "grades", Value: bson.D{{Key: "$sum", Value: "$grades"}}}}}}

	groupGradeDistributionStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.academic_session"},
			{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grades"}}},
		}},
	}

	switch {
	case prefix != "" && number == "" && sectionNumber == "" && !professor:
		// Filter on Course
		collection = courseCollection
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix}}})

	case prefix != "" && number != "" && sectionNumber == "" && !professor:
		// Filter on Course
		collection = courseCollection
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix, "course_number": number}}})

	case prefix != "" && number != "" && sectionNumber != "" && !professor:
		// Filter on Course then Section
		collection = courseCollection
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix, "course_number": number}}})
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"sections.section_number": sectionNumber}}})

	case prefix == "" && number == "" && sectionNumber == "" && professor:
		// Filter on Professor
		collection = professorCollection
		match := bson.M{}
		if firstName != "" {
			match["first_name"] = firstName
		}
		if lastName != "" {
			match["last_name"] = lastName
		}
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: match}})

	case prefix != "" && professor:
		// Filter on Section by Matching Course and Professor IDs
		collection = sectionCollection

		var profIDs []primitive.ObjectID
		match := bson.M{"professors": bson.M{}}
		if firstName != "" {
			match["first_name"] = firstName
		}
		if lastName != "" {
			match["last_name"] = lastName
		}
		cursor, err = professorCollection.Find(ctx, bson.D{{Key: "$match", Value: match}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GradeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &results); err != nil {
			panic(err)
		}

		for _, prof := range results {
			profID := prof["_id"].(primitive.ObjectID)
			profIDs = append(profIDs, profID)
		}

		courseMatch := bson.M{"subject_prefix": prefix}
		if number != "" {
			courseMatch["course_number"] = number
		}
		sectionMatch := bson.M{"course_reference": bson.M{"$in": profIDs}}
		if sectionNumber != "" {
			sectionMatch["section_number"] = sectionNumber
		}

		pipeline = append(pipeline, bson.D{{Key: "$match", Value: sectionMatch}})

	default:
		c.JSON(http.StatusBadRequest, responses.GradeResponse{Status: http.StatusBadRequest, Message: "error", Data: "Invalid query parameters."})
		return
	}

	// perform aggregation
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &grades); err != nil {
		panic(err)
	}

	// Handle different flags
	flag := c.Query("flag")
	if flag == "overall" {
		// combine all semester grade_distributions
		overallResponse := make([]int32, 14)
		for _, sem := range grades {
			if len(sem["grade_distribution"].([]int32)) != 14 {
				print("Length of Array: ")
				println(len(sem["grade_distribution"].([]int32)))
			}
			for i, grade := range sem["grade_distribution"].([]int32) {
				overallResponse[i] += grade
			}
		}
		c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: overallResponse})
	} else if flag == "semester" {
		c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: grades})
	} else {
		c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Endpoint broken"})
	}
}
