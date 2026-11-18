package shared

import (
	"context"
	"errors"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrInvalidGradeQuery = errors.New("Invalid query parameters.")

// GradeQuery is independent of HTTP, so both REST and GraphQL can use it.
type GradeQuery struct {
	Mode          string
	ID            *bson.ObjectID
	Prefix        string
	Number        string
	SectionNumber string
	FirstName     string
	LastName      string
}
type GradeResult struct {
	Overall      [14]int
	Semesters    []schema.GradeData
	SectionTypes []schema.TypedGradeData
}
type GradesRepo interface {
	Aggregate(context.Context, GradeQuery) (GradeResult, error)
}
type MongoGradesRepo struct {
	courses    *mongo.Collection
	sections   *mongo.Collection
	professors *mongo.Collection
}

func NewMongoGradesRepo(database *mongo.Database) *MongoGradesRepo {
	return &MongoGradesRepo{
		courses:    database.Collection("courses"),
		sections:   database.Collection("sections"),
		professors: database.Collection("professors"),
	}
}

var _ GradesRepo = (*MongoGradesRepo)(nil)

func gradesForID(ctx context.Context, database *mongo.Database, mode string, id bson.ObjectID) ([14]int, error) {
	result, err := NewMongoGradesRepo(database).Aggregate(ctx, GradeQuery{Mode: mode, ID: &id})
	return result.Overall, err
}
