package loaders

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"github.com/vikstrous/dataloadgen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestProfessorLoader(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	firstID, secondID := primitive.NewObjectID(), primitive.NewObjectID()

	mt.Run("batch order duplicates and cache", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, mt.DB.Name()+"."+mt.Coll.Name(), mtest.FirstBatch,
			bson.D{{Key: "_id", Value: secondID}, {Key: "first_name", Value: "Second"}},
			bson.D{{Key: "_id", Value: firstID}, {Key: "first_name", Value: "First"}},
		))
		reader := &professorReader{collection: mt.Coll}
		loader := dataloadgen.NewLoader(reader.getProfessors, dataloadgen.WithWait(time.Millisecond))
		ctx := context.WithValue(context.Background(), loadersKey, &Loaders{ProfessorLoader: loader})
		refs := []*model.Professor{{ID: firstID.Hex()}, {ID: secondID.Hex()}, {ID: firstID.Hex()}}

		professors, err := GetProfessors(ctx, refs)
		if err != nil {
			mt.Fatal(err)
		}
		for i, ref := range refs {
			if professors[i] == nil || professors[i].ID != ref.ID {
				mt.Fatalf("professor %d = %v, want ID %s", i, professors[i], ref.ID)
			}
		}
		if professors[0].FirstName != "First" || professors[1].FirstName != "Second" {
			mt.Fatalf("unexpected professor data: %v", professors)
		}
		cached, err := GetProfessor(ctx, refs[0])
		if err != nil || cached != professors[0] {
			mt.Fatalf("cached professor = %v, error = %v", cached, err)
		}

		commands := mt.GetAllStartedEvents()
		if len(commands) != 1 || commands[0].CommandName != "find" {
			mt.Fatalf("expected one find command, got %v", commands)
		}
		values, err := commands[0].Command.Lookup("filter", "_id", "$in").Array().Values()
		if err != nil {
			mt.Fatal(err)
		}
		ids := make([]primitive.ObjectID, len(values))
		for i, value := range values {
			ids[i] = value.ObjectID()
		}
		if !reflect.DeepEqual(ids, []primitive.ObjectID{firstID, secondID}) {
			mt.Fatalf("query IDs = %v, want two unique IDs", ids)
		}
	})

	mt.Run("missing professor keeps error at its position", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, mt.DB.Name()+"."+mt.Coll.Name(), mtest.FirstBatch,
			bson.D{{Key: "_id", Value: secondID}},
		))
		reader := &professorReader{collection: mt.Coll}
		professors, errs := reader.getProfessors(context.Background(), []primitive.ObjectID{firstID, secondID})
		if len(professors) != 2 || professors[0] != nil || professors[1].ID != secondID.Hex() {
			mt.Fatalf("unexpected results: %v", professors)
		}
		if len(errs) != 2 || !errors.Is(errs[0], mongo.ErrNoDocuments) || errs[1] != nil {
			mt.Fatalf("unexpected errors: %v", errs)
		}
	})

	mt.Run("sections share a batch", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, mt.DB.Name()+"."+mt.Coll.Name(), mtest.FirstBatch,
			bson.D{{Key: "_id", Value: firstID}},
			bson.D{{Key: "_id", Value: secondID}},
		))
		reader := &professorReader{collection: mt.Coll}
		loader := dataloadgen.NewLoader(reader.getProfessors, dataloadgen.WithWait(time.Second))
		ctx := context.WithValue(context.Background(), loadersKey, &Loaders{ProfessorLoader: loader})
		sections := []*model.Section{
			{Professors: []*model.Professor{{ID: firstID.Hex()}}},
			{Professors: []*model.Professor{{ID: secondID.Hex()}, {ID: firstID.Hex()}}},
		}
		var wg sync.WaitGroup
		for _, section := range sections {
			wg.Go(func() {
				professors, err := GetProfessors(ctx, section.Professors)
				if err != nil {
					mt.Error(err)
					return
				}
				for i, ref := range section.Professors {
					if professors[i] == nil || professors[i].ID != ref.ID {
						mt.Errorf("professor %d = %v, want ID %s", i, professors[i], ref.ID)
					}
				}
			})
		}
		wg.Wait()
		if commands := mt.GetAllStartedEvents(); len(commands) != 1 || commands[0].CommandName != "find" {
			mt.Fatalf("expected one professor query across sections, got %v", commands)
		}
	})

	mt.Run("database error reaches every key", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 2, Message: "query failed"}))
		reader := &professorReader{collection: mt.Coll}
		loader := dataloadgen.NewLoader(reader.getProfessors)
		_, err := loader.LoadAll(context.Background(), []primitive.ObjectID{firstID, secondID})
		var errs dataloadgen.ErrorSlice
		if !errors.As(err, &errs) || len(errs) != 2 {
			mt.Fatalf("expected two load errors, got %v", err)
		}
		for _, err := range errs {
			var commandErr mongo.CommandError
			if !errors.As(err, &commandErr) || commandErr.Code != 2 {
				mt.Fatalf("expected database error, got %v", err)
			}
		}
	})

	mt.Run("decode error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, mt.DB.Name()+"."+mt.Coll.Name(), mtest.FirstBatch,
			bson.D{{Key: "_id", Value: firstID}, {Key: "sections", Value: bson.A{42}}},
		))
		reader := &professorReader{collection: mt.Coll}
		professors, errs := reader.getProfessors(context.Background(), []primitive.ObjectID{firstID})
		if professors != nil || len(errs) != 1 || errs[0] == nil {
			mt.Fatalf("expected decode error, got %v, %v", professors, errs)
		}
	})
}

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
