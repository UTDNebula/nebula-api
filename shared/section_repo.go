package shared

import (
	"context"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// SectionRepo supplies decoded coursebook records to API handlers.
type SectionRepo interface {
	ByID(context.Context, bson.ObjectID) (schema.Section, error)
	Search(context.Context, bson.M, int64, int64) ([]schema.Section, error)
	Grades(context.Context, bson.ObjectID) ([14]int, error)
	Courses(context.Context, bson.M, map[string]bson.D) ([]schema.Course, error)
	Professors(context.Context, bson.M, map[string]bson.D) ([]schema.Professor, error)
}
