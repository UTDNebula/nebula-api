package controllers

import (
	"github.com/UTDNebula/nebula-api/rest/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Legacy evaluation handlers acquire the course collection on demand.
func courseCollection() *mongo.Collection {
	return configs.GetCollection("courses")
}

func sectionCollection() *mongo.Collection   { return configs.GetCollection("sections") }
func professorCollection() *mongo.Collection { return configs.GetCollection("professors") }
