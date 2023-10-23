package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Requirement interface {
	IsRequirement()
}

type CollectionRequirement struct {
	Name     string     `json:"name" bson:"name"`
	Required int        `json:"required"`
	Options  []bson.Raw `json:"options"`
}

func (CollectionRequirement) IsRequirement() {}

type CourseRequirement struct {
	ClassReference primitive.ObjectID `json:"class_reference" bson:"class_reference"`
	MinimumGrade   string             `json:"minimum_grade" bson:"minimum_grade"`
}

func (CourseRequirement) IsRequirement() {}

type ChoiceRequirement struct {
	Choices *CollectionRequirement `json:"choices"`
}

func (ChoiceRequirement) IsRequirement() {}

type ConsentRequirement struct {
	Granter string `json:"granter"`
}

func (ConsentRequirement) IsRequirement() {}

type CoreRequirement struct {
	CoreFlag string `json:"core_flag" bson:"core_flag"`
	Hours    int    `json:"hours"`
}

func (CoreRequirement) IsRequirement() {}

type ExamRequirement struct {
	ExamReference primitive.ObjectID `json:"exam_reference" bson:"exam_reference"`
	MinimumScore  int                `json:"minimum_score" bson:"minimum_score"`
}

func (ExamRequirement) IsRequirement() {}

type GPARequirement struct {
	Minimum float64 `json:"minimum"`
	Subset  string  `json:"subset"`
}

func (GPARequirement) IsRequirement() {}

type HoursRequirement struct {
	Required int                  `json:"required"`
	Options  []*CourseRequirement `json:"options"`
}

func (HoursRequirement) IsRequirement() {}

type LimitRequirement struct {
	MaxHours int `json:"max_hours" bson:"max_hours"`
}

func (LimitRequirement) IsRequirement() {}

type MajorRequirement struct {
	Major string `json:"major"`
}

func (MajorRequirement) IsRequirement() {}

type MinorRequirement struct {
	Minor string `json:"minor"`
}

func (MinorRequirement) IsRequirement() {}

type OtherRequirement struct {
	Description string `json:"description"`
	Condition   string `json:"condition"`
}

func (OtherRequirement) IsRequirement() {}

type SectionRequirement struct {
	SectionReference *Section `json:"section_reference" bson:"sections_reference"`
}

func (SectionRequirement) IsRequirement() {}
