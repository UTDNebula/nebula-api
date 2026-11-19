package loaders

import (
	"context"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type courseReader struct {
	collection *mongo.Collection
}

// getCourses implements a batch function that can retrieve many courses by ID for use in a dataloader
func (c *courseReader) getCourses(ctx context.Context, courseIDs []bson.ObjectID) ([]*model.Course, []error) {
	// Find documents matching courseIDs
	cursor, err := c.collection.Find(ctx, bson.M{"_id": bson.M{"$in": courseIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	// Decode
	var dbCourses []*model.DBCourse
	if err := cursor.All(ctx, &dbCourses); err != nil {
		return nil, []error{err}
	}

	// $in does not return in order of given CourseIDs
	// so we order with map.
	courseMap := make(map[bson.ObjectID]*model.Course, len(dbCourses))
	for _, dbCourse := range dbCourses {
		courseMap[dbCourse.ID] = model.TransformCourse(dbCourse)
	}

	courses := make([]*model.Course, len(courseIDs))
	for i, ID := range courseIDs {
		courses[i] = courseMap[ID]
	}

	return courses, nil
}

// For convenience in resolvers

// GetCourse returns single course by id efficiently
// Note that courseRef contains only its ObjectID.
// All other fields in courseRef are empty.
func GetCourse(ctx context.Context, courseRef *model.Course) (*model.Course, error) {
	courseID, err := bson.ObjectIDFromHex(courseRef.ID)
	if err != nil {
		return nil, err
	}

	loaders := For(ctx)
	return loaders.CourseLoader.Load(ctx, courseID)
}

// GetCourses returns many courses by ids efficiently
// Note that coursesRef contains only the courses' ObjectID.
// All other fields of each course in coursesRef are empty.
func GetCourses(ctx context.Context, coursesRef []*model.Course) ([]*model.Course, error) {
	courseIDs := make([]bson.ObjectID, len(coursesRef))
	for i, course := range coursesRef {
		courseID, err := bson.ObjectIDFromHex(course.ID)
		if err != nil {
			return nil, err
		}
		courseIDs[i] = courseID
	}

	loaders := For(ctx)
	return loaders.CourseLoader.LoadAll(ctx, courseIDs)
}
