package models

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @TODO: Fix Model - Cannot inline interface{}
type Requirement struct {
	// Used to determine requirement type
	// Type = "course" | "section"    | "exam"  | "major"  | "minor" | "gpa"  | "consent"
	//       | "other" | "collection" | "hours" | "choice" | "limit" | "core"
	Type        string      `json:"type" bson:"type" validate:"required"`
	Requirement interface{} `json:",inline" bson:",inline" validate:"required"`
}
type CourseRequirement struct {
	Class_reference primitive.ObjectID `bson:"class_reference,omitempty" json:"class_reference,omitempty" validate:"required"`
	Minimum_grade   string             `json:"minimum_grade,omitempty" validate:"required"`
}

type SectionRequirement struct {
	Section_reference primitive.ObjectID `bson:"section_reference,omitempty" json:"section_reference,omitempty" validate:"required"`
}

type ExamRequirement struct {
	Exam_reference primitive.ObjectID `bson:"exam_reference,omitempty" json:"exam_reference,omitempty" validate:"required"`
	Minimum_score  int                `json:"minimum_score,omitempty" validate:"required"`
}

type MajorRequirement struct {
	Major string `json:"major"`
}

type MinorRequirement struct {
	Minor string `json:"minor,omitempty" validate:"required"`
}

type GPARequirement struct {
	Minimum float64 `json:"minimum,omitempty" validate:"required"`
	Subset  string  `json:"subset,omitempty" validate:"required"`
}

type ConsentRequirement struct {
	Granter string `json:"granter,omitempty" validate:"required"`
}

type OtherRequirement struct {
	Description string `json:"description,omitempty" validate:"required"`
	Condition   string `json:"condition,omitempty" validate:"required"`
}

type CollectionRequirement struct {
	Name     string        `json:"name,omitempty" validate:"required"`
	Required int           `json:"required,omitempty" validate:"required"`
	Options  []Requirement `json:"options,omitempty" validate:"required"`
	Choices  []Requirement `json:"choices,omitempty" validate:"required"`
}

type HoursRequirement struct {
	Required int           `json:"required,omitempty" validate:"required"`
	Options  []Requirement `json:"options,omitempty" validate:"required"`
	Choices  []Requirement `json:"choices,omitempty" validate:"required"`
}

type ChoiceRequirement struct {
	Choices []Requirement `json:"choices,omitempty" validate:"required"`
}

type LimitRequirement struct {
	Max_hours int `json:"max_hours,omitempty" validate:"required"`
}

type CoreRequirement struct {
	Core_flag string `json:"core_flag,omitempty" validate:"required"`
	Hours     int    `json:"hours,omitempty" validate:"required"`
}

func (req *Requirement) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var rawValue bson.RawValue
	err := bson.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	err = rawValue.Unmarshal(&req)
	if err != nil {
		return err
	}

	var requirement struct {
		Requirement bson.RawValue
	}

	err = rawValue.Unmarshal(&requirement)
	if err != nil {
		return err
	}

	switch req.Type {
	case "course":
		reqType := CourseRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "section":
		reqType := SectionRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "exam":
		reqType := ExamRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "major":
		reqType := MajorRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "minor":
		reqType := MinorRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "gpa":
		reqType := GPARequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "consent":
		reqType := ConsentRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "other":
		reqType := OtherRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "collection":
		reqType := CollectionRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "hours":
		reqType := HoursRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "choice":
		reqType := ChoiceRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "limit":
		reqType := LimitRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	case "core":
		reqType := CoreRequirement{}
		err = requirement.Requirement.Unmarshal(&reqType)
		req.Requirement = reqType
	default:
		return errors.Errorf("Unknown requirement type %s", req.Type)
	}

	return err
}

/*
type Requirement struct {

	// Use to determine requirement type
	// Type = "course" | "section"    | "exam"  | "major"  | "minor" | "gpa"  | "consent"
	//       | "other" | "collection" | "hours" | "choice" | "limit" | "core"
	Type string `json:"type,omitempty" validate:"required"`

	// CourseRequirement
	Class_reference *primitive.ObjectID `bson:"class_reference,omitempty" json:"class_reference,omitempty"`
	Minimum_grade   string              `json:"minimum_grade,omitempty"`

	// SectionRequirement
	Section_reference *primitive.ObjectID `bson:"section_reference,omitempty" json:"section_reference,omitempty"`

	// ExamRequirement
	Exam_reference *primitive.ObjectID `bson:"exam_reference,omitempty" json:"exam_reference,omitempty"`
	Minimum_score  int                 `json:"minimum_score,omitempty"`

	// MajorRequirement
	Major string `json:"major,omitempty"`

	// MinorRequirement
	Minor string `json:"minor,omitempty"`

	// GPARequirement
	Minimum float64 `json:"minimum,omitempty"`
	Subset  string  `json:"subset,omitempty"`

	// ConsentRequirement
	Granter string `json:"granter,omitempty"`

	// OtherRequirement
	Description string `json:"description,omitempty"`
	Condition   string `json:"condition,omitempty"`

	// CollectionRequirement
	Name string `json:"name,omitempty"`

	// @TODO: deal with this self-referentialism, should effectively use the following
	// Options        []Requirement      `json:"options,omitempty"`
	// Choices 	      []Requirement      `json:"choices,omitempty"`

	// CollectionRequirement, HoursRequirement
	Required int           `json:"required,omitempty"`
	Options  []interface{} `json:"options,omitempty"`

	// CollectionRequirement, HoursRequirement, ChoiceRequirement
	Choices []interface{} `json:"choices,omitempty"`

	// LimitRequirement
	Max_hours int `json:"max_hours,omitempty"`

	// CoreRequirement
	Core_flag string `json:"core_flag,omitempty"`
	Hours     int    `json:"hours,omitempty"`
}
*/
