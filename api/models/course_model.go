package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// @TODO: Fix Model
type Course struct {
	Id     primitive.ObjectID     `bson:"_id" json:"_id" validate:"required"`
	Name   string                 `json:"name,omitempty" validate:"required"`
	Course map[string]interface{} `bson:",inline" json:",inline" validate:"required"`
}

/*
type Course struct {
	Id                       primitive.ObjectID   `json:"_id,omitempty" validate:"required"`
	Name                     string               `json:"name,omitempty" validate:"required"`
	Location                 string               `json:"location,omitempty" validate:"required"`
	Course_number            string               `json:"course_number,omitempty" validate:"required"`
	Subject_prefix           string               `json:"subject_prefix,omitempty" validate:"required"`
	Title                    string               `json:"title,omitempty" validate:"required"`
	Description              string               `json:"description,omitempty" validate:"required"`
	School                   string               `json:"school,omitempty" validate:"required"`
	Credit_hours             string               `json:"credit_hours,omitempty" validate:"required"`
	Class_level              string               `json:"class_level,omitempty" validate:"required"`
	Activity_type            string               `json:"activity_type,omitempty" validate:"required"`
	Grading                  string               `json:"grading,omitempty" validate:"required"`
	Internal_course_number   string               `json:"internal_course_number,omitempty" validate:"required"`
	Prerequisites            interface{}          `json:"prerequisites,omitempty" validate:"required"`
	Corequisites             interface{}          `json:"corequisites,omitempty" validate:"required"`
	Co_or_pre_requisites     interface{}          `json:"co_or_pre_requisites,omitempty" validate:"required"`
	Sections                 []primitive.ObjectID `json:"sections,omitempty" validate:"required"`
	Lecture_contact_hours    string               `json:"lecture_contact_hours,omitempty"`
	Laboratory_contact_hours string               `json:"laboratory_contact_hours,omitempty"`
	Offering_frequency       string               `json:"offering_frequency,omitempty"`
	Attributes               interface{}          `json:"attributes,omitempty" validate:"required"`
}
*/
