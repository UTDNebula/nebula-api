package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var coursesCollection *mongo.Collection = configs.GetCollection(configs.DB, "courses")

func GradesAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		representation := c.Request.URL.Query().Get("representation")

		var grades []map[string]interface{}

		bySectionPipeline := mongo.Pipeline{
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "sections"},
						{Key: "localField", Value: "sections"},
						{Key: "foreignField", Value: "_id"},
						{Key: "as", Value: "sections"},
					},
				},
			},
			bson.D{{Key: "$unwind", Value: "$sections"}},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "professors"},
						{Key: "localField", Value: "sections.professors"},
						{Key: "foreignField", Value: "_id"},
						{Key: "as", Value: "professors"},
					},
				},
			},
			bson.D{
				{Key: "$group",
					Value: bson.D{
						{Key: "_id", Value: "$sections._id"},
						{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$sections.grade_distribution"}}},
					},
				},
			},
		}

		sumSemesterPipeline :=
			mongo.Pipeline{
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "sections"},
							{Key: "localField", Value: "sections"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "sections"},
						},
					},
				},
				bson.D{{Key: "$unwind", Value: "$sections"}},
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "professors"},
							{Key: "localField", Value: "sections.professors"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "professors"},
						},
					},
				},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "_id", Value: "$sections.academic_session.name"},
							{Key: "grade_distribution", Value: "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{Key: "$unwind",
						Value: bson.D{
							{Key: "path", Value: "$grade_distribution"},
							{Key: "includeArrayIndex", Value: "ix"},
						},
					},
				},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id",
								Value: bson.D{
									{Key: "academic_session", Value: "$_id"},
									{Key: "ix", Value: "$ix"},
								},
							},
							{Key: "grade_distributions", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
						},
					},
				},
				bson.D{
					{Key: "$sort",
						Value: bson.D{
							{Key: "_id.ix", Value: 1},
							{Key: "_id", Value: 1},
						},
					},
				},
				bson.D{{Key: "$addFields", Value: bson.D{{Key: "grade_distributions", Value: bson.D{{Key: "$sum", Value: "$grade_distributions"}}}}}},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id", Value: "$_id.academic_session"},
							{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grade_distributions"}}},
						},
					},
				},
			}

		totalPipeline :=
			mongo.Pipeline{
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "sections"},
							{Key: "localField", Value: "sections"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "sections"},
						},
					},
				},
				bson.D{{Key: "$unwind", Value: "$sections"}},
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "professors"},
							{Key: "localField", Value: "sections.professors"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "professors"},
						},
					},
				},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "_id", Value: primitive.Null{}},
							{Key: "grade_distribution", Value: "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{Key: "$unwind",
						Value: bson.D{
							{Key: "path", Value: "$grade_distribution"},
							{Key: "includeArrayIndex", Value: "ix"},
						},
					},
				},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id", Value: "$ix"},
							{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
						},
					},
				},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
				bson.D{{Key: "$addFields", Value: bson.D{{Key: "grade_distribution", Value: bson.D{{Key: "$sum", Value: "$grade_distribution"}}}}}},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id", Value: primitive.Null{}},
							{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
						},
					},
				},
			}

		var cursor *mongo.Cursor
		var err error
		switch representation {
		case "section":
			cursor, err = courseCollection.Aggregate(ctx, bySectionPipeline)
		case "semester":
			cursor, err = courseCollection.Aggregate(ctx, sumSemesterPipeline)
		case "total":
			cursor, err = courseCollection.Aggregate(ctx, totalPipeline)
		default:
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "invalid representation field"})
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &grades); err != nil {
			panic(err)
		}

		// return result
		c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: grades})
	}
}

func GradesSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		var representation string
		query := bson.M{}
		var grades []map[string]interface{}

		queryParams := c.Request.URL.Query()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Query takes the following form:
		//  - If filtering by course parameter, use equivalent CourseSearch parameter.
		//  - If filtering by section parameter, add the prefix "sections." to the equivalent SectionSearch query parameter.
		//  - If filtering by professor parameter, add the prefix "professors." to the equivalent ProfessorSearch query parameter.

		for key, _ := range queryParams {
			if key[0:len("courses.")] == "courses." { // discard courses prefix becuase pipeline will aggregate courses
				courseKey := key[len("courses."):]
				query[courseKey] = c.Query(key)
			} else if key == "representation" {
				representation = c.Query(key)
			} else if key == "sections.course_reference" || key == "sections.professors" {
				objId, err := primitive.ObjectIDFromHex(c.Query(key))
				if err != nil {
					c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
					return
				} else {
					query[key] = objId
				}
			} else {
				query[key] = c.Query(key)
			}
		}

		bySectionPipeline := mongo.Pipeline{
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "sections"},
						{Key: "localField", Value: "sections"},
						{Key: "foreignField", Value: "_id"},
						{Key: "as", Value: "sections"},
					},
				},
			},
			bson.D{{Key: "$unwind", Value: "$sections"}},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "professors"},
						{Key: "localField", Value: "sections.professors"},
						{Key: "foreignField", Value: "_id"},
						{Key: "as", Value: "professors"},
					},
				},
			},
			bson.D{{Key: "$match", Value: query}},
			bson.D{
				{Key: "$group",
					Value: bson.D{
						{Key: "_id", Value: "$sections._id"},
						{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$sections.grade_distribution"}}},
					},
				},
			},
		}

		totalPipeline :=
			mongo.Pipeline{
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "sections"},
							{Key: "localField", Value: "sections"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "sections"},
						},
					},
				},
				bson.D{{Key: "$unwind", Value: "$sections"}},
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "professors"},
							{Key: "localField", Value: "sections.professors"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "professors"},
						},
					},
				},
				bson.D{{Key: "$match", Value: query}},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "_id", Value: primitive.Null{}},
							{Key: "grade_distribution", Value: "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{Key: "$unwind",
						Value: bson.D{
							{Key: "path", Value: "$grade_distribution"},
							{Key: "includeArrayIndex", Value: "ix"},
						},
					},
				},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id", Value: "$ix"},
							{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
						},
					},
				},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
				bson.D{{Key: "$addFields", Value: bson.D{{Key: "grade_distribution", Value: bson.D{{Key: "$sum", Value: "$grade_distribution"}}}}}},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id", Value: primitive.Null{}},
							{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
						},
					},
				},
			}

		sumSemesterPipeline :=
			mongo.Pipeline{
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "sections"},
							{Key: "localField", Value: "sections"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "sections"},
						},
					},
				},
				bson.D{{Key: "$unwind", Value: "$sections"}},
				bson.D{
					{Key: "$lookup",
						Value: bson.D{
							{Key: "from", Value: "professors"},
							{Key: "localField", Value: "sections.professors"},
							{Key: "foreignField", Value: "_id"},
							{Key: "as", Value: "professors"},
						},
					},
				},
				bson.D{
					{Key: "$match", Value: query},
				},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "_id", Value: "$sections.academic_session.name"},
							{Key: "grade_distribution", Value: "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{Key: "$unwind",
						Value: bson.D{
							{Key: "path", Value: "$grade_distribution"},
							{Key: "includeArrayIndex", Value: "ix"},
						},
					},
				},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id",
								Value: bson.D{
									{Key: "academic_session", Value: "$_id"},
									{Key: "ix", Value: "$ix"},
								},
							},
							{Key: "grade_distributions", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
						},
					},
				},
				bson.D{
					{Key: "$sort",
						Value: bson.D{
							{Key: "_id.ix", Value: 1},
							{Key: "_id", Value: 1},
						},
					},
				},
				bson.D{{Key: "$addFields", Value: bson.D{{Key: "grade_distributions", Value: bson.D{{Key: "$sum", Value: "$grade_distributions"}}}}}},
				bson.D{
					{Key: "$group",
						Value: bson.D{
							{Key: "_id", Value: "$_id.academic_session"},
							{Key: "grade_distribution", Value: bson.D{{Key: "$push", Value: "$grade_distributions"}}},
						},
					},
				},
			}

		var cursor *mongo.Cursor
		var err error

		switch representation {
		case "section":
			cursor, err = coursesCollection.Aggregate(ctx, bySectionPipeline)
		case "semester":
			cursor, err = courseCollection.Aggregate(ctx, sumSemesterPipeline)
		case "total":
			cursor, err = coursesCollection.Aggregate(ctx, totalPipeline)
		default:
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "invalid representation field"})
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProfessorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		// retrieve and parse all valid documents
		if err = cursor.All(ctx, &grades); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, responses.DegreeResponse{Status: http.StatusOK, Message: "success", Data: grades})
	}
}
