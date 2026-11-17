package shared

import (
	"context"
	"errors"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ErrEventsNotFound means that no event document exists for the requested date.
var ErrEventsNotFound = errors.New("no events found for date")

// EventRepository retrieves a date's events without HTTP or database driver results.
// Implementations return ErrEventsNotFound when the date has no document.
type EventRepository[T any] interface {
	GetEventsByDate(context.Context, string) (schema.MultiBuildingEvents[T], error)
}

type mongoEventRepository[T any] struct {
	collection *mongo.Collection
}

// NewMongoEventRepository uses an existing collection. It does not connect to MongoDB.
func NewMongoEventRepository[T any](collection *mongo.Collection) EventRepository[T] {
	return &mongoEventRepository[T]{collection: collection}
}

func (r *mongoEventRepository[T]) GetEventsByDate(ctx context.Context, date string) (schema.MultiBuildingEvents[T], error) {
	var events schema.MultiBuildingEvents[T]
	err := r.collection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return events, ErrEventsNotFound
	}
	return events, err
}
