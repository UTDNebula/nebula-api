package loaders

import (
	"context"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type professorReader struct {
	collection *mongo.Collection
}

// getProfessors fetches professors by ID and returns them in the requested order.
func (p *professorReader) getProfessors(ctx context.Context, professorIDs []primitive.ObjectID) ([]*model.Professor, []error) {
	cursor, err := p.collection.Find(ctx, bson.M{"_id": bson.M{"$in": professorIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	var dbProfessors []*model.DBProfessor
	if err := cursor.All(ctx, &dbProfessors); err != nil {
		return nil, []error{err}
	}

	professorMap := make(map[primitive.ObjectID]*model.Professor, len(dbProfessors))
	for _, dbProfessor := range dbProfessors {
		professorMap[dbProfessor.ID] = model.TransformProfessor(dbProfessor)
	}

	professors := make([]*model.Professor, len(professorIDs))
	errs := make([]error, len(professorIDs))
	for i, id := range professorIDs {
		professors[i] = professorMap[id]
		if professors[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return professors, errs
}

// GetProfessor loads a professor from a reference that contains only its ID.
func GetProfessor(ctx context.Context, professorRef *model.Professor) (*model.Professor, error) {
	professorID, err := primitive.ObjectIDFromHex(professorRef.ID)
	if err != nil {
		return nil, err
	}

	return For(ctx).ProfessorLoader.Load(ctx, professorID)
}

// GetProfessors loads professors from references that contain only their IDs.
func GetProfessors(ctx context.Context, professorsRef []*model.Professor) ([]*model.Professor, error) {
	professorIDs := make([]primitive.ObjectID, len(professorsRef))
	for i, professor := range professorsRef {
		professorID, err := primitive.ObjectIDFromHex(professor.ID)
		if err != nil {
			return nil, err
		}
		professorIDs[i] = professorID
	}

	return For(ctx).ProfessorLoader.LoadAll(ctx, professorIDs)
}
