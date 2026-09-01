package controllers

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/UTDNebula/nebula-api/rest/configs"

	"github.com/UTDNebula/nebula-api/rest/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var sections []schema.Section

	// build query key value pairs (only one value per key)
	query, err := getQuery[schema.Section]("Search", c)
	if err != nil {
		return
	}

	offset, limit, err := configs.GetLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}
	opts := options.Find().SetSkip(offset).SetLimit(limit)

	// get cursor for query results
	cursor, err := sectionCollection.Find(ctx, query, opts)
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
	respond(c, http.StatusOK, "success", sections)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
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
	respond(c, http.StatusOK, "success", section)
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
func SectionCourseSearch(c *gin.Context) {
	sectionAggregate[schema.Course]("Search", c)
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
func SectionCourseById(c *gin.Context) {
	sectionAggregate[schema.Course]("ById", c)
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
func SectionProfessorSearch(c *gin.Context) {
	sectionAggregate[schema.Professor]("Search", c)
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
func SectionProfessorById(c *gin.Context) {
	sectionAggregate[schema.Professor]("ById", c)
}

// sectionAggregate returns the list of aggregated objects from list of filtered sections
func sectionAggregate[T any](flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var queryResults []T
	var sectionQuery bson.M

	// Determine the section query
	sectionQuery, err := getQuery[schema.Section](flag, c)
	if err != nil {
		return
	}

	// Determine the offset and limit for pagination & delete offset fields
	paginateMap, err := configs.GetAggregateLimit(&sectionQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Error offset is not type integer", err.Error())
		return
	}

	// Pipeline to query the field from the filtered sections
	schemaType := strings.Split(reflect.TypeFor[[]T]().String(), ".")[1]
	sectionQueryPipeline := buildSectionPipeline(schemaType, sectionQuery, paginateMap)

	// perform aggregation on the pipeline
	cursor, err := sectionCollection.Aggregate(ctx, sectionQueryPipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &queryResults); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", queryResults)
}

// buildSectionPipeline builds the pipeline to aggregate targets from filtered sections
func buildSectionPipeline(schemaType string, sectionQuery bson.M, paginateMap map[string]bson.D) mongo.Pipeline {
	field := typeToField[schemaType]

	filterSection := mongo.Pipeline{
		bson.D{{Key: "$match", Value: sectionQuery}},

		bson.D{{Key: "$sort", Value: getSort("Section")}},
		paginateMap["former_offset"],
		paginateMap["limit"],
	}

	var lookup = mongo.Pipeline{
		// Lookup the target objects from sections
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: field},
			{Key: "localField", Value: field},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: field},
		}}},
	}

	extract := mongo.Pipeline{
		// Unwind the target objects
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + field},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		// Replace the sections with the target objects
		bson.D{{Key: "$replaceWith", Value: "$" + field}},
	}

	paginate := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: getSort(schemaType)}},

		paginateMap["latter_offset"],
		paginateMap["limit"],
	}

	pipeline := filterSection
	pipeline = append(pipeline, lookup...)
	pipeline = append(pipeline, extract...)
	pipeline = append(pipeline, paginate...)
	return pipeline
}
