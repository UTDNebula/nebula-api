package graph

import (
	"context"
	"time"

	"graphql/configs"
	"graphql/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *queryResolver) Professors(ctx context.Context, filter *model.ProfessorFilter, offset *int32) ([]*model.Professor, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var professors []*model.Professor
	var dbProfessors []*model.DBProfessor
	var err error

	var professorQuery bson.M
	if filter != nil {
		bsonBytes, err := bson.Marshal(filter)
		if err != nil {
			return nil, err
		}
		if err = bson.Unmarshal(bsonBytes, &professorQuery); err != nil {
			return nil, err
		}
	} else {
		professorQuery = bson.M{} // no filter
	}

	paginate := options.Find().SetSkip(int64(*offset)).SetLimit(configs.GetEnvLimit())

	cursor, err := r.ProfCollection.Find(timeoutCtx, professorQuery, paginate)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeoutCtx)

	if err = cursor.All(timeoutCtx, &dbProfessors); err != nil {
		return nil, err
	}

	for _, dbProfessor := range dbProfessors {
		professors = append(professors, model.TransformProfessor(dbProfessor))
	}

	return professors, err
}

func (r *queryResolver) Professor(ctx context.Context, id string) (*model.Professor, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	dbProfessor := &model.DBProfessor{}

	err = r.ProfCollection.FindOne(
		timeoutCtx, bson.M{"_id": objectId}).Decode(dbProfessor)
	if err != nil {
		return nil, err
	}

	return model.TransformProfessor(dbProfessor), err

}
