package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID                     string                 `json:"_id" bson:"_id"`
	CourseNumber           string                 `json:"course_number" bson:"course_number"`
	SubjectPrefix          string                 `json:"subject_prefix" bson:"subject_prefix"`
	Title                  string                 `json:"title"`
	Description            string                 `json:"description"`
	EnrollmentReqs         string                 `json:"enrollment_reqs" bson:"enrollment_reqs"`
	School                 string                 `json:"school"`
	CreditHours            string                 `json:"credit_hours" bson:"credit_hours"`
	ClassLevel             string                 `json:"class_level" bson:"class_level"`
	ActivityType           string                 `json:"activity_type" bson:"activity_type"`
	Grading                string                 `json:"grading"`
	InternalCourseNumber   string                 `json:"internal_course_number" bson:"internal_course_number"`
	Prerequisites          *CollectionRequirement `json:"prerequisites"`
	Corequisites           *CollectionRequirement `json:"corequisites"`
	CoOrPreRequisites      *CollectionRequirement `json:"co_or_pre_requisites" bson:"co_or_pre_requisites"`
	Sections               []primitive.ObjectID   `json:"sections"`
	LectureContactHours    string                 `json:"lecture_contact_hours" bson:"lecture_contact_hours"`
	LaboratoryContactHours string                 `json:"laboratory_contact_hours" bson:"laboratory_contact_hours"`
	OfferingFrequency      string                 `json:"offering_frequency" bson:"offering_frequency"`
	CatalogYear            string                 `json:"catalog_year" bson:"catalog_year"`
	Attributes             interface{}            `json:"attributes"`
}

func (Course) IsOutcome() {}
