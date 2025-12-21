package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var trendsCourseCollection *mongo.Collection = configs.GetCollection("trends_course_sections")
var trendsProfCollection *mongo.Collection = configs.GetCollection("trends_prof_sections")

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var courseSections []schema.Section
	courseQuery := bson.M{"_id": c.Query("subject_prefix") + c.Query("course_number")}
	var err error

	// Pipeline to query the Sections + Professors from the filtered courses
	pipeline := mongo.Pipeline{
		// filter the courses
		bson.D{{Key: "$match", Value: courseQuery}},

		// unwind the sections
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$sections"},
				{Key: "preserveNullAndEmptyArrays", Value: false}, // avoid course documents that can't be replaced
			}},
		},

		// lookup the professors of the sections
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "sections.professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.professor_details"},
			}},
		},

		// lookup the course of the sections
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "sections.course_reference"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.course_details"},
			}},
		},

		// replace the courses with sections
		bson.D{{Key: "$replaceWith", Value: "$sections"}},

		// keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	// perform aggregation on the pipeline
	cursor, err := trendsCourseCollection.Aggregate(ctx, pipeline)
	if err != nil {
		// return error for any aggregation problem
		respondWithInternalError(c, err)
		return
	}
	// parse the array of sections of the course
	if err = cursor.All(ctx, &courseSections); err != nil {
		panic(err)
	}
	respond(c, http.StatusOK, "success", courseSections)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	professorQuery, _ := schema.FilterQuery[schema.Professor](c)

	defer cancel()

	pipeline := mongo.Pipeline{
		// Match professor by first/last name
		bson.D{{Key: "$match", Value: professorQuery}},

		// Expand sections array into individual documents
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$sections"},
				{Key: "preserveNullAndEmptyArrays", Value: false},
			}},
		},

		// Lookup course info using sections.course_reference
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "sections.course_reference"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.course_details"},
			}},
		},

		// Lookup professor info using sections.course_reference
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "sections.professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.professor_details"},
			}},
		},

		// replace the courses with sections
		bson.D{{Key: "$replaceWith", Value: "$sections"}},

		// keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	cursor, err := trendsProfCollection.Aggregate(ctx, pipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	var results []schema.Section

	if err := cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", results)

}
