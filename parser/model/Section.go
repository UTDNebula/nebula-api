package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Section struct {
	ID                  primitive.ObjectID `bson:"_id"`
	SectionNumber       string             `bson:"section_number"`
	CourseReference     primitive.ObjectID `bson:"course_reference"`
	SectionCorequisites bson.D             `bson:"section_corequisites"`
	Session             AcademicSession    `bson:"academic_session"`
	Professors          bson.A             `bson:"professors"`
	TeachingAssistants  bson.A             `bson:"teaching_assistants"`
	InternalClassNumber string             `bson:"internal_class_number"`
	InstructionMode     string             `bson:"instruction_mode"`
	Meetings            bson.A             `bson:"meetings"`
	CoreFlags           bson.A             `bson:"core_flags"`
	SyllabusUri         string             `bson:"syllabus_uri"`
	GradeDistribution   bson.A             `bson:"grade_distribution"`
	Attributes          bson.D             `bson:"attributes"`
	V                   int                `bson:"__v"`
}
