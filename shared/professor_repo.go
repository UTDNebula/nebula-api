package shared

import (
	"context"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ProfessorRepo supplies decoded coursebook records to API handlers.
type ProfessorRepo interface {
	ByID(context.Context, bson.ObjectID) (schema.Professor, error)
	Search(context.Context, bson.M, int64, int64) ([]schema.Professor, error)
	All(context.Context) ([]schema.Professor, error)
	TrendsSections(context.Context, bson.M) ([]schema.Section, error)
	Grades(context.Context, bson.ObjectID) ([14]int, error)
	Courses(context.Context, bson.M, map[string]bson.D) ([]schema.Course, error)
	Sections(context.Context, bson.M, map[string]bson.D) ([]schema.Section, error)
}
