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
				{"$lookup",
					bson.D{
						{"from", "sections"},
						{"localField", "sections"},
						{"foreignField", "_id"},
						{"as", "sections"},
					},
				},
			},
			bson.D{{"$unwind", "$sections"}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "professors"},
						{"localField", "sections.professors"},
						{"foreignField", "_id"},
						{"as", "professors"},
					},
				},
			},
			bson.D{
				{"$group",
					bson.D{
						{"_id", "$sections._id"},
						{"grade_distribution", bson.D{{"$push", "$sections.grade_distribution"}}},
					},
				},
			},
		}

		sumSemesterPipeline :=
			mongo.Pipeline{
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "sections"},
							{"localField", "sections"},
							{"foreignField", "_id"},
							{"as", "sections"},
						},
					},
				},
				bson.D{{"$unwind", "$sections"}},
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "professors"},
							{"localField", "sections.professors"},
							{"foreignField", "_id"},
							{"as", "professors"},
						},
					},
				},
				bson.D{
					{"$project",
						bson.D{
							{"_id", "$sections.academic_session.name"},
							{"grade_distribution", "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{"$unwind",
						bson.D{
							{"path", "$grade_distribution"},
							{"includeArrayIndex", "ix"},
						},
					},
				},
				bson.D{
					{"$group",
						bson.D{
							{"_id",
								bson.D{
									{"academic_session", "$_id"},
									{"ix", "$ix"},
								},
							},
							{"grade_distributions", bson.D{{"$push", "$grade_distribution"}}},
						},
					},
				},
				bson.D{
					{"$sort",
						bson.D{
							{"_id.ix", 1},
							{"_id", 1},
						},
					},
				},
				bson.D{{"$addFields", bson.D{{"grade_distributions", bson.D{{"$sum", "$grade_distributions"}}}}}},
				bson.D{
					{"$group",
						bson.D{
							{"_id", "$_id.academic_session"},
							{"grade_distribution", bson.D{{"$push", "$grade_distributions"}}},
						},
					},
				},
			}

		totalPipeline :=
			mongo.Pipeline{
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "sections"},
							{"localField", "sections"},
							{"foreignField", "_id"},
							{"as", "sections"},
						},
					},
				},
				bson.D{{"$unwind", "$sections"}},
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "professors"},
							{"localField", "sections.professors"},
							{"foreignField", "_id"},
							{"as", "professors"},
						},
					},
				},
				bson.D{
					{"$project",
						bson.D{
							{"_id", primitive.Null{}},
							{"grade_distribution", "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{"$unwind",
						bson.D{
							{"path", "$grade_distribution"},
							{"includeArrayIndex", "ix"},
						},
					},
				},
				bson.D{
					{"$group",
						bson.D{
							{"_id", "$ix"},
							{"grade_distribution", bson.D{{"$push", "$grade_distribution"}}},
						},
					},
				},
				bson.D{{"$sort", bson.D{{"_id", 1}}}},
				bson.D{{"$addFields", bson.D{{"grade_distribution", bson.D{{"$sum", "$grade_distribution"}}}}}},
				bson.D{
					{"$group",
						bson.D{
							{"_id", primitive.Null{}},
							{"grade_distribution", bson.D{{"$push", "$grade_distribution"}}},
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

		for key, _ := range queryParams {
			if key[0:len("courses.")] == "courses." { // discard courses prefix becuase pipeline will aggregate courses
				courseKey := key[len("courses."):]
				query[courseKey] = c.Query(key)
			} else if key == "representation" {
				representation = c.Query(key)
			} else {
				query[key] = c.Query(key)
			}
		}

		bySectionPipeline := mongo.Pipeline{
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "sections"},
						{"localField", "sections"},
						{"foreignField", "_id"},
						{"as", "sections"},
					},
				},
			},
			bson.D{{"$unwind", "$sections"}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "professors"},
						{"localField", "sections.professors"},
						{"foreignField", "_id"},
						{"as", "professors"},
					},
				},
			},
			bson.D{{"$match", query}},
			bson.D{
				{"$group",
					bson.D{
						{"_id", "$sections._id"},
						{"grade_distribution", bson.D{{"$push", "$sections.grade_distribution"}}},
					},
				},
			},
		}

		totalPipeline :=
			mongo.Pipeline{
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "sections"},
							{"localField", "sections"},
							{"foreignField", "_id"},
							{"as", "sections"},
						},
					},
				},
				bson.D{{"$unwind", "$sections"}},
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "professors"},
							{"localField", "sections.professors"},
							{"foreignField", "_id"},
							{"as", "professors"},
						},
					},
				},
				bson.D{{"$match", query}},
				bson.D{
					{"$project",
						bson.D{
							{"_id", primitive.Null{}},
							{"grade_distribution", "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{"$unwind",
						bson.D{
							{"path", "$grade_distribution"},
							{"includeArrayIndex", "ix"},
						},
					},
				},
				bson.D{
					{"$group",
						bson.D{
							{"_id", "$ix"},
							{"grade_distribution", bson.D{{"$push", "$grade_distribution"}}},
						},
					},
				},
				bson.D{{"$sort", bson.D{{"_id", 1}}}},
				bson.D{{"$addFields", bson.D{{"grade_distribution", bson.D{{"$sum", "$grade_distribution"}}}}}},
				bson.D{
					{"$group",
						bson.D{
							{"_id", primitive.Null{}},
							{"grade_distribution", bson.D{{"$push", "$grade_distribution"}}},
						},
					},
				},
			}

		sumSemesterPipeline :=
			mongo.Pipeline{
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "sections"},
							{"localField", "sections"},
							{"foreignField", "_id"},
							{"as", "sections"},
						},
					},
				},
				bson.D{{"$unwind", "$sections"}},
				bson.D{
					{"$lookup",
						bson.D{
							{"from", "professors"},
							{"localField", "sections.professors"},
							{"foreignField", "_id"},
							{"as", "professors"},
						},
					},
				},
				bson.D{
					{"$match", query},
				},
				bson.D{
					{"$project",
						bson.D{
							{"_id", "$sections.academic_session.name"},
							{"grade_distribution", "$sections.grade_distribution"},
						},
					},
				},
				bson.D{
					{"$unwind",
						bson.D{
							{"path", "$grade_distribution"},
							{"includeArrayIndex", "ix"},
						},
					},
				},
				bson.D{
					{"$group",
						bson.D{
							{"_id",
								bson.D{
									{"academic_session", "$_id"},
									{"ix", "$ix"},
								},
							},
							{"grade_distributions", bson.D{{"$push", "$grade_distribution"}}},
						},
					},
				},
				bson.D{
					{"$sort",
						bson.D{
							{"_id.ix", 1},
							{"_id", 1},
						},
					},
				},
				bson.D{{"$addFields", bson.D{{"grade_distributions", bson.D{{"$sum", "$grade_distributions"}}}}}},
				bson.D{
					{"$group",
						bson.D{
							{"_id", "$_id.academic_session"},
							{"grade_distribution", bson.D{{"$push", "$grade_distributions"}}},
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
