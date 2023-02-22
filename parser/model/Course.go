package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID                     primitive.ObjectID   `bson:"_id"`
	CourseNumber           string               `bson:"course_number"`
	SubjectPrefix          string               `bson:"subject_prefix"`
	Title                  string               `bson:"title"`
	Description            string               `bson:"description"`
	School                 string               `bson:"school"`
	CreditHours            string               `bson:"credit_hours"`
	ClassLevel             string               `bson:"class_level"`
	ActivityType           string               `bson:"activity_type"`
	Grading                string               `bson:"grading"`
	InternalCourseNumber   string               `bson:"internal_course_number"`
	Prerequisites          bson.D               `bson:"prerequisites"`
	Corequisites           bson.D               `bson:"Corequisites"`
	CoOrPreRequisites      bson.D               `bson:"co_or_pre_requisites"`
	Sections               []primitive.ObjectID `bson:"sections"`
	LectureContactHours    string               `bson:"lecture_contact_hours"`
	LaboratoryContactHours string               `bson:"laboratory_contact_hours"`
	OfferingFrequency      string               `bson:"offering_frequency"`
	V                      int                  `bson:"__v"`
}
