package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
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

// 5 Functions

// @Id				gradeAggregationBySemester
// @Router			/grades/semester [get]
// @Tags			Grades
// @Description	"Returns grade distributions aggregated by semester"
// @Produce		json
// @Param			prefix			query		string									false	"The course's subject prefix"
// @Param			number			query		string									false	"The course's official number"
// @Param			first_name		query		string									false	"The professor's first name"
// @Param			last_name		query		string									false	"The professors's last name"
// @Param			section_number	query		string									false	"The number of the section"
// @Success		200				{object}	schema.APIResponse[[]schema.GradeData]	"An array of grade distributions for each semester included"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400				{object}	schema.APIResponse[string]				"A string describing the error"
func (h *GradesController) GradeAggregationSemester(c *gin.Context) {
	h.aggregate("semester", c)
}

// @Id				gradeAggregationSectionType
// @Router			/grades/semester/sectionType [get]
// @Tags			Grades
// @Description	"Returns the grade distributions aggregated by semester and broken down into section type"
// @Produce		json
// @Param			prefix			query		string										false	"The course's subject prefix"
// @Param			number			query		string										false	"The course's official number"
// @Param			first_name		query		string										false	"The professor's first name"
// @Param			last_name		query		string										false	"The professors's last name"
// @Param			section_number	query		string										false	"The number of the section"
// @Success		200				{object}	schema.APIResponse[[]schema.TypedGradeData]	"An array of grade distributions for each section type for each semester included"
// @Failure		500				{object}	schema.APIResponse[string]					"A string describing the error"
// @Failure		400				{object}	schema.APIResponse[string]					"A string describing the error"
func (h *GradesController) GradesAggregationSectionType(c *gin.Context) {
	h.aggregate("section_type", c)
}

// @Id				gradeAggregationOverall
// @Router			/grades/overall [get]
// @Tags			Grades
// @Description	"Returns the overall grade distribution"
// @Produce		json
// @Param			prefix			query		string						false	"The course's subject prefix"
// @Param			number			query		string						false	"The course's official number"
// @Param			first_name		query		string						false	"The professor's first name"
// @Param			last_name		query		string						false	"The professors's last name"
// @Param			section_number	query		string						false	"The number of the section"
// @Success		200				{object}	schema.APIResponse[[]int]	"A grade distribution array"
// @Failure		500				{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400				{object}	schema.APIResponse[string]	"A string describing the error"
func (h *GradesController) GradesAggregationOverall(c *gin.Context) {
	h.aggregate("overall", c)
}

// @Id				GradesByCourseID
// @Router			/course/{id}/grades [get]
// @Tags			Courses
// @Description	"Returns the overall grade distribution for a course"
// @Produce		json
// @Param			id	path		string						true	"ID of course to get grades for"
// @Success		200	{object}	schema.APIResponse[[]int]	"A grade distribution array for the course"
// @Failure		500	{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]	"A string describing the error"
func (h *CourseController) GradesByCourseID(c *gin.Context) {
	gradesByID(c, h.repo.Grades)
}

// @Id				GradesBySectionID
// @Router			/section/{id}/grades [get]
// @Tags			Sections
// @Description	"Returns the overall grade distribution for a section"
// @Produce		json
// @Param			id	path		string						true	"ID of section to get grades for"
// @Success		200	{object}	schema.APIResponse[[]int]	"A grade distribution array for the section"
// @Failure		500	{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]	"A string describing the error"
func (h *SectionController) GradesBySectionID(c *gin.Context) {
	gradesByID(c, h.repo.Grades)
}

// @Id				GradesByProfessorID
// @Router			/professor/{id}/grades [get]
// @Tags			Professors
// @Description	"Returns the overall grade distribution for a professor"
// @Produce		json
// @Param			id	path		string						true	"ID of professor to get grades for"
// @Success		200	{object}	schema.APIResponse[[]int]	"A grade distribution array for the professor"
// @Failure		500	{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]	"A string describing the error"
func (h *ProfessorController) GradesByProfessorID(c *gin.Context) {
	gradesByID(c, h.repo.Grades)
}

type GradesController struct{ repo shared.GradesRepo }

func NewGradesController(repo shared.GradesRepo) *GradesController {
	return &GradesController{repo: repo}
}
func (h *GradesController) aggregate(mode string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	result, err := h.repo.Aggregate(ctx, shared.GradeQuery{
		Mode:          mode,
		Prefix:        c.Query("prefix"),
		Number:        c.Query("number"),
		SectionNumber: c.Query("section_number"),
		FirstName:     c.Query("first_name"),
		LastName:      c.Query("last_name"),
	})
	if errors.Is(err, shared.ErrInvalidGradeQuery) {
		respond(c, http.StatusBadRequest, "error", err.Error())
		return
	}
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	switch mode {
	case "overall":
		respond(c, http.StatusOK, "success", result.Overall)
	case "semester":
		respond(c, http.StatusOK, "success", result.Semesters)
	case "section_type":
		respond(c, http.StatusOK, "success", result.SectionTypes)
	}
}
