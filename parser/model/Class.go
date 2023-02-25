package model

import "go.mongodb.org/mongo-driver/bson"

type Class struct {
	Subject           string
	CatalogNumber     string
	Section           string
	GradeDistribution bson.A
}
