package shared

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func findRecord[T any](ctx context.Context, collection *mongo.Collection, id bson.ObjectID) (T, error) {
	var record T
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return record, ErrNotFound
	}
	return record, err
}
func findRecords[T any](ctx context.Context, collection *mongo.Collection, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var records []T
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}
func aggregateRecords[T any](ctx context.Context, collection *mongo.Collection, pipeline mongo.Pipeline) ([]T, error) {
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
