package model

import "go.mongodb.org/mongo-driver/bson"

type Exam interface {
	IsExam()
	GetID() string
}

type Outcome interface {
	IsOutcome()
}

type ALEKSExam struct {
	ID        string              `json:"_id" bson:"_id"`
	Placement []*PossibleOutcomes `json:"placement"`
}

func (ALEKSExam) IsExam()            {}
func (this ALEKSExam) GetID() string { return this.ID }

type APExam struct {
	ID     string              `json:"_id" bson:"_id"`
	Name   string              `json:"name"`
	Yields []*PossibleOutcomes `json:"yields"`
}

func (APExam) IsExam()            {}
func (this APExam) GetID() string { return this.ID }

type CLEPExam struct {
	ID     string              `json:"_id" bson:"_id"`
	Name   string              `json:"name"`
	Yields []*PossibleOutcomes `json:"yields"`
}

func (CLEPExam) IsExam()            {}
func (this CLEPExam) GetID() string { return this.ID }

type CSPlacementExam struct {
	ID     string              `json:"_id" bson:"_id"`
	Yields []*PossibleOutcomes `json:"yields"`
}

func (CSPlacementExam) IsExam()            {}
func (this CSPlacementExam) GetID() string { return this.ID }

type IBExam struct {
	ID     string              `json:"_id" bson:"_id"`
	Name   string              `json:"name"`
	Level  string              `json:"level"`
	Yields []*PossibleOutcomes `json:"yields"`
}

func (IBExam) IsExam()            {}
func (this IBExam) GetID() string { return this.ID }

type PossibleOutcomes struct {
	Requirement      bson.Raw   `json:"requirement,omitempty"`
	PossibleOutcomes []bson.Raw `json:"possible_outcomes" bson:"outcome"`
}
