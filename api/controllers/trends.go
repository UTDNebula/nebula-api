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

	trendsCourseCollection := configs.GetCollection("trends_course_sections")
	courseQuery := bson.M{"_id": c.Query("subject_prefix") + c.Query("course_number")}

	// Pipeline to query the Sections + Professors from the filtered courses
	pipeline := buildTrendsPipeline(courseQuery)

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
	defer cancel()

	var results []schema.Section

	trendsProfCollection := configs.GetCollection("trends_prof_sections")
	professorQuery, err := schema.FilterQuery[schema.Professor](c)
	if err != nil {
		return
	}

	pipeline := buildTrendsPipeline(professorQuery)

	cursor, err := trendsProfCollection.Aggregate(ctx, pipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if err := cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", results)

}

// buildTrendsPipeline build the pipeline to embed the list of professors and course details
// directly in the queried list of sections
func buildTrendsPipeline(fromObjQuery bson.M) mongo.Pipeline {
	return mongo.Pipeline{
		// Match professors from the query
		bson.D{{Key: "$match", Value: fromObjQuery}},

		// Expand sections array into individual documents
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$sections"},
				{Key: "preserveNullAndEmptyArrays", Value: false}, // Avoid documents that can't be replaced
			}},
		},

		// Lookup courses info using sections.course_reference
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "sections.course_reference"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.course_details"},
			}},
		},

		// Lookup professors info using sections.course_reference
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
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}
}
