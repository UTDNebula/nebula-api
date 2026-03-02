package graph

import (
	"context"
	"fmt"
	"graphql/graph/model"

)

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context) (*model.MultiBuildingEvents, error) {
	panic(fmt.Errorf("not implemented: Events - events"))
}