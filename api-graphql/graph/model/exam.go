package model

type Exam interface {
	IsExam()
}

type Outcome interface {
	IsOutcome()
}

type PossibleOutcomes struct {
	Requirement      Requirement `json:"requirement"`
	PossibleOutcomes [][]Outcome `json:"possible_outcomes"`
}

type ALEKSExam struct {
	ID        string             `json:"_id"`
	Placement []PossibleOutcomes `json:"placement"`
}

func (ALEKSExam) IsExam() {}

type APExam struct {
	ID     string             `json:"_id"`
	Name   string             `json:"name"`
	Yields []PossibleOutcomes `json:"yields"`
}

func (APExam) IsExam() {}

type CLEPExam struct {
	ID     string             `json:"_id"`
	Name   string             `json:"name"`
	Yields []PossibleOutcomes `json:"yields"`
}

func (CLEPExam) IsExam() {}

type CSPlacementExam struct {
	ID     string             `json:"_id"`
	Yields []PossibleOutcomes `json:"yields"`
}

func (CSPlacementExam) IsExam() {}

type IBExam struct {
	ID     string             `json:"_id"`
	Name   string             `json:"name"`
	Level  string             `json:"level"`
	Yields []PossibleOutcomes `json:"yields"`
}

func (IBExam) IsExam() {}
