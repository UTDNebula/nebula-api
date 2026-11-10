package graph

import (
	"context"
	"time"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"github.com/UTDNebula/nebula-api/graphql/loaders"
	"github.com/UTDNebula/nebula-api/shared/configs"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *queryResolver) Professors(ctx context.Context, filter *model.ProfessorFilter, offset int32) ([]*model.Professor, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var professors []*model.Professor
	
	professorQuery := bson.M{}
	if filter != nil {
		bsonBytes, err := bson.Marshal(filter)
		if err != nil {
			return nil, err
		}
		if err = bson.Unmarshal(bsonBytes, &professorQuery); err != nil {
			return nil, err
		}
	} 

	paginate := options.Find().SetSkip(int64(offset)).SetLimit(configs.GetEnvLimit())

	cursor, err := r.ProfCollection.Find(timeoutCtx, professorQuery, paginate)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeoutCtx)

	if err = cursor.All(timeoutCtx, &professors); err != nil {
		return nil, err
	}

	return professors, err
}

func (r *queryResolver) Professor(ctx context.Context, id bson.ObjectID) (*model.Professor, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	professor := &model.Professor{}
	var err error

	result := r.ProfCollection.FindOne(timeoutCtx, bson.M{"_id": id})
	if err = result.Decode(professor); err != nil {
		return nil, err
	}

	return professor, err
}

// Sections is the resolver for the section field referenced from professor
func (r *professorResolver) Sections(ctx context.Context, obj *model.Professor) ([]*model.Section, error) {
	return loaders.GetSections(ctx, obj.Sections)
}
