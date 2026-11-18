package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
	trendsSectionSearch("Course", c)
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
	trendsSectionSearch("Professor", c)
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
func TrendsCombinedSectionSearch(c *gin.Context) {
	trendsSectionSearch("Combined", c)
}

// @Id				trendsGrades
// @Router			/grades/trends [post]
// @Tags			Grades
// @Description	"Returns grade distributions for a batch of courses, grouped by course, professor, semester, and section type. If professors are given, only the course-professor pairs that actually exist are returned; otherwise every professor of the given courses is returned. Specialized convenience endpoint for UTD Trends internal use."
// @Accept			json
// @Produce		json
// @Param			request	body		schema.TrendsGradesRequest						true	"Courses and optional professors"
// @Success		200		{object}	schema.APIResponse[[]schema.TrendsCourseGrades]	"Grade distributions grouped by course"
// @Failure		400		{object}	schema.APIResponse[string]						"A string describing the error"
// @Failure		500		{object}	schema.APIResponse[string]						"A string describing the error"
func TrendsGrades(c *gin.Context) {
	var req schema.TrendsGradesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	trendsQuery, err := trendsGradesQuery(req)
	if err != nil {
		respond(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	cursor, err := configs.GetCollection("trends_course_and_prof_sections").Aggregate(ctx, buildTrendsGradesPipeline(trendsQuery))
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	var courseGrades []schema.TrendsCourseGrades
	if err := cursor.All(ctx, &courseGrades); err != nil {
		respondWithInternalError(c, err)
		return
	}

	if courseGrades == nil {
		courseGrades = []schema.TrendsCourseGrades{}
	}
	for i := range courseGrades {
		courseGrades[i].Id = formatTrendsCourse(courseGrades[i].Id)
	}

	respond(c, http.StatusOK, "success", courseGrades)
}

// trendsGradesQuery builds a $match query over trends_course_and_prof_sections whose _id is
// {course, first_name, last_name}
// Both branches are bounded on _id so they are served by the default _id index.
func trendsGradesQuery(req schema.TrendsGradesRequest) (bson.M, error) {
	var courses []string
	seenCourses := map[string]bool{}
	for _, rawCourse := range req.Courses {

		course := strings.ToUpper(strings.Join(strings.Fields(rawCourse), ""))
		if course == "" {
			return nil, fmt.Errorf("course %q is empty", rawCourse)
		}
		if !seenCourses[course] {
			seenCourses[course] = true
			courses = append(courses, course)
		}
	}

	if len(req.Professors) == 0 {
		// Every _id with the given course sorts between {course: X} and {course: X + "\x00"}
		courseRanges := bson.A{}
		for _, course := range courses {
			courseRanges = append(courseRanges, bson.M{"_id": bson.D{
				{Key: "$gte", Value: bson.D{{Key: "course", Value: course}}},
				{Key: "$lt", Value: bson.D{{Key: "course", Value: course + "\x00"}}},
			}})
		}
		return bson.M{"$or": courseRanges}, nil
	}

	// Embedded _id documents only match exactly when the field order is identical, so bson.D is required
	ids := bson.A{}
	seenNames := map[[2]string]bool{}
	for _, professor := range req.Professors {
		words := strings.Fields(professor)
		if len(words) < 2 {
			return nil, fmt.Errorf("professor %q must include a first and last name", professor)
		}
		// Both first and last names may contain spaces, so try every split point
		for split := 1; split < len(words); split++ {
			name := [2]string{strings.Join(words[:split], " "), strings.Join(words[split:], " ")}
			if seenNames[name] {
				continue
			}
			seenNames[name] = true
			for _, course := range courses {
				ids = append(ids, bson.D{
					{Key: "course", Value: course},
					{Key: "first_name", Value: name[0]},
					{Key: "last_name", Value: name[1]},
				})
			}
		}
	}
	return bson.M{"_id": bson.M{"$in": ids}}, nil
}

// buildTrendsGradesPipeline sums the grade distributions of the matched sections per
// course, professor, semester, and section type, then nests them in that order
func buildTrendsGradesPipeline(trendsQuery bson.M) mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{Key: "$match", Value: trendsQuery}},

		bson.D{{Key: "$unwind", Value: "$sections"}},

		bson.D{{Key: "$match", Value: bson.D{
			{Key: "sections.grade_distribution.0", Value: bson.D{{Key: "$exists", Value: true}}},
		}}},

		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "course", Value: "$_id.course"},
				{Key: "professor", Value: bson.D{{Key: "$concat", Value: bson.A{"$_id.first_name", " ", "$_id.last_name"}}}},
				{Key: "semester", Value: "$sections.academic_session.name"},
				{Key: "type", Value: sectionTypeSwitch("$sections.section_number")},
			}},
			{Key: "grade_distributions", Value: bson.D{{Key: "$push", Value: "$sections.grade_distribution"}}},
		}}},

		// Element-wise sum of the grade distributions
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "grade_distribution", Value: bson.D{{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$grade_distributions"},
				{Key: "initialValue", Value: bson.A{}},
				{Key: "in", Value: bson.D{{Key: "$map", Value: bson.D{
					{Key: "input", Value: bson.D{{Key: "$zip", Value: bson.D{
						{Key: "inputs", Value: bson.A{"$$value", "$$this"}},
						{Key: "useLongestLength", Value: true},
						{Key: "defaults", Value: bson.A{0, 0}},
					}}}},
					{Key: "as", Value: "pair"},
					{Key: "in", Value: bson.D{{Key: "$sum", Value: "$$pair"}}},
				}}}},
			}}}},
		}}},

		// Nest types into semesters, semesters into professors, and professors into courses.
		// Each $sort keeps the pushed arrays ordered.
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id.course", Value: 1},
			{Key: "_id.professor", Value: 1},
			{Key: "_id.semester", Value: 1},
			{Key: "_id.type", Value: 1},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "course", Value: "$_id.course"},
				{Key: "professor", Value: "$_id.professor"},
				{Key: "semester", Value: "$_id.semester"},
			}},
			{Key: "types", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "type", Value: "$_id.type"},
				{Key: "grade_distribution", Value: "$grade_distribution"},
			}}}},
		}}},

		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id.course", Value: 1},
			{Key: "_id.professor", Value: 1},
			{Key: "_id.semester", Value: 1},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "course", Value: "$_id.course"},
				{Key: "professor", Value: "$_id.professor"},
			}},
			{Key: "semesters", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "semester", Value: "$_id.semester"},
				{Key: "types", Value: "$types"},
			}}}},
		}}},

		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id.course", Value: 1},
			{Key: "_id.professor", Value: 1},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.course"},
			{Key: "professors", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "professor", Value: "$_id.professor"},
				{Key: "semesters", Value: "$semesters"},
			}}}},
		}}},

		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}
}

// formatTrendsCourse turns a stored course key into the request format
func formatTrendsCourse(course string) string {
	if i := strings.IndexFunc(course, unicode.IsDigit); i > 0 {
		return course[:i] + " " + course[i:]
	}
	return course
}

// trendsSectionSearch handles trends-based section routes for both course and professor query.
// Reduce the repetitiveness of routes whose aggregation behaviors are identical.
// This is subject to change as requests might be more complex.
func trendsSectionSearch(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var detailedSections []schema.Section

	// If the branches become complicated, refactor to make it more generic.
	var trendsCollection *mongo.Collection
	var trendsQuery bson.M
	var err error

	switch flag {
	case "Course":
		trendsCollection = configs.GetCollection("trends_course_sections")
		// UTDTrends only uses prefix and number
		trendsQuery = bson.M{
			"_id": c.Query("subject_prefix") + c.Query("course_number"),
		}
	case "Professor":
		trendsCollection = configs.GetCollection("trends_prof_sections")
		trendsQuery, err = schema.FilterQuery[schema.Professor](c.Request.URL.Query())
		if err != nil {
			return
		}
	case "Combined":
		trendsCollection = configs.GetCollection("trends_course_and_prof_sections")
		trendsQuery = bson.M{
			"_id": bson.M{
				"course":     c.Query("subject_prefix") + c.Query("course_number"),
				"prof_first": c.Query("first_name"),
				"prof_last":  c.Query("last_name"),
			},
		}
	default:
		// This should never happen, but act as a fallback
		err = fmt.Errorf("invalid flag for trendsSectionSearch: %s", flag)
		respondWithInternalError(c, err)
	}

	cursor, err := trendsCollection.Aggregate(ctx, buildTrendsPipeline(trendsQuery))
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &detailedSections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", detailedSections)
}

// buildTrendsPipeline build the pipeline to embed the list of professors and course details
// directly in the queried list of sections
func buildTrendsPipeline(srcObjQuery bson.M) mongo.Pipeline {
	return mongo.Pipeline{
		// Match the original objects from the query
		bson.D{{Key: "$match", Value: srcObjQuery}},

		// Expand sections array into individual documents
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$sections"},
				// Avoid documents that can't be replaced
				{Key: "preserveNullAndEmptyArrays", Value: false},
			}},
		},

		// Embed the course details to the list of sections
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "sections.course_reference"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "sections.course_details"},
			}},
		},

		// Embed the professor details to the list of sections
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
