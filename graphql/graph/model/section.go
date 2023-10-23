package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Section struct {
	ID                  string                 `json:"_id" bson:"_id"`
	SectionNumber       string                 `json:"section_number" bson:"section_number"`
	CourseReference     primitive.ObjectID     `json:"course_reference" bson:"course_reference"`
	SectionCorequisites *CollectionRequirement `json:"section_corequisites" bson:"sections_corequisites"`
	AcademicSession     *AcademicSession       `json:"academic_session" bson:"academic_session"`
	Professors          []primitive.ObjectID   `json:"professors"`
	TeachingAssistants  []*Assistant           `json:"teaching_assistants" bson:"teaching_assistants"`
	InternalClassNumber string                 `json:"internal_class_number" bson:"internal_class_number"`
	InstructionMode     string                 `json:"instruction_mode" bson:"instruction_mode"`
	Meetings            []*Meeting             `json:"meetings"`
	CoreFlags           []string               `json:"core_flags" bson:"core_flags"`
	SyllabusURI         string                 `json:"syllabus_uri" bson:"syllabus_uri"`
	GradeDistribution   []int                  `json:"grade_distribution" bson:"grade_distribution"`
	Attributes          *Attributes            `json:"attributes"`
}

type Attributes struct {
	RawAttributes []string `json:"raw_attributes" bson:"raw_attributes"`
}
