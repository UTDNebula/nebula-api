package shared

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

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
