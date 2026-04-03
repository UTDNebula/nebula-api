package graph

import (
	"context"
	"graphql/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Rooms is the resolver for the rooms field.
func (r *queryResolver) Rooms(ctx context.Context) ([]*model.BuildingRooms, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var buildingRooms []*model.BuildingRooms
	var err error

	cursor, err := r.BuildingCollection.Find(timeoutCtx, bson.M{})
	if err != nil {
		return buildingRooms, err
	}
	defer cursor.Close(timeoutCtx)

	if err = cursor.All(ctx, &buildingRooms); err != nil {
		return buildingRooms, err
	}

	return buildingRooms, err
}
