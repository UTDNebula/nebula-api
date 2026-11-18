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

// @Id				professorSearch
// @Router			/professor [get]
// @Tags			Professors
// @Description	"Returns paginated list of professors matching the query's string-typed key-value pairs. See offset for more details on pagination."
// @Produce		json
// @Param			offset							query		number									false	"The starting position of the current page of professors (e.g. For starting at the 17th professor, offset=16)."
// @Param			first_name						query		string									false	"The professor's first name"
// @Param			last_name						query		string									false	"The professor's last name"
// @Param			titles							query		string									false	"One of the professor's title"
// @Param			email							query		string									false	"The professor's email address"
// @Param			phone_number					query		string									false	"The professor's phone number"
// @Param			office.building					query		string									false	"The building of the location of the professor's office"
// @Param			office.room						query		string									false	"The room of the location of the professor's office"
// @Param			office.map_uri					query		string									false	"A hyperlink to the UTD room locator of the professor's office"
// @Param			profile_uri						query		string									false	"A hyperlink pointing to the professor's official university profile"
// @Param			image_uri						query		string									false	"A link to the image used for the professor on the professor's official university profile"
// @Param			office_hours.start_date			query		string									false	"The start date of one of the office hours meetings of the professor"
// @Param			office_hours.end_date			query		string									false	"The end date of one of the office hours meetings of the professor"
// @Param			office_hours.meeting_days		query		string									false	"One of the days that one of the office hours meetings of the professor"
// @Param			office_hours.start_time			query		string									false	"The time one of the office hours meetings of the professor starts"
// @Param			office_hours.end_time			query		string									false	"The time one of the office hours meetings of the professor ends"
// @Param			office_hours.modality			query		string									false	"The modality of one of the office hours meetings of the professor"
// @Param			office_hours.location.building	query		string									false	"The building of one of the office hours meetings of the professor"
// @Param			office_hours.location.room		query		string									false	"The room of one of the office hours meetings of the professor"
// @Param			office_hours.location.map_uri	query		string									false	"A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Success		200								{object}	schema.APIResponse[[]schema.Professor]	"A list of professors"
// @Failure		500								{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400								{object}	schema.APIResponse[string]				"A string describing the error"
func (h *ProfessorController) ProfessorSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// build query key value pairs (only one value per key)
	query, err := getQuery[schema.Professor]("Search", c)
	if err != nil {
		return
	}
	offset, ok := searchOffset(c)
	if !ok {
		return
	}
	records, err := h.repo.Search(ctx, query, offset, configs.GetEnvLimit())
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", records)
}

// @Id				professorById
// @Router			/professor/{id} [get]
// @Tags			Professors
// @Description	"Returns the professor with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the professor to get"
// @Success		200	{object}	schema.APIResponse[schema.Professor]	"A professor"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func (h *ProfessorController) ProfessorById(c *gin.Context) {
	// parse object id from id parameter
	id, err := objectIDFromParam(c, "id")
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// find and parse matching professor
	record, err := h.repo.ByID(ctx, *id)
	if errors.Is(err, shared.ErrNotFound) {
		respond(c, http.StatusNotFound, "error", "No professors with given ID")
		return
	}
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", record)
}

// @Id				professorAll
// @Router			/professor/all [get]
// @Tags			Professors
// @Description	"Returns all professors"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]schema.Professor]	"All professors"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
func (h *ProfessorController) ProfessorAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	records, err := h.repo.All(ctx)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", records)
}

// @Id				professorCourseSearch
// @Router			/professor/courses [get]
// @Tags			Professors
// @Description	"Returns paginated list of the courses of all the professors matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset					query		number									false	"The starting position of the current page of professors (e.g. For starting at the 17th professor, former_offset=16)."
// @Param			latter_offset					query		number									false	"The starting position of the current page of courses (e.g. For starting at the 4th course, latter_offset=3)."
// @Param			first_name						query		string									false	"The professor's first name"
// @Param			last_name						query		string									false	"The professor's last name"
// @Param			titles							query		string									false	"One of the professor's title"
// @Param			email							query		string									false	"The professor's email address"
// @Param			phone_number					query		string									false	"The professor's phone number"
// @Param			office.building					query		string									false	"The building of the location of the professor's office"
// @Param			office.room						query		string									false	"The room of the location of the professor's office"
// @Param			office.map_uri					query		string									false	"A hyperlink to the UTD room locator of the professor's office"
// @Param			profile_uri						query		string									false	"A hyperlink pointing to the professor's official university profile"
// @Param			image_uri						query		string									false	"A link to the image used for the professor on the professor's official university profile"
// @Param			office_hours.start_date			query		string									false	"The start date of one of the office hours meetings of the professor"
// @Param			office_hours.end_date			query		string									false	"The end date of one of the office hours meetings of the professor"
// @Param			office_hours.meeting_days		query		string									false	"One of the days that one of the office hours meetings of the professor"
// @Param			office_hours.start_time			query		string									false	"The time one of the office hours meetings of the professor starts"
// @Param			office_hours.end_time			query		string									false	"The time one of the office hours meetings of the professor ends"
// @Param			office_hours.modality			query		string									false	"The modality of one of the office hours meetings of the professor"
// @Param			office_hours.location.building	query		string									false	"The building of one of the office hours meetings of the professor"
// @Param			office_hours.location.room		query		string									false	"The room of one of the office hours meetings of the professor"
// @Param			office_hours.location.map_uri	query		string									false	"A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Success		200								{object}	schema.APIResponse[[]schema.Professor]	"A list of courses"
// @Failure		500								{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400								{object}	schema.APIResponse[string]				"A string describing the error"
func (h *ProfessorController) ProfessorCourseSearch(c *gin.Context) {
	respondRelated[schema.Professor, schema.Course]("Search", c, h.repo.Courses, "offset is not type integer")
}

// @Id				professorCourseById
// @Router			/professor/{id}/courses [get]
// @Tags			Professors
// @Description	"Returns all the courses taught by the professor with given ID"
// @Produce		json
// @Param			id	path		string								true	"ID of the professor to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Course]	"A list of courses"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]			"A string describing the error"
func (h *ProfessorController) ProfessorCourseById(c *gin.Context) {
	respondRelated[schema.Professor, schema.Course]("ById", c, h.repo.Courses, "offset is not type integer")
}

// @Id				professorSectionSearch
// @Router			/professor/sections [get]
// @Tags			Professors
// @Description	"Returns paginated list of the sections of all the professors matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset					query		number									false	"The starting position of the current page of professors (e.g. For starting at the 17th professor, former_offset=16)."
// @Param			latter_offset					query		number									false	"The starting position of the current page of sections (e.g. For starting at the 4th section, latter_offset=3)."
// @Param			first_name						query		string									false	"The professor's first name"
// @Param			last_name						query		string									false	"The professor's last name"
// @Param			titles							query		string									false	"One of the professor's title"
// @Param			email							query		string									false	"The professor's email address"
// @Param			phone_number					query		string									false	"The professor's phone number"
// @Param			office.building					query		string									false	"The building of the location of the professor's office"
// @Param			office.room						query		string									false	"The room of the location of the professor's office"
// @Param			office.map_uri					query		string									false	"A hyperlink to the UTD room locator of the professor's office"
// @Param			profile_uri						query		string									false	"A hyperlink pointing to the professor's official university profile"
// @Param			image_uri						query		string									false	"A link to the image used for the professor on the professor's official university profile"
// @Param			office_hours.start_date			query		string									false	"The start date of one of the office hours meetings of the professor"
// @Param			office_hours.end_date			query		string									false	"The end date of one of the office hours meetings of the professor"
// @Param			office_hours.meeting_days		query		string									false	"One of the days that one of the office hours meetings of the professor"
// @Param			office_hours.start_time			query		string									false	"The time one of the office hours meetings of the professor starts"
// @Param			office_hours.end_time			query		string									false	"The time one of the office hours meetings of the professor ends"
// @Param			office_hours.modality			query		string									false	"The modality of one of the office hours meetings of the professor"
// @Param			office_hours.location.building	query		string									false	"The building of one of the office hours meetings of the professor"
// @Param			office_hours.location.room		query		string									false	"The room of one of the office hours meetings of the professor"
// @Param			office_hours.location.map_uri	query		string									false	"A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Success		200								{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500								{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400								{object}	schema.APIResponse[string]				"A string describing the error"
func (h *ProfessorController) ProfessorSectionSearch(c *gin.Context) {
	respondRelated[schema.Professor, schema.Section]("Search", c, h.repo.Sections, "offset is not type integer")
}

// @Id				professorSectionById
// @Router			/professor/{id}/sections [get]
// @Tags			Professors
// @Description	"Returns all the sections taught by the professor with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the professor to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func (h *ProfessorController) ProfessorSectionById(c *gin.Context) {
	respondRelated[schema.Professor, schema.Section]("ById", c, h.repo.Sections, "offset is not type integer")
}

// ProfessorController handles requests using a supplied repository.
type ProfessorController struct{ repo shared.ProfessorRepo }

func NewProfessorController(repo shared.ProfessorRepo) *ProfessorController {
	return &ProfessorController{repo: repo}
}
