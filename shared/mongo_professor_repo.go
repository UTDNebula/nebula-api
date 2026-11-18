package shared

import (
	"context"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoProfessorRepo struct{ collection *mongo.Collection }

func NewMongoProfessorRepo(collection *mongo.Collection) *MongoProfessorRepo {
	return &MongoProfessorRepo{collection: collection}
}

var _ ProfessorRepo = (*MongoProfessorRepo)(nil)

func (r *MongoProfessorRepo) ByID(ctx context.Context, id bson.ObjectID) (schema.Professor, error) {
	return findRecord[schema.Professor](ctx, r.collection, id)
}
func (r *MongoProfessorRepo) Search(ctx context.Context, filter bson.M, offset, limit int64) ([]schema.Professor, error) {
	return findRecords[schema.Professor](ctx, r.collection, filter, options.Find().SetSkip(offset).SetLimit(limit))
}
func (r *MongoProfessorRepo) Grades(ctx context.Context, id bson.ObjectID) ([14]int, error) {
	return gradesForID(ctx, r.collection.Database(), "professor_endpoint", id)
}

func (r *MongoProfessorRepo) All(ctx context.Context) ([]schema.Professor, error) {
	return findRecords[schema.Professor](ctx, r.collection, bson.M{})
}
func (r *MongoProfessorRepo) TrendsSections(ctx context.Context, filter bson.M) ([]schema.Section, error) {
	return NewMongoTrendsRepo(r.collection.Database().Collection("trends_prof_sections")).Sections(ctx, filter)
}

func (r *MongoProfessorRepo) Courses(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Course, error) {
	return aggregateRecords[schema.Course](ctx, r.collection, buildProfessorPipeline("courses", filter, pagination))
}

func (r *MongoProfessorRepo) Sections(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Section, error) {
	return aggregateRecords[schema.Section](ctx, r.collection, buildProfessorPipeline("sections", filter, pagination))
}
