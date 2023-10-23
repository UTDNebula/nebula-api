package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Professor struct {
	ID          string               `json:"_id" bson:"_id"`
	FirstName   string               `json:"first_name" bson:"first_name"`
	LastName    string               `json:"last_name" bson:"last_name"`
	Titles      []string             `json:"titles"`
	Email       string               `json:"email"`
	PhoneNumber *string              `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Office      *Location            `json:"office,omitempty" bson:",omitempty"`
	ProfileURI  *string              `json:"profile_uri,omitempty" bson:"profile_uri,omitempty"`
	ImageURI    *string              `json:"image_uri,omitempty" bson:"image_uri,omitempty"`
	OfficeHours []*Meeting           `json:"office_hours" bson:"office_hours"`
	Sections    []primitive.ObjectID `json:"sections"`
}
