package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"github.com/UTDNebula/nebula-api/graphql/loaders"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *queryResolver) Events(ctx context.Context, date string, building *string, room *string) (model.EventResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var events *model.MultiBuildingEvents

	err := r.EventCollection.FindOne(timeoutCtx, bson.M{"date": date}).Decode(events)
	if err != nil {
		return nil, err
	}

	if building == nil {
		return events, nil
	}
	
	events = &model.MultiBuildingEvents{}
	for _, b := range events.Buildings {
		if b.Building == *building {
			if room == nil {
				return b, nil
			}

			for _, r := range b.Rooms {
				if r.Room == *room {
					return r, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("date not found")
}

// Section implements [SectionWithTimeResolver].
func (r *sectionWithTimeResolver) Section(ctx context.Context, obj *model.SectionWithTime) (*model.Section, error) {
	return loaders.GetSection(ctx, obj.Section)
}