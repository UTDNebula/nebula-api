package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.27

import (
	"context"
	"fmt"

	"github.com/LocatingWizard/nebula_api_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Requirement is the resolver for the requirement field.
func (r *possibleOutcomesResolver) Requirement(ctx context.Context, obj *model.PossibleOutcomes) (model.Requirement, error) {
	bytes, err := bson.Marshal(obj.Requirement)
	if err != nil {
		return nil, err
	}

	requirementType := obj.Requirement.Lookup("type").StringValue()

	switch requirementType {
	case "course":
		var t model.CourseRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "section":
		var t model.SectionRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "exam":
		var t model.ExamRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "major":
		var t model.MajorRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "minor":
		var t model.MinorRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "gpa":
		var t model.GPARequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "consent":
		var t model.ConsentRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "collection":
		var t model.CollectionRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "hours":
		var t model.HoursRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "other":
		var t model.OtherRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "choice":
		var t model.ChoiceRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "limit":
		var t model.LimitRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	case "core":
		var t model.CoreRequirement
		bson.Unmarshal(bytes, &t)
		return t, nil
	default:
		return nil, fmt.Errorf("unkown requirement type: %v", err)
	}
}

// PossibleOutcomes is the resolver for the possible_outcomes field.
func (r *possibleOutcomesResolver) PossibleOutcomes(ctx context.Context, obj *model.PossibleOutcomes) ([][]model.Outcome, error) {
	var out [][]model.Outcome
	for _, v := range obj.PossibleOutcomes {
		var temp []model.Outcome
		vals, err := v.Values()
		if err != nil {
			return nil, err
		}
		for _, w := range vals {
			itemType := w.Type
			if itemType == bsontype.ObjectID {
				id := w.ObjectID()
				course, err := r.Query().CourseByID(ctx, id.Hex())
				if err != nil {
					return nil, err
				}
				temp = append(temp, course)
			} else if itemType == bsontype.EmbeddedDocument {
				var t model.Credit
				bson.Unmarshal(w.Value, &t)
				temp = append(temp, t)
			}
		}
		out = append(out, temp)
	}
	return out, nil
}

// PossibleOutcomes returns PossibleOutcomesResolver implementation.
func (r *Resolver) PossibleOutcomes() PossibleOutcomesResolver { return &possibleOutcomesResolver{r} }

type possibleOutcomesResolver struct{ *Resolver }
