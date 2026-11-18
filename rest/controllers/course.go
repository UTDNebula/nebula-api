package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

type CourseController struct {
	repo shared.CourseRepo
}

func NewCourseController(repo shared.CourseRepo) *CourseController {
	return &CourseController{
		repo: repo,
	}
}

// @Id				courseSearch
// @Router			/course [get]
// @Tags			Courses
// @Description	"Returns paginated list of courses matching the query's string-typed key-value pairs. See offset for more details on pagination."
// @Produce		json
// @Param			offset					query		number								false	"The starting position of the current page of courses (e.g. For starting at the 17th course, offset=16)."
// @Param			course_number			query		string								false	"The course's official number"
// @Param			subject_prefix			query		string								false	"The course's subject prefix"
// @Param			title					query		string								false	"The course's title"
// @Param			description				query		string								false	"The course's description"
// @Param			school					query		string								false	"The course's school"
// @Param			credit_hours			query		string								false	"The number of credit hours awarded by successful completion of the course"
// @Param			class_level				query		string								false	"The level of education that this course course corresponds to"
// @Param			activity_type			query		string								false	"The type of class this course corresponds to"
// @Param			grading					query		string								false	"The grading status of this course"
// @Param			internal_course_number	query		string								false	"The internal (university) number used to reference this course"
// @Param			lecture_contact_hours	query		string								false	"The weekly contact hours in lecture for a course"
// @Param			offering_frequency		query		string								false	"The frequency of offering a course"
// @Success		200						{object}	schema.APIResponse[[]schema.Course]	"A list of courses"
// @Failure		500						{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400						{object}	schema.APIResponse[string]			"A string describing the error"
func (h *CourseController) CourseSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// Build query key value pairs (only one value per key)
	query, err := getQuery[schema.Course]("Search", c)
	if err != nil {
		return
	}
	offset, ok := searchOffset(c)
	if !ok {
		return
	}
	courses, err := h.repo.Search(ctx, query, offset, configs.GetEnvLimit())
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", courses)
}

// @Id				courseById
// @Router			/course/{id} [get]
// @Tags			Courses
// @Description	"Returns the course with given ID"
// @Produce		json
// @Param			id	path		string								true	"ID of the course to get"
// @Success		200	{object}	schema.APIResponse[schema.Course]	"A course"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
func (h *CourseController) CourseById(c *gin.Context) {
	// parse object id from id parameter
	id, err := objectIDFromParam(c, "id")
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// find and parse matching course
	course, err := h.repo.ByID(ctx, *id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			respond(c, http.StatusNotFound, "error", "No courses with given ID")
		} else {
			respondWithInternalError(c, err)
		}
		return
	}

	// return result
	respond(c, http.StatusOK, "success", course)
}

// @Id				courseAll
// @Router			/course/all [get]
// @Tags			Courses
// @Description	"Returns all courses"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]schema.Course]	"All courses"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
func (h *CourseController) CourseAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	courses, err := h.repo.All(ctx)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", courses)
}

// @Id				courseSectionSearch
// @Router			/course/sections [get]
// @Tags			Courses
// @Description	"Returns paginated list of sections of all the courses matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset			query		number									false	"The starting position of the current page of courses (e.g. For starting at the 17th course, former_offset=16)."
// @Param			latter_offset			query		number									false	"The starting position of the current page of sections (e.g. For starting at the 4th section, latter_offset=3)."
// @Param			course_number			query		string									false	"The course's official number"
// @Param			subject_prefix			query		string									false	"The course's subject prefix"
// @Param			title					query		string									false	"The course's title"
// @Param			description				query		string									false	"The course's description"
// @Param			school					query		string									false	"The course's school"
// @Param			credit_hours			query		string									false	"The number of credit hours awarded by successful completion of the course"
// @Param			class_level				query		string									false	"The level of education that this course course corresponds to"
// @Param			activity_type			query		string									false	"The type of class this course corresponds to"
// @Param			grading					query		string									false	"The grading status of this course"
// @Param			internal_course_number	query		string									false	"The internal (university) number used to reference this course"
// @Param			lecture_contact_hours	query		string									false	"The weekly contact hours in lecture for a course"
// @Param			offering_frequency		query		string									false	"The frequency of offering a course"
// @Success		200						{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500						{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400						{object}	schema.APIResponse[string]				"A string describing the error"
func (h *CourseController) CourseSectionSearch(c *gin.Context) {
	respondRelated[schema.Course, schema.Section]("Search", c, h.repo.Sections, "Error offset is not type integer")
}

// @Id				courseSectionById
// @Router			/course/{id}/sections [get]
// @Tags			Courses
// @Description	"Returns the all of the sections of the course with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the course to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func (h *CourseController) CourseSectionById(c *gin.Context) {
	respondRelated[schema.Course, schema.Section]("ById", c, h.repo.Sections, "Error offset is not type integer")
}

// @Id				courseProfessorSearch
// @Router			/course/professors [get]
// @Tags			Courses
// @Description	"Returns paginated list of professors of all the courses matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset			query		number									false	"The starting position of the current page of courses (e.g. For starting at the 17th course, former_offset=16)."
// @Param			latter_offset			query		number									false	"The starting position of the current page of professors (e.g. For starting at the 4th professor, latter_offset=3)."
// @Param			course_number			query		string									false	"The course's official number"
// @Param			subject_prefix			query		string									false	"The course's subject prefix"
// @Param			title					query		string									false	"The course's title"
// @Param			description				query		string									false	"The course's description"
// @Param			school					query		string									false	"The course's school"
// @Param			credit_hours			query		string									false	"The number of credit hours awarded by successful completion of the course"
// @Param			class_level				query		string									false	"The level of education that this course course corresponds to"
// @Param			activity_type			query		string									false	"The type of class this course corresponds to"
// @Param			grading					query		string									false	"The grading status of this course"
// @Param			internal_course_number	query		string									false	"The internal (university) number used to reference this course"
// @Param			lecture_contact_hours	query		string									false	"The weekly contact hours in lecture for a course"
// @Param			offering_frequency		query		string									false	"The frequency of offering a course"
// @Success		200						{object}	schema.APIResponse[[]schema.Professor]	"A list of professors"
// @Failure		500						{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400						{object}	schema.APIResponse[string]				"A string describing the error"
func (h *CourseController) CourseProfessorSearch(c *gin.Context) {
	respondRelated[schema.Course, schema.Professor]("Search", c, h.repo.Professors, "Error offset is not type integer")
}

// @Id				courseProfessorById
// @Router			/course/{id}/professors [get]
// @Tags			Courses
// @Description	"Returns the all of the professors of the course with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the course to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Professor]	"A list of professors"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func (h *CourseController) CourseProfessorById(c *gin.Context) {
	respondRelated[schema.Course, schema.Professor]("ById", c, h.repo.Professors, "Error offset is not type integer")
}
