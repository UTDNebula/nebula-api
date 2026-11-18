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

// @Id				sectionSearch
// @Router			/section [get]
// @Tags			Sections
// @Description	"Returns paginated list of sections matching the query's string-typed key-value pairs. See offset for more details on pagination."
// @Produce		json
// @Param			offset							query		number									false	"The starting position of the current page of sections (e.g. For starting at the 17th professor, offset=16)."
// @Param			section_number					query		string									false	"The section's official number"
// @Param			academic_session.name			query		string									false	"The name of the academic session of the section"
// @Param			academic_session.start_date		query		string									false	"The date of classes starting for the section"
// @Param			academic_session.end_date		query		string									false	"The date of classes ending for the section"
// @Param			teaching_assistants.first_name	query		string									false	"The first name of one of the teaching assistants of the section"
// @Param			teaching_assistants.last_name	query		string									false	"The last name of one of the teaching assistants of the section"
// @Param			teaching_assistants.role		query		string									false	"The role of one of the teaching assistants of the section"
// @Param			teaching_assistants.email		query		string									false	"The email of one of the teaching assistants of the section"
// @Param			internal_class_number			query		string									false	"The internal (university) number used to reference this section"
// @Param			instruction_mode				query		string									false	"The instruction modality for this section"
// @Param			meetings.start_date				query		string									false	"The start date of one of the section's meetings"
// @Param			meetings.end_date				query		string									false	"The end date of one of the section's meetings"
// @Param			meetings.meeting_days			query		string									false	"One of the days that one of the section's meetings"
// @Param			meetings.start_time				query		string									false	"The time one of the section's meetings starts"
// @Param			meetings.end_time				query		string									false	"The time one of the section's meetings ends"
// @Param			meetings.modality				query		string									false	"The modality of one of the section's meetings"
// @Param			meetings.location.building		query		string									false	"The building of one of the section's meetings"
// @Param			meetings.location.room			query		string									false	"The room of one of the section's meetings"
// @Param			meetings.location.map_uri		query		string									false	"A hyperlink to the UTD room locator of one of the section's meetings"
// @Param			core_flags						query		string									false	"One of core requirement codes this section fulfills"
// @Param			syllabus_uri					query		string									false	"A link to the syllabus on the web"
// @Success		200								{object}	schema.APIResponse[[]schema.Section]	"A list of sections"
// @Failure		500								{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400								{object}	schema.APIResponse[string]				"A string describing the error"
func (h *SectionController) SectionSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// build query key value pairs (only one value per key)
	query, err := getQuery[schema.Section]("Search", c)
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

// @Id				sectionById
// @Router			/section/{id} [get]
// @Tags			Sections
// @Description	"Returns the section with given ID"
// @Produce		json
// @Param			id	path		string								true	"ID of the section to get"
// @Success		200	{object}	schema.APIResponse[schema.Section]	"A section"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]			"A string describing the error"
func (h *SectionController) SectionById(c *gin.Context) {
	// parse object id from id parameter
	id, err := objectIDFromParam(c, "id")
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// find and parse matching section
	record, err := h.repo.ByID(ctx, *id)
	if errors.Is(err, shared.ErrNotFound) {
		respond(c, http.StatusNotFound, "error", "No sections with given ID")
		return
	}
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", record)
}

// @Id				sectionCourseSearch
// @Router			/section/courses [get]
// @Tags			Sections
// @Description	"Returns paginated list of courses of all the sections matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset					query		number								false	"The starting position of the current page of sections (e.g. For starting at the 16th section, former_offset=16)."
// @Param			latter_offset					query		number								false	"The starting position of the current page of courses (e.g. For starting at the 16th course, latter_offset=16)."
// @Param			section_number					query		string								false	"The section's official number"
// @Param			academic_session.name			query		string								false	"The name of the academic session of the section"
// @Param			academic_session.start_date		query		string								false	"The date of classes starting for the section"
// @Param			academic_session.end_date		query		string								false	"The date of classes ending for the section"
// @Param			teaching_assistants.first_name	query		string								false	"The first name of one of the teaching assistants of the section"
// @Param			teaching_assistants.last_name	query		string								false	"The last name of one of the teaching assistants of the section"
// @Param			teaching_assistants.role		query		string								false	"The role of one of the teaching assistants of the section"
// @Param			teaching_assistants.email		query		string								false	"The email of one of the teaching assistants of the section"
// @Param			internal_class_number			query		string								false	"The internal (university) number used to reference this section"
// @Param			instruction_mode				query		string								false	"The instruction modality for this section"
// @Param			meetings.start_date				query		string								false	"The start date of one of the section's meetings"
// @Param			meetings.end_date				query		string								false	"The end date of one of the section's meetings"
// @Param			meetings.meeting_days			query		string								false	"One of the days that one of the section's meetings"
// @Param			meetings.start_time				query		string								false	"The time one of the section's meetings starts"
// @Param			meetings.end_time				query		string								false	"The time one of the section's meetings ends"
// @Param			meetings.modality				query		string								false	"The modality of one of the section's meetings"
// @Param			meetings.location.building		query		string								false	"The building of one of the section's meetings"
// @Param			meetings.location.room			query		string								false	"The room of one of the section's meetings"
// @Param			meetings.location.map_uri		query		string								false	"A hyperlink to the UTD room locator of one of the section's meetings"
// @Param			core_flags						query		string								false	"One of core requirement codes this section fulfills"
// @Param			syllabus_uri					query		string								false	"A link to the syllabus on the web"
// @Success		200								{object}	schema.APIResponse[[]schema.Course]	"A list of courses"
// @Failure		500								{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400								{object}	schema.APIResponse[string]			"A string describing the error"
func (h *SectionController) SectionCourseSearch(c *gin.Context) {
	respondRelated[schema.Section, schema.Course]("Search", c, h.repo.Courses, "Error offset is not type integer")
}

// @Id				sectionCourseById
// @Router			/section/{id}/course [get]
// @Tags			Sections
// @Description	"Returns the course of the section with given ID"
// @Produce		json
// @Param			id	path		string								true	"ID of the section to get"
// @Success		200	{object}	schema.APIResponse[schema.Course]	"A course"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		404	{object}	schema.APIResponse[string]			"A string describing the error"
func (h *SectionController) SectionCourseById(c *gin.Context) {
	courses, ok := relationshipResults[schema.Section, schema.Course]("ById", c, h.repo.Courses, "Error offset is not type integer")
	if !ok {
		return
	}
	if len(courses) == 0 {
		respond(c, http.StatusNotFound, "error", "No section with given ID")
		return
	}
	respond(c, http.StatusOK, "success", courses[0])
}

// @Id				sectionProfessorSearch
// @Router			/section/professors [get]
// @Tags			Sections
// @Description	"Returns paginated list of professors of all the sections matching the query's string-typed key-value pairs. See former_offset and latter_offset for pagination details."
// @Produce		json
// @Param			former_offset					query		number									false	"The starting position of the current page of sections (e.g. For starting at the 16th sections, former_offset=16)."
// @Param			latter_offset					query		number									false	"The starting position of the current page of professors (e.g. For starting at the 16th professor, latter_offset=16)."
// @Param			section_number					query		string									false	"The section's official number"
// @Param			academic_session.name			query		string									false	"The name of the academic session of the section"
// @Param			academic_session.start_date		query		string									false	"The date of classes starting for the section"
// @Param			academic_session.end_date		query		string									false	"The date of classes ending for the section"
// @Param			teaching_assistants.first_name	query		string									false	"The first name of one of the teaching assistants of the section"
// @Param			teaching_assistants.last_name	query		string									false	"The last name of one of the teaching assistants of the section"
// @Param			teaching_assistants.role		query		string									false	"The role of one of the teaching assistants of the section"
// @Param			teaching_assistants.email		query		string									false	"The email of one of the teaching assistants of the section"
// @Param			internal_class_number			query		string									false	"The internal (university) number used to reference this section"
// @Param			instruction_mode				query		string									false	"The instruction modality for this section"
// @Param			meetings.start_date				query		string									false	"The start date of one of the section's meetings"
// @Param			meetings.end_date				query		string									false	"The end date of one of the section's meetings"
// @Param			meetings.meeting_days			query		string									false	"One of the days that one of the section's meetings"
// @Param			meetings.start_time				query		string									false	"The time one of the section's meetings starts"
// @Param			meetings.end_time				query		string									false	"The time one of the section's meetings ends"
// @Param			meetings.modality				query		string									false	"The modality of one of the section's meetings"
// @Param			meetings.location.building		query		string									false	"The building of one of the section's meetings"
// @Param			meetings.location.room			query		string									false	"The room of one of the section's meetings"
// @Param			meetings.location.map_uri		query		string									false	"A hyperlink to the UTD room locator of one of the section's meetings"
// @Param			core_flags						query		string									false	"One of core requirement codes this section fulfills"
// @Param			syllabus_uri					query		string									false	"A link to the syllabus on the web"
// @Success		200								{object}	schema.APIResponse[[]schema.Professor]	"A list of professor"
// @Failure		500								{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400								{object}	schema.APIResponse[string]				"A string describing the error"
func (h *SectionController) SectionProfessorSearch(c *gin.Context) {
	respondRelated[schema.Section, schema.Professor]("Search", c, h.repo.Professors, "offset is not type integer")
}

// @Id				sectionProfessorById
// @Router			/section/{id}/professors [get]
// @Tags			Sections
// @Description	"Returns the paginated list of professors of the section with given ID"
// @Produce		json
// @Param			id	path		string									true	"ID of the section to get"
// @Success		200	{object}	schema.APIResponse[[]schema.Professor]	"A list of professors"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func (h *SectionController) SectionProfessorById(c *gin.Context) {
	respondRelated[schema.Section, schema.Professor]("ById", c, h.repo.Professors, "offset is not type integer")
}

// SectionController handles requests using a supplied repository.
type SectionController struct{ repo shared.SectionRepo }

func NewSectionController(repo shared.SectionRepo) *SectionController {
	return &SectionController{repo: repo}
}
