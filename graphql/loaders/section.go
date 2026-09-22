package loaders

import (
	"context"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type sectionReader struct {
	collection *mongo.Collection 
}

// getSections implements a batch function that can retrieve many users by ID for use in a dataloader
func (s *sectionReader) getSections(ctx context.Context, sectionIDs []primitive.ObjectID) ([]*model.Section, []error) {
	// Find documents matching sectionIDs
	cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": sectionIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	// Decode
	var dbSections []*model.DBSection
	if err := cursor.All(ctx, &dbSections); err != nil {
		return nil, []error{err}
	}

	// $in does not return in order of given sectionIDs
	// so we order with map.
	sectionMap := make(map[primitive.ObjectID]*model.Section, len(dbSections))
	for _, dbSection := range dbSections {
			sectionMap[dbSection.ID] = model.TransformSection(dbSection)
	}

	sections := make([]*model.Section, len(sectionIDs))
	for i, ID := range sectionIDs {
			sections[i] = sectionMap[ID]
	}

	return sections, nil
}

// For convenience in resolvers

// GetSection returns single user by id efficiently
func GetSection(ctx context.Context, sectionID primitive.ObjectID) (*model.Section, error) {
	loaders := For(ctx)
	return loaders.SectionLoader.Load(ctx, sectionID)
}

// GetSections returns many users by ids efficiently
func GetSections(ctx context.Context, sectionIDs []primitive.ObjectID) ([]*model.Section, error) {
	loaders := For(ctx)
	return loaders.SectionLoader.LoadAll(ctx, sectionIDs)
}
