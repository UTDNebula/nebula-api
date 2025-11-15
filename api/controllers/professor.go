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

var professorCollection *mongo.Collection = configs.GetCollection("professors")
var trendsProfCollection *mongo.Collection = configs.GetCollection("trends_prof_sections")

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
func ProfessorSearch(c *gin.Context) {
	//name := c.Query("name")            // value of specific query parameter: string
	//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var professors []schema.Professor

	// build query key value pairs (only one value per key)
	query, err := getQuery[schema.Professor]("Search", c)
	if err != nil {
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	// get cursor for query results
	cursor, err := professorCollection.Find(ctx, query, optionLimit)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &professors); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// return result
	respond(c, http.StatusOK, "success", professors)
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
func ProfessorById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var professor schema.Professor

	// parse object id from id parameter
	query, err := getQuery[schema.Professor]("ById", c)
	if err != nil {
		return
	}

	// find and parse matching professor
	err = professorCollection.FindOne(ctx, query).Decode(&professor)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// return result
	respond(c, http.StatusOK, "success", professor)
}

// @Id				professorAll
// @Router			/professor/all [get]
// @Tags			Professors
// @Description	"Returns all professors"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]schema.Professor]	"All professors"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
func ProfessorAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var professors []schema.Professor

	cursor, err := professorCollection.Find(ctx, bson.M{})

	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &professors); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// return result
	respond(c, http.StatusOK, "success", professors)
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
func ProfessorCourseSearch() gin.HandlerFunc {
	// Wrapper of professorCourse() with flag of Search
	return func(c *gin.Context) {
		professorCourse("Search", c)
	}
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
func ProfessorCourseById() gin.HandlerFunc {
	// Essentially wrapper of professorCourse() with flag of ById
	return func(c *gin.Context) {
		professorCourse("ById", c)
	}
}

// Pipeline builder for professor aggregate endpoints
func professorPipeline(endpoint string, professorQuery bson.M, paginateMap map[string]int64) mongo.Pipeline {
	// common stages
	baseStages := mongo.Pipeline{
		// filter the professors
		bson.D{{Key: "$match", Value: professorQuery}},

		// paginate the professors before pulling the courses from those professor
		bson.D{{Key: "$skip", Value: paginateMap["former_offset"]}}, // skip to the specified offset
		bson.D{{Key: "$limit", Value: paginateMap["limit"]}},        // limit to the specified number of professors

		// lookup the array of sections from sections collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},

		// keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	// course pagination stages
	paginationStages := mongo.Pipeline{
		bson.D{{Key: "$skip", Value: paginateMap["latter_offset"]}},
		bson.D{{Key: "$limit", Value: paginateMap["limit"]}},
	}

	if endpoint == "courses" {
		courseStages := mongo.Pipeline{
			// project the courses referenced by each section in the array
			bson.D{{Key: "$project", Value: bson.D{{Key: "courses", Value: "$sections.course_reference"}}}},

			// lookup the array of courses from coures collection
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "courses"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "courses"},
			}}},

			// unwind the courses
			bson.D{{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$courses"},
				{Key: "preserveNullAndEmptyArrays", Value: false}, // to avoid the professor documents that can't be replaced
			}}},

			// replace the combination of ids and courses with the courses entirely
			bson.D{{Key: "$replaceWith", Value: "$courses"}},
		}

		return append(append(baseStages, courseStages...), paginationStages...)
	}

	if endpoint == "sections" {
		sectionStages := mongo.Pipeline{
			// project the sections
			bson.D{{Key: "$project", Value: bson.D{{Key: "sections", Value: "$sections"}}}},

			// unwind the sections
			bson.D{{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$sections"},
				{Key: "preserveNullAndEmptyArrays", Value: false}, // to avoid the professor documents that can't be replaced
			}}},

			// replace the combination of ids and sections with the sections entirely
			bson.D{{Key: "$replaceWith", Value: "$sections"}},
		}

		return append(append(baseStages, sectionStages...), paginationStages...)
	}

	return append(baseStages, paginationStages...) // fallback (shouldn't happen because we call with either courses or sections)
}

// Get all of the courses of the professors depending on the type of flag
func professorCourse(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var professorCourses []schema.Course // array of courses of the professors (or single professor with Id)
	var professorQuery bson.M            // query filter the professor
	var err error

	// determine the professor's query
	professorQuery, err = getQuery[schema.Professor](flag, c)
	if err != nil {
		return
	}

	// determine the offset and limit for pagination stage
	// and delete "offset" field in professorQuery
	paginateMap, err := configs.GetAggregateLimit(&professorQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	// Pipeline to query the courses from the filtered professors (or a single professor)
	professorCoursePipeline := professorPipeline("courses", professorQuery, paginateMap)

	// Perform aggreration on the pipeline
	cursor, err := professorCollection.Aggregate(ctx, professorCoursePipeline)
	if err != nil {
		// return the error with there's something wrong with the aggregation
		respondWithInternalError(c, err)
		return
	}
	// Parse the array of courses from these professors
	if err = cursor.All(ctx, &professorCourses); err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", professorCourses)
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
func ProfessorSectionSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		professorSection("Search", c)
	}
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
func ProfessorSectionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		professorSection("ById", c)
	}
}

// Get all of the sections of the professors depending on the type of flag
func professorSection(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var professorSections []schema.Section // array of sections of the professors (or single professor with Id)
	var professorQuery bson.M              // query filter the professor
	var err error

	// determine the professor's query
	professorQuery, err = getQuery[schema.Professor](flag, c)
	if err != nil {
		return
	}

	// determine the offset and limit for pagination stage
	paginateMap, err := configs.GetAggregateLimit(&professorQuery, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	// Pipeline to query the courses from the filtered professors (or a single professor)
	professorSectionPipeline := professorPipeline("sections", professorQuery, paginateMap)

	// Perform aggreration on the pipeline
	cursor, err := professorCollection.Aggregate(ctx, professorSectionPipeline)
	if err != nil {
		// return the error with there's something wrong with the aggregation
		respondWithInternalError(c, err)
		return
	}
	// Parse the array of sections from these professors
	if err = cursor.All(ctx, &professorSections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", professorSections)
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
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$sections"},
			{Key: "preserveNullAndEmptyArrays", Value: false}, // avoid course documents that can't be replaced
		}}},

		// Lookup course info using sections.course_reference
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "courses"},
			{Key: "localField", Value: "sections.course_reference"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections.course_details"},
		}}},

		// Lookup professor info using sections.course_reference
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "professors"},
			{Key: "localField", Value: "sections.professors"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections.professor_details"},
		}}},

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
