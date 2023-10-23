package graph

//go:generate go run github.com/99designs/gqlgen generate

import "go.mongodb.org/mongo-driver/mongo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *mongo.Database
}
