package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Requirement struct {
	Type string `bson:"type" json:"type"`
}

type CourseRequirement struct {
	Requirement
	ClassReference string `bson:"class_reference" json:"class_reference"`
	MinimumGrade   string `bson:"minimum_grade" json:"minimum_grade"`
}

func NewCourseRequirement(classRef string, minGrade string) *CourseRequirement {
	return &CourseRequirement{Requirement{"course"}, classRef, minGrade}
}

type SectionRequirement struct {
	Requirement
	SectionReference primitive.ObjectID `bson:"section_reference" json:"section_reference"`
}

func NewSectionRequirement(sectionRef primitive.ObjectID) *SectionRequirement {
	return &SectionRequirement{Requirement{"section"}, sectionRef}
}

type ExamRequirement struct {
	Requirement
	ExamReference string  `bson:"exam_reference" json:"exam_reference"`
	MinimumScore  float64 `bson:"minimum_score" json:"minimum_score"`
}

func NewExamRequirement(examRef string, minScore float64) *ExamRequirement {
	return &ExamRequirement{Requirement{"exam"}, examRef, minScore}
}

type MajorRequirement struct {
	Requirement
	Major string `bson:"major" json:"major"`
}

func NewMajorRequirement(major string) *MajorRequirement {
	return &MajorRequirement{Requirement{"major"}, major}
}

type MinorRequirement struct {
	Requirement
	Minor string `bson:"minor" json:"minor"`
}

func NewMinorRequirement(minor string) *MinorRequirement {
	return &MinorRequirement{Requirement{"minor"}, minor}
}

type GPARequirement struct {
	Requirement
	Minimum float64 `bson:"minimum" json:"minimum"`
	Subset  string  `bson:"subset" json:"subset"`
}

func NewGPARequirement(min float64, subset string) *GPARequirement {
	return &GPARequirement{Requirement{"gpa"}, min, subset}
}

type ConsentRequirement struct {
	Requirement
	Granter string `bson:"granter" json:"granter"`
}

func NewConsentRequirement(granter string) *ConsentRequirement {
	return &ConsentRequirement{Requirement{"consent"}, granter}
}

type OtherRequirement struct {
	Requirement
	Description string `bson:"description" json:"description"`
	Condition   string `bson:"condition" json:"condition"`
}

func NewOtherRequirement(description, condition string) *OtherRequirement {
	return &OtherRequirement{Requirement{"other"}, description, condition}
}

type CollectionRequirement struct {
	Requirement
	Name     string        `bson:"name" json:"name"`
	Required int           `bson:"required" json:"required"`
	Options  []interface{} `bson:"options" json:"options"`
}

func NewCollectionRequirement(name string, required int, options []interface{}) *CollectionRequirement {
	return &CollectionRequirement{Requirement{"collection"}, name, required, options}
}

type HoursRequirement struct {
	Requirement
	Required int                  `bson:"required" json:"required"`
	Options  []*CourseRequirement `bson:"options" json:"options"`
}

func NewHoursRequirement(required int, options []*CourseRequirement) *HoursRequirement {
	return &HoursRequirement{Requirement{"hours"}, required, options}
}

type ChoiceRequirement struct {
	Requirement
	Choices *CollectionRequirement `bson:"choices" json:"choices"`
}

func NewChoiceRequirement(choices *CollectionRequirement) *ChoiceRequirement {
	return &ChoiceRequirement{Requirement{"choice"}, choices}
}

type LimitRequirement struct {
	Requirement
	MaxHours int `bson:"max_hours" json:"max_hours"`
}

func NewLimitRequirement(maxHours int) *LimitRequirement {
	return &LimitRequirement{Requirement{"limit"}, maxHours}
}

type CoreRequirement struct {
	Requirement
	CoreFlag string `bson:"core_flag" json:"core_flag"`
	Hours    int    `bson:"hours" json:"hours"`
}

func NewCoreRequirement(coreFlag string, hours int) *CoreRequirement {
	return &CoreRequirement{Requirement{"core"}, coreFlag, hours}
}

type Degree struct {
	Subtype            string                 `bson:"subtype" json:"subtype"`
	School             string                 `bson:"school" json:"school"`
	Name               string                 `bson:"name" json:"name"`
	Year               string                 `bson:"year" json:"year"`
	Abbreviation       string                 `bson:"abbreviation" json:"abbreviation"`
	MinimumCreditHours int                    `bson:"minimum_credit_hours" json:"minimum_credit_hours"`
	CatalogUri         string                 `bson:"catalog_uri" json:"catalog_uri"`
	Requirements       *CollectionRequirement `bson:"requirements" json:"requirements"`
}
