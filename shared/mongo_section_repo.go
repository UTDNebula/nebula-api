package shared

import (
	"context"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoSectionRepo struct{ collection *mongo.Collection }

func NewMongoSectionRepo(collection *mongo.Collection) *MongoSectionRepo {
	return &MongoSectionRepo{collection: collection}
}

var _ SectionRepo = (*MongoSectionRepo)(nil)

func (r *MongoSectionRepo) ByID(ctx context.Context, id bson.ObjectID) (schema.Section, error) {
	return findRecord[schema.Section](ctx, r.collection, id)
}
func (r *MongoSectionRepo) Search(ctx context.Context, filter bson.M, offset, limit int64) ([]schema.Section, error) {
	return findRecords[schema.Section](ctx, r.collection, filter, options.Find().SetSkip(offset).SetLimit(limit))
}
func (r *MongoSectionRepo) Grades(ctx context.Context, id bson.ObjectID) ([14]int, error) {
	return gradesForID(ctx, r.collection.Database(), "section_endpoint", id)
}

func (r *MongoSectionRepo) Courses(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Course, error) {
	return aggregateRecords[schema.Course](ctx, r.collection, buildSectionPipeline("courses", filter, pagination))
}

func (r *MongoSectionRepo) Professors(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Professor, error) {
	return aggregateRecords[schema.Professor](ctx, r.collection, buildSectionPipeline("professors", filter, pagination))
}
