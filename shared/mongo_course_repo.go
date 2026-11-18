package shared

import (
	"context"
	"errors"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoCourseRepo fetches course data from MongoDB.
type MongoCourseRepo struct {
	collection *mongo.Collection
}

// NewMongoCourseRepo receives an existing MongoDB collection.
func NewMongoCourseRepo(
	collection *mongo.Collection,
) *MongoCourseRepo {
	return &MongoCourseRepo{
		collection: collection,
	}
}

func (r *MongoCourseRepo) ByID(
	ctx context.Context,
	id bson.ObjectID,
) (schema.Course, error) {
	var course schema.Course

	filter := bson.M{
		"_id": id,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&course)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return schema.Course{}, ErrNotFound
	}

	if err != nil {
		return schema.Course{}, err
	}

	return course, nil
}

var _ CourseRepo = (*MongoCourseRepo)(nil)

func (r *MongoCourseRepo) Search(ctx context.Context, filter bson.M, offset, limit int64) ([]schema.Course, error) {
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var courses []schema.Course
	if err := cursor.All(ctx, &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *MongoCourseRepo) All(ctx context.Context) ([]schema.Course, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var courses []schema.Course
	if err := cursor.All(ctx, &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *MongoCourseRepo) Sections(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Section, error) {
	return aggregateRecords[schema.Section](ctx, r.collection, buildCoursePipeline("sections", filter, pagination))
}

func (r *MongoCourseRepo) Professors(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Professor, error) {
	return aggregateRecords[schema.Professor](ctx, r.collection, buildCoursePipeline("professors", filter, pagination))
}

func (r *MongoCourseRepo) Grades(ctx context.Context, id bson.ObjectID) ([14]int, error) {
	return gradesForID(ctx, r.collection.Database(), "course_endpoint", id)
}
func (r *MongoCourseRepo) TrendsSections(ctx context.Context, filter bson.M) ([]schema.Section, error) {
	return NewMongoTrendsRepo(r.collection.Database().Collection("trends_course_sections")).Sections(ctx, filter)
}
