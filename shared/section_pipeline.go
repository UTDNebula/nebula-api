package shared

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// buildSectionPipeline builds the pipeline to aggregate courses or professors from filtered sections
func buildSectionPipeline(endpoint string, sectionQuery bson.M, paginate map[string]bson.D) mongo.Pipeline {
	baseStages := mongo.Pipeline{
		// filter the sections
		bson.D{{Key: "$match", Value: sectionQuery}},

		// paginate the sections before pulling courses from those sections
		paginate["former_offset"],
		paginate["limit"],
	}

	var lookupStages mongo.Pipeline
	switch endpoint {
	case "courses":
		lookupStages = mongo.Pipeline{
			// lookup the course referenced by sections from the course collection
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "courses"},
				{Key: "localField", Value: "course_reference"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "courses"},
			}}},
		}

	case "professors":
		lookupStages = mongo.Pipeline{
			// lookup the professors referenced by sections from the course collection
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "professors"},
				{Key: "localField", Value: "professors"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "professors"},
			}}},
		}

	default:
		panic("invalid endpoint for buildSectionPipeline: " + endpoint)
	}

	replaceStages := mongo.Pipeline{
		// project to remove every other fields except for courses/professors
		bson.D{{Key: "$project", Value: bson.D{
			{Key: endpoint, Value: "$" + endpoint},
		}}},

		// unwind the courses/professors
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + endpoint},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		// replace the combinations of id and course/professor with courses/professors entirely
		bson.D{{Key: "$replaceWith", Value: "$" + endpoint}},

		// remove duplicate courses/professors
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "item", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$replaceWith", Value: "$item"}},
	}

	paginateStages := mongo.Pipeline{
		// keep order deterministic between calls
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},

		// paginate the courses/professors
		paginate["latter_offset"],
		paginate["limit"],
	}

	return append(append(append(baseStages, lookupStages...), replaceStages...), paginateStages...)
}
