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
)

// We want to Filter (Match) ASAP

// --------------------------------------------------------

// Aggregate By Course -> Section:
// ---- Prefix
// ---- Prefix, Number
// ---- Prefix, Number, SectionNumber

// Aggregate By Professor -> Section
// ---- Professor

// Aggregate By Find Course, Find Professor: then Match Section
// ---- Prefix, Professor
// ---- Prefix, Number, Professor
// ---- Prefix, Number, Professor, SectionNumber

// --------------------------------------------------------

// Filter on Course
// ---- Prefix
// ---- Prefix, Number

// Filter on Course then Section
// ---- Prefix, Number, SectionNumber

// Filter on Professor
// ---- Professor

// Filter on Section by Matching Course and Professor IDs
// ---- Prefix, Professor
// ---- Prefix, Number, Professor
// ---- Prefix, Number, Professor, SectionNumber

// 4 Functions

func GradesAggregation(flag string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var grades []map[string]interface{}
		var results []map[string]interface{}

		var cursor *mongo.Cursor
		var collection *mongo.Collection
		var pipeline mongo.Pipeline

		var sectionMatch bson.D
		var courseMatch bson.D
		var courseFind bson.D
		var professorMatch bson.D
		var professorFind bson.D

		var err error

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// @TODO: Recommend forcing using first_name and last_name to ensure single professors per query.
		// All professors sharing the name will be aggregated together in the current implementation
		prefix := c.Query("prefix")
		number := c.Query("number")
		section_number := c.Query("section_number")
		first_name := c.Query("first_name")
		last_name := c.Query("last_name")

		professor := (first_name != "" || last_name != "")

		lookupSectionsStage := bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sections"},
				{Key: "localField", Value: "sections"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections"},
			}},
		}

		unwindSectionsStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$sections"}}}}

		projectGradeDistributionStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: "$sections.academic_session.name"},
				{Key: "grade_distribution", Value: "$sections.grade_distribution"},
			}},
		}

		projectGradeDistributionWithSectionsStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: "$academic_session.name"},
				{Key: "grade_distribution", Value: "$grade_distribution"},
			}},
		}

		unwindGradeDistributionStage := bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$grade_distribution"},
				{Key: "includeArrayIndex", Value: "ix"},
			}},
		}

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
		case prefix != "" && number == "" && section_number == "" && !professor:
			// Filter on Course
			collection = courseCollection
			courseMatch = bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix}}}
			pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

		case prefix != "" && number != "" && section_number == "" && !professor:
			// Filter on Course
			collection = courseCollection
			courseMatch := bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix, "course_number": number}}}
			pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

		case prefix != "" && number != "" && section_number != "" && !professor:
			// Filter on Course then Section
			collection = courseCollection
			courseMatch := bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix, "course_number": number}}}
			sectionMatch := bson.D{{Key: "$match", Value: bson.M{"sections.section_number": section_number}}}
			pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage, unwindSectionsStage, sectionMatch, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

		case prefix == "" && number == "" && section_number == "" && professor:
			// Filter on Professor
			collection = professorCollection

			// Build professorMatch
			if last_name == "" {
				professorMatch = bson.D{{Key: "$match", Value: bson.M{"first_name": first_name}}}
			} else if first_name == "" {
				professorMatch = bson.D{{Key: "$match", Value: bson.M{"last_name": last_name}}}
			} else {
				professorMatch = bson.D{{Key: "$match", Value: bson.M{"first_name": first_name, "last_name": last_name}}}
			}

			// Build grades pipeline
			pipeline = mongo.Pipeline{professorMatch, lookupSectionsStage, unwindSectionsStage, projectGradeDistributionStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

		case prefix != "" && professor:
			// Filter on Section by Matching Course and Professor IDs

			// Here we get the valid course ids and professor ids
			// and then we perform the grades aggregation against the sections collection,
			// matching on the course_reference and professor

			var profIDs []primitive.ObjectID
			var courseIDs []primitive.ObjectID

			collection = sectionCollection

			// Find valid professor ids
			if last_name == "" {
				professorFind = bson.D{{Key: "first_name", Value: first_name}}
			} else if first_name == "" {
				professorFind = bson.D{{Key: "last_name", Value: last_name}}
			} else {
				professorFind = bson.D{{Key: "first_name", Value: first_name}, {Key: "last_name", Value: last_name}}
			}

			cursor, err = professorCollection.Find(ctx, professorFind)
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

			// Get valid course ids
			if number == "" {
				courseFind = bson.D{{Key: "subject_prefix", Value: prefix}}
			} else {
				courseFind = bson.D{{Key: "subject_prefix", Value: prefix}, {Key: "course_number", Value: number}}
			}

			cursor, err = courseCollection.Find(ctx, courseFind)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GradeResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
				return
			}
			if err = cursor.All(ctx, &results); err != nil {
				panic(err)
			}

			for _, course := range results {
				courseID := course["_id"].(primitive.ObjectID)
				courseIDs = append(courseIDs, courseID)
			}

			// Build sectionMatch
			if section_number == "" {
				sectionMatch =
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "course_reference", Value: bson.D{{Key: "$in", Value: courseIDs}}},
						{Key: "professors", Value: bson.D{{Key: "$in", Value: profIDs}}},
					}}}
			} else {
				sectionMatch =
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "course_reference", Value: bson.D{{Key: "$in", Value: courseIDs}}},
						{Key: "professors", Value: bson.D{{Key: "$in", Value: profIDs}}},
						{Key: "section_number", Value: section_number},
					}}}
			}

			// Build grades pipeline
			pipeline = mongo.Pipeline{sectionMatch, projectGradeDistributionWithSectionsStage, unwindGradeDistributionStage, groupGradesStage, sortGradesStage, sumGradesStage, groupGradeDistributionStage}

		default:
			c.JSON(http.StatusBadRequest, responses.GradeResponse{Status: http.StatusBadRequest, Message: "error", Data: "Invalid query parameters."})
			return
		}

		// peform aggregation
		cursor, err = collection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &grades); err != nil {
			panic(err)
		}

		if flag == "overall" {
			// combine all semester grade_distributions
			overallResponse := make([]int32, 14)
			for _, sem := range grades {
				if len(sem["grade_distribution"].(primitive.A)) != 14 {
					print("Length of Array: ")
					println(len(sem["grade_distribution"].(primitive.A)))
				}
				for i, grade := range sem["grade_distribution"].(primitive.A) {
					overallResponse[i] += grade.(int32)
				}
			}
			c.JSON(http.StatusOK, responses.GradeResponse{Status: http.StatusOK, Message: "success", Data: overallResponse})
		} else if flag == "semester" {
			c.JSON(http.StatusOK, responses.GradeResponse{Status: http.StatusOK, Message: "success", Data: grades})
		} else {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Endpoint broken"})
		}

	}
}
