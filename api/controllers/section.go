package controllers

import (
	"context"
	"errors"
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
	//name := c.Query("name")            // value of specific query parameter: string
	//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sections []schema.Section

	// build query key value pairs (only one value per key)
	query, err := getQuery[schema.Section]("Search", c)
	if err != nil {
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond[string](c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	// get cursor for query results
	cursor, err := sectionCollection.Find(ctx, query, optionLimit)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &sections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// return result
	respond[[]schema.Section](c, http.StatusOK, "success", sections)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var section schema.Section

	// parse object id from id parameter
	query, err := getQuery[schema.Section]("ById", c)
	if err != nil {
		return
	}

	// find and parse matching section
	err = sectionCollection.FindOne(ctx, query).Decode(&section)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No sections with given ID")
		} else {
			respondWithInternalError(c, err)
		}
		return
	}

	// return result
	respond[schema.Section](c, http.StatusOK, "success", section)
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

// Get an array of courses from sections, filtered based on the the flag
func sectionCourse(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sectionCourses []schema.Course
	var sectionQuery bson.M
	var err error
	if sectionQuery, err = getQuery[schema.Section](flag, c); err != nil {
		return
	}

	rawPaginateMap, err := configs.GetAggregateLimit(&sectionQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	paginateMap := make(map[string]int)
	for k, v := range rawPaginateMap {
		paginateMap[k] = int(v)
	}

	pipeline := buildSectionPipeline(sectionQuery, paginateMap, "courses", flag == "ById")
	cursor, err := sectionCollection.Aggregate(ctx, pipeline)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if flag == "ById" {
		var course schema.Course
		if cursor.Next(ctx) {
			if err := cursor.Decode(&course); err != nil {
				respondWithInternalError(c, err)
				return
			}
			respond[*schema.Course](c, http.StatusOK, "success", &course)
			return
		}
		respond[interface{}](c, http.StatusOK, "success", nil)
	} else {
		if err := cursor.All(ctx, &sectionCourses); err != nil {
			respondWithInternalError(c, err)
			return
		}
		respond[[]schema.Course](c, http.StatusOK, "success", sectionCourses)
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

// Get an array of professors sections,
func sectionProfessor(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sectionProfessors []schema.Professor
	var sectionQuery bson.M
	var err error
	if sectionQuery, err = getQuery[schema.Section](flag, c); err != nil {
		return
	}

	rawPaginateMap, err := configs.GetAggregateLimit(&sectionQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	paginateMap := make(map[string]int)
	for k, v := range rawPaginateMap {
		paginateMap[k] = int(v)
	}

	pipeline := buildSectionPipeline(sectionQuery, paginateMap, "professors", flag == "ById")
	cursor, err := sectionCollection.Aggregate(ctx, pipeline)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if err = cursor.All(ctx, &sectionProfessors); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", sectionProfessors)

}

func buildSectionPipeline(
	sectionQuery bson.M,
	paginateMap map[string]int,
	lookupType string,
	single bool,
) mongo.Pipeline {
	localField := "course_reference"
	field := lookupType

	if lookupType == "professors" {
		localField = "professor_id"
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: sectionQuery}},
	}
	if !single {
		pipeline = append(pipeline,
			bson.D{{Key: "$skip", Value: paginateMap["former_offset"]}},
			bson.D{{Key: "$limit", Value: paginateMap["limit"]}},
		)
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: lookupType},
			{Key: "localField", Value: localField},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: field},
		}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: field, Value: "$" + field}}}},
	)
	// unwind/replaceWith so the aggregation yields the joined document itself.
	pipeline = append(pipeline,
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + field},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},
		bson.D{{Key: "$replaceWith", Value: "$" + field}},
	)

	// For non-single (search) requests, apply sorting and latter_offset pagination
	// after unwinding so we return a paginated list of the joined documents.
	if !single {
		pipeline = append(pipeline,
			bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
			bson.D{{Key: "$skip", Value: paginateMap["latter_offset"]}},
			bson.D{{Key: "$limit", Value: paginateMap["limit"]}},
		)
	}
	return pipeline
}
