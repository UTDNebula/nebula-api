package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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
func (h *CourseController) TrendsCourseSectionSearch(c *gin.Context) {
	// UTDTrends only uses prefix and number
	trendsSections(c, bson.M{"_id": c.Query("subject_prefix") + c.Query("course_number")}, h.repo.TrendsSections)
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
func (h *ProfessorController) TrendsProfessorSectionSearch(c *gin.Context) {
	query, err := schema.FilterQuery[schema.Professor](c.Request.URL.Query())
	if err != nil {
		respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}
	trendsSections(c, query, h.repo.TrendsSections)
}

// @Id				trendsCombinedSectionSearch
// @Router			/combined/sections/trends [get]
// @Tags			Combined
// @Description	"Returns sections matching both course and professor criteria from the combined trends collection. Specialized high-speed convenience endpoint for UTD Trends internal use; limited query flexibility."
// @Produce		json
// @Param			course_number	query		string									true	"The course's official number"
// @Param			subject_prefix	query		string									true	"The course's subject prefix"
// @Param			first_name		query		string									true	"The professor's first name"
// @Param			last_name		query		string									true	"The professor's last name"
// @Success		200				{object}	schema.APIResponse[[]schema.Section]	"A list of Sections"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
func (h *TrendsController) TrendsCombinedSectionSearch(c *gin.Context) {
	query := bson.M{"_id": bson.M{
		"course":     c.Query("subject_prefix") + c.Query("course_number"),
		"prof_first": c.Query("first_name"),
		"prof_last":  c.Query("last_name"),
	}}
	trendsSections(c, query, h.repo.Sections)
}

type TrendsController struct{ repo shared.TrendsRepo }

func NewTrendsController(repo shared.TrendsRepo) *TrendsController {
	return &TrendsController{repo: repo}
}
func trendsSections(c *gin.Context, filter bson.M, fetch func(context.Context, bson.M) ([]schema.Section, error)) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	sections, err := fetch(ctx, filter)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", sections)
}
