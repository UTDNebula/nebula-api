package graph

import "go.mongodb.org/mongo-driver/mongo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CourseCollection        *mongo.Collection
	SectionCollection       *mongo.Collection
	ProfCollection          *mongo.Collection
	BuildingCollection      *mongo.Collection
	CometCalendarCollection *mongo.Collection
	EventCollection         *mongo.Collection
}
