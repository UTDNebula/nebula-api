package loaders

import (
	"context"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type professorReader struct {
	collection *mongo.Collection
}

// getProfessors fetches professors by ID and returns them in the requested order.
func (p *professorReader) getProfessors(ctx context.Context, professorIDs []bson.ObjectID) ([]*model.Professor, []error) {
	cursor, err := p.collection.Find(ctx, bson.M{"_id": bson.M{"$in": professorIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	var professors []*model.Professor
	if err := cursor.All(ctx, &professors); err != nil {
		return nil, []error{err}
	}

	professorMap := make(map[bson.ObjectID]*model.Professor, len(professors))
	for _, professor := range professors {
		professorMap[professor.ID] = professor
	}

	orderedProfessors := make([]*model.Professor, len(professorIDs))
	errs := make([]error, len(professorIDs))
	for i, id := range professorIDs {
		orderedProfessors[i] = professorMap[id]
		if orderedProfessors[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return orderedProfessors, errs
}

// GetProfessor loads a professor from a reference that contains only its ID.
func GetProfessor(ctx context.Context, professorRef bson.ObjectID) (*model.Professor, error) {
	loaders := For(ctx)
	return loaders.ProfessorLoader.Load(ctx, professorRef)
}

// GetProfessors loads professors from references that contain only their IDs.
func GetProfessors(ctx context.Context, professorsRef []bson.ObjectID) ([]*model.Professor, error) {
  loaders := For(ctx)
	return loaders.ProfessorLoader.LoadAll(ctx, professorsRef)
}
