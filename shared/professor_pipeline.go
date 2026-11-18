package shared

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Pipeline builder for professor aggregate endpoints
func buildProfessorPipeline(endpoint string, professorQuery bson.M, paginate map[string]bson.D) mongo.Pipeline {
	// common stages
	baseStages := mongo.Pipeline{
		// filter the professors
		bson.D{{Key: "$match", Value: professorQuery}},

		// paginate the professors before pulling the courses/sections from those professor
		paginate["former_offset"],
		paginate["limit"],

		// lookup the array of sections from sections collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},
	}

	var middleStages mongo.Pipeline

	switch endpoint {
	case "courses":
		middleStages = mongo.Pipeline{
			// project the courses referenced by each section in the array
			bson.D{{Key: "$project", Value: bson.D{{Key: "courses", Value: "$sections.course_reference"}}}},

			// lookup the array of courses from coures collection
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "courses"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "courses"},
			}}},
		}

	case "sections":
		middleStages = mongo.Pipeline{
			// project the sections
			bson.D{{Key: "$project", Value: bson.D{{Key: "sections", Value: "$sections"}}}},
		}

	default:
		panic("invalid endpoint for professorPipeline: " + endpoint)
	}

	// common pagination stages
	paginateStages := mongo.Pipeline{
		// unwind the courses/sections
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + endpoint},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		// replace the combination of ids and courses/sections with the courses/sections entirely
		bson.D{{Key: "$replaceWith", Value: "$" + endpoint}},

		// keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},

		// paginate the courses/sections
		paginate["latter_offset"],
		paginate["limit"],
	}

	return append(append(baseStages, middleStages...), paginateStages...)
}
