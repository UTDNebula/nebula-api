package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Models that interact with MongoDB
// I.e., these models have ID reference to other models

type Course struct {
	ID                     bson.ObjectID          `json:"_id" bson:"_id"`
	SubjectPrefix          string                 `json:"subject_prefix" bson:"subject_prefix"`
	CourseNumber           string                 `json:"course_number" bson:"course_number"`
	Title                  string                 `json:"title" bson:"title"`
	Description            string                 `json:"description" bson:"description"`
	EnrollmentReqs         string                 `json:"enrollment_reqs" bson:"enrollment_reqs"`
	School                 string                 `json:"school" bson:"school"`
	CreditHours            string                 `json:"credit_hours" bson:"credit_hours"`
	ClassLevel             string                 `json:"class_level" bson:"class_level"`
	ActivityType           string                 `json:"activity_type" bson:"activity_type"`
	Grading                string                 `json:"grading" bson:"grading"`
	InternalCourseNumber   string                 `json:"internal_course_number" bson:"internal_course_number"`
	Prerequisites          *CollectionRequirement `json:"prerequisites,omitempty" bson:"prerequisites,omitempty"`
	Corequisites           *CollectionRequirement `json:"corequisites,omitempty" bson:"corequisites,omitempty"`
	CoOrPreRequisites      *CollectionRequirement `json:"co_or_pre_requisites,omitempty" bson:"co_or_pre_requisites,omitempty"`
	Sections               []bson.ObjectID        `json:"sections" bson:"sections"`
	LectureContactHours    string                 `json:"lecture_contact_hours" bson:"lecture_contact_hours"`
	LaboratoryContactHours string                 `json:"laboratory_contact_hours" bson:"laboratory_contact_hours"`
	OfferingFrequency      string                 `json:"offering_frequency" bson:"offering_frequency"`
	CatalogYear            string                 `json:"catalog_year" bson:"catalog_year"`
	Attributes             any                    `json:"attributes,omitempty" bson:"attributes,omitempty"`
}

type Professor struct {
	ID          bson.ObjectID   `json:"_id" bson:"_id"`
	FirstName   string          `json:"first_name" bson:"first_name"`
	LastName    string          `json:"last_name" bson:"last_name"`
	Titles      []string        `json:"titles" bson:"titles"`
	Email       string          `json:"email" bson:"email"`
	PhoneNumber string          `json:"phone_number" bson:"phone_number"`
	Office      *Location       `json:"office" bson:"office"`
	ProfileURI  string          `json:"profile_uri" bson:"profile_uri"`
	ImageURI    string          `json:"image_uri" bson:"image_uri"`
	OfficeHours any             `json:"office_hours,omitempty" bson:"office_hours,omitempty"`
	Sections    []bson.ObjectID `json:"sections" bson:"sections"`
}

type Section struct {
	ID                  bson.ObjectID          `json:"_id" bson:"_id"`
	SectionNumber       string                 `json:"section_number" bson:"section_number"`
	CourseReference     bson.ObjectID          `json:"course_reference" bson:"course_reference"`
	SectionCorequisites *CollectionRequirement `json:"section_corequisites,omitempty" bson:"section_corequisites,omitempty"`
	AcademicSession     *AcademicSession       `json:"academic_session" bson:"academic_session"`
	Professors          []bson.ObjectID        `json:"professors" bson:"professors"`
	TeachingAssistants  []*Assistant           `json:"teaching_assistants" bson:"teaching_assistants"`
	InternalClassNumber string                 `json:"internal_class_number" bson:"internal_class_number"`
	InstructionMode     string                 `json:"instruction_mode" bson:"instruction_mode"`
	Meetings            []*Meeting             `json:"meetings" bson:"meetings"`
	CoreFlags           []string               `json:"core_flags" bson:"core_flags"`
	SyllabusURI         string                 `json:"syllabus_uri" bson:"syllabus_uri"`
	GradeDistribution   []int32                `json:"grade_distribution,omitempty" bson:"grade_distribution,omitempty"`
	Attributes          any                    `json:"attributes,omitempty" bson:"attributes,omitempty"`
}

type SectionWithTime struct {
	Section   bson.ObjectID `json:"section" bson:"section"`
	StartTime string        `json:"start_time" bson:"start_time"`
	EndTime   string        `json:"end_time" bson:"end_time"`
}

func (SectionWithTime) IsEventResult() {}
