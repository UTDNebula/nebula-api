package loaders

import (
	"context"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type sectionReader struct {
	collection *mongo.Collection
}

// getSections implements a batch function that can retrieve many users by ID for use in a dataloader
func (s *sectionReader) getSections(ctx context.Context, sectionIDs []bson.ObjectID) ([]*model.Section, []error) {
	// Find documents matching sectionIDs
	cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": sectionIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	// Decode
	var sections []*model.Section
	if err := cursor.All(ctx, &sections); err != nil {
		return nil, []error{err}
	}

	// $in does not return in order of given sectionIDs
	// so we order with map.
	sectionMap := make(map[bson.ObjectID]*model.Section, len(sections))
	for _, section := range sections {
		sectionMap[section.ID] = section
	}

	orderedSections := make([]*model.Section, len(sectionIDs))
	for i, ID := range sectionIDs {
		orderedSections[i] = sectionMap[ID]
	}

	return orderedSections, nil
}

// For convenience in resolvers

// GetSection returns single section by id efficiently
func GetSection(ctx context.Context, sectionRef bson.ObjectID) (*model.Section, error) {
	loaders := For(ctx)
	return loaders.SectionLoader.Load(ctx, sectionRef)
}

// GetSections returns many sections by ids efficiently
func GetSections(ctx context.Context, sectionsRef []bson.ObjectID) ([]*model.Section, error) {
	loaders := For(ctx)
	return loaders.SectionLoader.LoadAll(ctx, sectionsRef)
}
