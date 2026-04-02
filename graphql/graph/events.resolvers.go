package graph

import (
	"context"
	"fmt"
	"graphql/graph/model"
)

// Buildings is the resolver for the buildings field.
func (r *multiBuildingEventsResolver) Buildings(ctx context.Context, obj *model.MultiBuildingEvents) ([]*model.SingleBuildingEvents, error) {
	panic(fmt.Errorf("not implemented: Buildings - buildings"))
}

// Events implements [QueryResolver].
func (r *queryResolver) Events(ctx context.Context, filter *model.EventFilter, offset *int32) ([]*model.MultiBuildingEvents, error) {
	panic("unimplemented")
}

// SectionEvents is the resolver for the section_events field.
func (r *roomEventsResolver) SectionEvents(ctx context.Context, obj *model.RoomEvents) ([]*model.SectionWithTime, error) {
	panic(fmt.Errorf("not implemented: SectionEvents - section_events"))
}

// Rooms is the resolver for the rooms field.
func (r *singleBuildingEventsResolver) Rooms(ctx context.Context, obj *model.SingleBuildingEvents) ([]*model.RoomEvents, error) {
	panic(fmt.Errorf("not implemented: Rooms - rooms"))
}