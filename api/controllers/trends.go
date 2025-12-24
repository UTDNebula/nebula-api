package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Id				trendsCourseSectionSearch
// @Router			/course/sections/trends [get]
// @Tags			Courses
// @Description	"Returns all of the given course's sections with Course and Professor data embedded. Specialized high-speed convenience endpoint for UTD Trends internal use; limited query flexibility."
// @Produce		json
// @Param			course_number	query		string									true	"The course's official number"
// @Param			subject_prefix	query		string									true	"The course's subject prefix"
// @Success		200				{object}	schema.APIResponse[[]schema.Section]	"A list of Sections"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
func TrendsCourseSectionSearch(c *gin.Context) {
	trendsSectionSearch("Course", c)
}

// @Id				trendsProfessorSectionSearch
// @Router			/professor/sections/trends [get]
// @Tags			Professors
// @Description	"Returns all of the given professor's sections with Course and Professor data embedded. Specialized high-speed convenience endpoint for UTD Trends internal use; limited query flexibility."
// @Produce		json
// @Param			first_name	query		string									true	"The professor's first name"
// @Param			last_name	query		string									true	"The professor's last name"
// @Success		200			{object}	schema.APIResponse[[]schema.Section]	"A list of Sections"
// @Failure		500			{object}	schema.APIResponse[string]				"A string describing the error"
func TrendsProfessorSectionSearch(c *gin.Context) {
	trendsSectionSearch("Professor", c)
}

// trendsSectionSearch handles trends-based section routes for both course and professor query.
//
// This is to reduce the repetitiveness of routes whose aggregation behaviors are basically similar.
// This is subject to change as requests may be more complex in the future.
func trendsSectionSearch(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var detailedSections []schema.Section

	// If the branches become complicated, refactor to make it more generic.
	var trendsCollection *mongo.Collection
	var trendsQuery bson.M
	var err error

	switch flag {
	case "Course":
		trendsCollection = configs.GetCollection("trends_course_sections")
		// UTDTrends only uses prefix and number
		trendsQuery = bson.M{
			"_id": c.Query("subject_prefix") + c.Query("course_number"),
		}
	case "Professor":
		trendsCollection = configs.GetCollection("trends_prof_sections")
		trendsQuery, err = schema.FilterQuery[schema.Professor](c)
		if err != nil {
			return
		}
	default:
		// This should never happen, but act as a fallback
		err = fmt.Errorf("invalid flag for trendsSectionSearch: %s", flag)
		respondWithInternalError(c, err)
	}

	trendsPipeline := buildTrendsPipeline(trendsQuery)

	cursor, err := trendsCollection.Aggregate(ctx, trendsPipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &detailedSections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", detailedSections)
}

// buildTrendsPipeline build the pipeline to embed the list of professors and course details
// directly in the queried list of sections
func buildTrendsPipeline(srcObjQuery bson.M) mongo.Pipeline {
	return mongo.Pipeline{
		// Match the original objects from the query
		bson.D{{Key: "$match", Value: srcObjQuery}},

		// Expand sections array into individual documents
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$sections"},
				{Key: "preserveNullAndEmptyArrays", Value: false}, // Avoid documents that can't be replaced
			}},
		},

		// Embed the course details to the list of sections
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "sections.course_reference"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.course_details"},
			}},
		},

		// Embed the professor details to the list of sections
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "sections.professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.professor_details"},
			}},
		},

		// Replace the original objects with sections
		bson.D{{Key: "$replaceWith", Value: "$sections"}},

		// Keep order deterministic between calls
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}},
		},
	}
}
