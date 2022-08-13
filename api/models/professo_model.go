package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Professor struct {
	ID          primitive.ObjectID   `bson:"_id" json:"_id"`
	FirstName   string               `bson:"first_name" json:"first_name"`
	LastName    string               `bson:"last_name" json:"last_name"`
	Titles      []string             `bson:"titles" json:"titles"`
	Email       string               `bson:"email" json:"email"`
	PhoneNumber string               `bson:"phone_number" json:"phone_number"`
	Office      Location             `bson:"office" json:"office"`
	ProfileURI  string               `bson:"profile_uri" json:"profile_uri"`
	ImageURI    string               `bson:"image_uri" json:"image_uri"`
	OfficeHours []Meeting            `bson:"office_hours" json:"office_hours"`
	Sections    []primitive.ObjectID `bson:"sections" json:"sections"`
	V           int                  `bson:"__v" json:"__v"`
}

type Location struct {
	Building string `bson:"building" json:"building"`
	Room     string `bson:"room" json:"room"`
	MapURI   string `bson:"map_uri" json:"map_uri"`
}

type Meeting struct {
	StartDate   string   `bson:"start_date" json:"start_date"`
	EndDate     string   `bson:"end_date" json:"end_date"`
	MeetingDays []string `bson:"meeting_days" json:"meeting_days"`
	StartTime   string   `bson:"start_time" json:"start_time"`
	EndTime     string   `bson:"end_time" json:"end_time"`
	Modality    string   `bson:"modality" json:"modality"`
	Location    Location `bson:"location" json:"location"`
}
