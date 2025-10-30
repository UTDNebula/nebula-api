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

var sectionCollection *mongo.Collection = configs.GetCollection("sections")

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
func SectionSearch(c *gin.Context) {
	findAndRespond[schema.Section](c, sectionCollection, 10*time.Second)
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
func SectionById(c *gin.Context) {
	findOneByIdAndRespond[schema.Section](c, sectionCollection, 10*time.Second)
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
func SectionCourseSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		sectionCourse("Search", c)
	}
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
func SectionCourseById() gin.HandlerFunc {
	return func(c *gin.Context) {
		sectionCourse("ById", c)
	}
}

// Get an array of courses from sections, filtered based on the flag
func sectionCourse(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sectionCourses []schema.Course
	var sectionQuery bson.M
	var err error
	if sectionQuery, err = getQuery[schema.Section](flag, c); err != nil {
		return
	}

	paginateMap, err := configs.GetAggregateLimit(&sectionQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	// Pipeline to query an array of courses from filtered sections
	sectionCoursePipeline := mongo.Pipeline{
		// Filter the sections
		bson.D{{Key: "$match", Value: sectionQuery}},
	}

	// Paginate the sections before pulling courses from those sections
	formerStages := buildPaginationStages(paginateMap["former_offset"], paginateMap["limit"])
	sectionCoursePipeline = append(sectionCoursePipeline, formerStages...)

	// Lookup the course referenced by sections from the course collection
	sectionCoursePipeline = append(sectionCoursePipeline, buildLookupStage("courses", "course_reference", "_id", "course_reference"))

	// Project to remove every other field except for courses
	sectionCoursePipeline = append(sectionCoursePipeline, buildProjectStage(bson.D{{Key: "courses", Value: "$course_reference"}}))

	// Unwind the courses
	sectionCoursePipeline = append(sectionCoursePipeline, buildUnwindStage("$courses", false))

	// Replace the combinations of id and course with courses entirely
	sectionCoursePipeline = append(sectionCoursePipeline, buildReplaceWithStage("$courses"))

	// Keep order deterministic between calls
	sectionCoursePipeline = append(sectionCoursePipeline, buildSortStage(bson.D{{Key: "_id", Value: 1}}))

	// Paginate the courses
	latterStages := buildPaginationStages(paginateMap["latter_offset"], paginateMap["limit"])
	sectionCoursePipeline = append(sectionCoursePipeline, latterStages...)

	cursor, err := sectionCollection.Aggregate(ctx, sectionCoursePipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Parse the array of courses
	if err = cursor.All(ctx, &sectionCourses); err != nil {
		respondWithInternalError(c, err)
		return
	}

	switch flag {
	case "Search":
		respond(c, http.StatusOK, "success", sectionCourses)
	case "ById":
		// Each section is only referenced by only one course, so returning a single course is ideal
		// A better way of handling this might be needed in the future
		respond(c, http.StatusOK, "success", sectionCourses[0])
	}

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
func SectionProfessorSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		sectionProfessor("Search", c)
	}
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
func SectionProfessorById() gin.HandlerFunc {
	return func(c *gin.Context) {
		sectionProfessor("ById", c)
	}
}

// Get an array of professors from sections
func sectionProfessor(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sectionProfessors []schema.Professor
	var sectionQuery bson.M
	var err error
	if sectionQuery, err = getQuery[schema.Section](flag, c); err != nil {
		return
	}

	paginateMap, err := configs.GetAggregateLimit(&sectionQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	// Pipeline to query an array of professors from filtered sections
	sectionProfessorPipeline := mongo.Pipeline{
		// Filter the sections
		bson.D{{Key: "$match", Value: sectionQuery}},
	}

	// Paginate the sections before pulling professors from those sections
	formerStages := buildPaginationStages(paginateMap["former_offset"], paginateMap["limit"])
	sectionProfessorPipeline = append(sectionProfessorPipeline, formerStages...)

	// Lookup the professors from the professors collection
	sectionProfessorPipeline = append(sectionProfessorPipeline, buildLookupStage("professors", "professors", "_id", "professors"))

	// Project to extract professors
	sectionProfessorPipeline = append(sectionProfessorPipeline, buildProjectStage(bson.D{{Key: "professors", Value: "$professors"}}))

	// Unwind the professors
	sectionProfessorPipeline = append(sectionProfessorPipeline, buildUnwindStage("$professors", false))

	// Replace the root with professors
	sectionProfessorPipeline = append(sectionProfessorPipeline, buildReplaceWithStage("$professors"))

	// Keep order deterministic between calls
	sectionProfessorPipeline = append(sectionProfessorPipeline, buildSortStage(bson.D{{Key: "_id", Value: 1}}))

	// Paginate the professors
	latterStages := buildPaginationStages(paginateMap["latter_offset"], paginateMap["limit"])
	sectionProfessorPipeline = append(sectionProfessorPipeline, latterStages...)

	cursor, err := sectionCollection.Aggregate(ctx, sectionProfessorPipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Parse the array of professors
	if err = cursor.All(ctx, &sectionProfessors); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", sectionProfessors)

}