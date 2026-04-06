package graph

import (
	"context"
	"fmt"
	"graphql/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Events(ctx context.Context, date string, building *string, room *string) (model.EventResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var dbEvents model.DBMultiBuildingEvents
	var events *model.MultiBuildingEvents

	err := r.EventCollection.FindOne(timeoutCtx, bson.M{"date": date}).Decode(&dbEvents)
	if err != nil {
		return nil, err
	}

	events = model.TransformMultiBuildingEvents(&dbEvents)
	if building == nil {
		return events, nil
	}

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
	return r.Query().Section(ctx, obj.Section.ID)
}