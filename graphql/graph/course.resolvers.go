package graph

import (
	"context"
	"time"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"github.com/UTDNebula/nebula-api/graphql/loaders"
	"github.com/UTDNebula/nebula-api/shared/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Courses is the resolver for the courses field.
func (r *queryResolver) Courses(ctx context.Context, filter *model.CourseFilter, offset int32) ([]*model.Course, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var courses []*model.Course
	var dbCourses []*model.DBCourse

	// Build the mongo query from the filter's struct
	var courseQuery bson.M
	if filter != nil {
		bsonBytes, err := bson.Marshal(filter)
		if err != nil {
			return nil, err
		}
		if err = bson.Unmarshal(bsonBytes, &courseQuery); err != nil {
			return nil, err
		}
	}
	// Paginate the list of courses
	paginate := options.Find().SetSkip(int64(offset)).SetLimit(configs.GetEnvLimit())

	// Query from Database
	cursor, err := r.CourseCollection.Find(timeoutCtx, courseQuery, paginate)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeoutCtx)

	if err = cursor.All(timeoutCtx, &dbCourses); err != nil {
		return nil, err
	}

	// Transform Mongo course to GraphQL course
	for _, dbCourse := range dbCourses {
		courses = append(courses, model.TransformCourse(dbCourse))
	}

	return courses, err
}

// Course is the resolver for the course field.
func (r *queryResolver) Course(ctx context.Context, id string) (*model.Course, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var dbCourse model.DBCourse
	result := r.CourseCollection.FindOne(timeoutCtx, bson.M{"_id": objectId})
	if err = result.Decode(&dbCourse); err != nil {
		return nil, err
	}

	return model.TransformCourse(&dbCourse), err
}

// Sections is the resolver for the section field referenced from course
func (r *courseResolver) Sections(ctx context.Context, obj *model.Course) ([]*model.Section, error) {
	return loaders.GetSections(ctx, obj.Sections)
}
