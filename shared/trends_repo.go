package shared

import (
	"context"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TrendsRepo interface {
	Sections(context.Context, bson.M) ([]schema.Section, error)
}
type MongoTrendsRepo struct{ collection *mongo.Collection }

func NewMongoTrendsRepo(collection *mongo.Collection) *MongoTrendsRepo {
	return &MongoTrendsRepo{collection: collection}
}
func (r *MongoTrendsRepo) Sections(ctx context.Context, filter bson.M) ([]schema.Section, error) {
	return aggregateRecords[schema.Section](ctx, r.collection, buildTrendsPipeline(filter))
}

var _ TrendsRepo = (*MongoTrendsRepo)(nil)
