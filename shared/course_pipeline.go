package shared

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// buildCoursePipeline builds the pipeline to aggregate the list of specified objects from list of courses
func buildCoursePipeline(endpoint string, courseQuery bson.M, paginate map[string]bson.D) mongo.Pipeline {
	baseStages := mongo.Pipeline{
		bson.D{{Key: "$match", Value: courseQuery}},

		// Skip to the offset, then limit to the number of courses
		paginate["former_offset"],
		paginate["limit"],

		// Lookup the list of sections from the courses
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},
	}

	var lookupStages, dedupStages mongo.Pipeline
	switch endpoint {
	case "sections":
		// No extra stages middle stages

	case "professors":
		// Lookup the list of professors from the list of sections
		lookupStages = mongo.Pipeline{
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "sections.professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "professors"},
			}}},
		}

		// Remove the duplicate professors
		dedupStages = mongo.Pipeline{
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "professors", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
			}}},

			bson.D{{Key: "$replaceWith", Value: "$professors"}},
		}

	default:
		panic("invalid endpoint for coursePipeline: " + endpoint)
	}

	replaceStages := mongo.Pipeline{
		// Unwind the target object of the sections
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + endpoint},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		// Replace the courses with the target objects
		bson.D{{Key: "$replaceWith", Value: "$" + endpoint}},
	}

	middleStages := append(append(lookupStages, replaceStages...), dedupStages...)

	paginateStages := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},

		paginate["latter_offset"],
		paginate["limit"],
	}

	return append(append(baseStages, middleStages...), paginateStages...)
}
