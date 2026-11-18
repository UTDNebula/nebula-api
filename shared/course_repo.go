package shared

import (
	"context"
	"errors"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrNotFound = errors.New("record not found")

// CourseRepo describes the course operations a controller can request.
type CourseRepo interface {
	ByID(
		ctx context.Context,
		id bson.ObjectID,
	) (schema.Course, error)
	Search(ctx context.Context, filter bson.M, offset, limit int64) ([]schema.Course, error)
	All(ctx context.Context) ([]schema.Course, error)
	Sections(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Section, error)
	Professors(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Professor, error)
	Grades(ctx context.Context, id bson.ObjectID) ([14]int, error)
	TrendsSections(ctx context.Context, filter bson.M) ([]schema.Section, error)
}
