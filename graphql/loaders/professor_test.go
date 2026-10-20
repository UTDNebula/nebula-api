package loaders

import (
	"context"
	"testing"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"github.com/vikstrous/dataloadgen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetProfessorsInvalidID(t *testing.T) {
	// ID validation must finish before the loader can make a database query.
	refs := []*model.Professor{{ID: primitive.NewObjectID().Hex()}, {ID: "invalid"}}
	if _, err := GetProfessors(context.Background(), refs); err == nil {
		t.Fatal("expected invalid ID error")
	}
	if _, err := GetProfessor(context.Background(), refs[1]); err == nil {
		t.Fatal("expected invalid ID error")
	}
}

func TestGetProfessorsEmpty(t *testing.T) {
	loader := dataloadgen.NewLoader(func(context.Context, []primitive.ObjectID) ([]*model.Professor, []error) {
		t.Error("empty references must not query the database")
		return nil, nil
	})
	ctx := context.WithValue(context.Background(), loadersKey, &Loaders{ProfessorLoader: loader})
	professors, err := GetProfessors(ctx, nil)
	if err != nil || professors == nil || len(professors) != 0 {
		t.Fatalf("expected empty list, got %v, %v", professors, err)
	}
}
