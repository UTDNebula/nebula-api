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
	var dbSections []*model.DBSection
	if err := cursor.All(ctx, &dbSections); err != nil {
		return nil, []error{err}
	}

	// $in does not return in order of given sectionIDs
	// so we order with map.
	sectionMap := make(map[bson.ObjectID]*model.Section, len(dbSections))
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

// GetSection returns single section by id efficiently
// Note that sectionRef contains only its ObjectID.
// All other fields in sectionRef are empty.
func GetSection(ctx context.Context, sectionRef *model.Section) (*model.Section, error) {
	sectionID, err := bson.ObjectIDFromHex(sectionRef.ID)
	if err != nil {
		return nil, err
	}

	loaders := For(ctx)
	return loaders.SectionLoader.Load(ctx, sectionID)
}

// GetSections returns many sections by ids efficiently
// Note that sectionsRef contains only the sections' ObjectID.
// All other fields of each section in sectionsRef are empty.
func GetSections(ctx context.Context, sectionsRef []*model.Section) ([]*model.Section, error) {
	sectionIDs := make([]bson.ObjectID, len(sectionsRef))
	for i, section := range sectionsRef {
		sectionID, err := bson.ObjectIDFromHex(section.ID)
		if err != nil {
			return nil, err
		}
		sectionIDs[i] = sectionID
	}

	loaders := For(ctx)
	return loaders.SectionLoader.LoadAll(ctx, sectionIDs)
}
