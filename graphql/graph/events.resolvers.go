package graph

import (
	"context"
	"graphql/configs"
	"graphql/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Events is the resolver for the events field.
