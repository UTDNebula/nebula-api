package models

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @TODO: Fix Model - Cannot inline interface{}

// @TODO: Choose implementation
// Good but non-conforming to schema - CANT INLINE INTERFACES
// Includes - @TODO: Schema Rename: Outcome -> PotentialOutcomes

type Exam struct {
	Id   primitive.ObjectID     `bson:"_id" json:"_id" validate:"required"`
	Type string                 `json:"type" validate:"required"`
	Exam map[string]interface{} `bson:",inline" json:",inline"`
}

/*
type APExam struct {
	Name   string              `json:"name" validate:"required"`
	Yields []PotentialOutcomes `json:"yields" validate:"required"`
}

type ALEKSExam struct {
	Placement []PotentialOutcomes `json:"placement" validate:"required"`
}

type CLEPExam struct {
	Name   string              `json:"name" validate:"required"`
	Yields []PotentialOutcomes `json:"yields" validate:"required"`
}

type IBExam struct {
	Name   string              `json:"name" validate:"required"`
	Level  string              `json:"level" validate:"required"`
	Yields []PotentialOutcomes `json:"yields" validate:"required"`
}

type CSPlacementExam struct {
	Yields []PotentialOutcomes `json:"yields" validate:"required"`
}

/*
// Bad but conforming to schema
type Exam struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id" validate:"required"`
	Type   string             `json:"type" validate:"required"`
	Name   string             `json:"name,omitempty"`
	Level  string             `json:"level,omitempty"`
	Yields []Yield            `json:"yields" validate:"required"`
}
*/

// @TODO: Schema Rename: Outcome -> PotentialOutcomes
type PotentialOutcomes struct {
	Requirement Requirement `bson:"requirement" json:"requirement" validate:"required"`

	// @TODO: Handle Outcome types
	//      - Outcome : [](primitive.ObjectID | {string, int})

	// messy and broken:
	// Outcomes [][]Outcome `json:"outcome,omitempty" validate:"required"`

	Outcomes [][]interface{} `json:"outcomes" validate:"required"`
}

/*
// @TODO: Handle Outcome types

// messy and broken
type Outcome struct {
	CourseId     *primitive.ObjecID  `bson:",omitempty" json:",omitempty"`
	Category     string              `json:"category,omitempty"`
	Credit_hours int                 `json:"credit_hours,omitempty"`
}

// ---------------------------

// Potential solution ::
*/
type Outcome struct {
	Type string `bson:"type" json:"type" validate:"required"`
	// @TODO: CANT INLINE AN INTERFACE - How to resolve?
	Outcome interface{} `bson:",inline" json:",inline" validate:"required"`
}

type Credit struct {
	Category     string `bson:"category json:"category" validate:"required"`
	Credit_hours int    `bson:"credit_hours json:"credit_hours" validate:"required"`
}

// Custom Exam Unmarshalling function
func (ex *Exam) UnmarshalBSON(t bsontype.Type, data []byte) error {
	var rawValue bson.RawValue
	err := bson.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	err = rawValue.Unmarshal(&ex)
	if err != nil {
		return err
	}

	var exam struct {
		Exam bson.RawValue
	}

	err = rawValue.Unmarshal(&exam)
	if err != nil {
		return err
	}

	switch ex.Type {
	case "AP":
		var exType map[string]interface{}
		err = exam.Exam.Unmarshal(&exType)
		ex.Exam = exType
	case "ALEKS":
		var exType map[string]interface{}
		err = exam.Exam.Unmarshal(&exType)
		ex.Exam = exType
	case "CLEP":
		var exType map[string]interface{}
		err = exam.Exam.Unmarshal(&exType)
		ex.Exam = exType
	case "IB":
		var exType map[string]interface{}
		err = exam.Exam.Unmarshal(&exType)
		ex.Exam = exType
	case "CS placement":
		var exType map[string]interface{}
		err = exam.Exam.Unmarshal(&exType)
		ex.Exam = exType
	default:
		return errors.Errorf("Unknown exam type %s", ex.Type)
	}

	return err
}

func (out *Outcome) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var rawValue bson.RawValue
	err := bson.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	err = rawValue.Unmarshal(&out)
	if err != nil {
		return err
	}

	var outcome struct {
		Outcome bson.RawValue
	}

	err = rawValue.Unmarshal(&outcome)
	if err != nil {
		return err
	}

	switch out.Type {
	case "course":
		outType := primitive.ObjectID{}
		err = outcome.Outcome.Unmarshal(&outType)
		out.Outcome = outType
	case "credit":
		outType := Credit{}
		err = outcome.Outcome.Unmarshal(&outType)
		out.Outcome = outType
	default:
		return errors.Errorf("Unknown outcome type %s", out.Type)
	}

	return err
}

//*/

// // @TODO: Choose implementation
// // Good but non-conforming to schema - CANT INLINE INTERFACES
// // Includes - @TODO: Schema Rename: Outcome -> PotentialOutcomes
// type Exam struct {
// 	Id   primitive.ObjectID `bson:"_id" json:"_id" validate:"required"`
// 	Type string             `json:"type" validate:"required"`
// 	Exam interface{}        `bson:",inline" json:",inline" validate:"required"`
// }

// type APExam struct {
// 	Name   string              `json:"name" validate:"required"`
// 	Yields []PotentialOutcomes `json:"yields" validate:"required"`
// }

// type ALEKSExam struct {
// 	Placement []PotentialOutcomes `json:"placement" validate:"required"`
// }

// type CLEPExam struct {
// 	Name   string              `json:"name" validate:"required"`
// 	Yields []PotentialOutcomes `json:"yields" validate:"required"`
// }

// type IBExam struct {
// 	Name   string              `json:"name" validate:"required"`
// 	Level  string              `json:"level" validate:"required"`
// 	Yields []PotentialOutcomes `json:"yields" validate:"required"`
// }

// type CSPlacementExam struct {
// 	Yields []PotentialOutcomes `json:"yields" validate:"required"`
// }

// /*
// // Bad but conforming to schema
// type Exam struct {
// 	Id     primitive.ObjectID `bson:"_id" json:"_id" validate:"required"`
// 	Type   string             `json:"type" validate:"required"`
// 	Name   string             `json:"name,omitempty"`
// 	Level  string             `json:"level,omitempty"`
// 	Yields []Yield            `json:"yields" validate:"required"`
// }
// */

// // @TODO: Schema Rename: Outcome -> PotentialOutcomes
// type PotentialOutcomes struct {
// 	Requirement Requirement `bson:"requirement" json:"requirement" validate:"required"`

// 	// @TODO: Handle Outcome types
// 	//      - Outcome : [](primitive.ObjectID | {string, int})

// 	// messy and broken:
// 	// Outcomes [][]Outcome `json:"outcome,omitempty" validate:"required"`

// 	Outcomes [][]interface{} `json:"outcomes" validate:"required"`
// }

// /*
// // @TODO: Handle Outcome types

// // messy and broken
// type Outcome struct {
// 	CourseId     *primitive.ObjecID  `bson:",omitempty" json:",omitempty"`
// 	Category     string              `json:"category,omitempty"`
// 	Credit_hours int                 `json:"credit_hours,omitempty"`
// }

// // ---------------------------

// // Potential solution ::
// */
// type Outcome struct {
// 	Type string `bson:"type" json:"type" validate:"required"`
// 	// @TODO: CANT INLINE AN INTERFACE - How to resolve?
// 	Outcome interface{} `bson:",inline" json:",inline" validate:"required"`
// }

// type Credit struct {
// 	Category     string `bson:"category json:"category" validate:"required"`
// 	Credit_hours int    `bson:"credit_hours json:"credit_hours" validate:"required"`
// }

// // Custom Exam Unmarshalling function
// func (ex *Exam) UnmarshalBSON(t bsontype.Type, data []byte) error {
// 	var rawValue bson.RawValue
// 	err := bson.Unmarshal(data, &rawValue)
// 	if err != nil {
// 		return err
// 	}

// 	err = rawValue.Unmarshal(&ex)
// 	if err != nil {
// 		return err
// 	}

// 	var exam struct {
// 		Exam bson.RawValue
// 	}

// 	err = rawValue.Unmarshal(&exam)
// 	if err != nil {
// 		return err
// 	}

// 	switch ex.Type {
// 	case "AP":
// 		exType := APExam{}
// 		err = exam.Exam.Unmarshal(&exType)
// 		ex.Exam = exType
// 	case "ALEKS":
// 		exType := ALEKSExam{}
// 		err = exam.Exam.Unmarshal(&exType)
// 		ex.Exam = exType
// 	case "CLEP":
// 		exType := CLEPExam{}
// 		err = exam.Exam.Unmarshal(&exType)
// 		ex.Exam = exType
// 	case "IB":
// 		exType := IBExam{}
// 		err = exam.Exam.Unmarshal(&exType)
// 		ex.Exam = exType
// 	case "CS placement":
// 		exType := CSPlacementExam{}
// 		err = exam.Exam.Unmarshal(&exType)
// 		ex.Exam = exType
// 	default:
// 		return errors.Errorf("Unknown exam type %s", ex.Type)
// 	}

// 	return err
// }

// func (out *Outcome) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
// 	var rawValue bson.RawValue
// 	err := bson.Unmarshal(data, &rawValue)
// 	if err != nil {
// 		return err
// 	}

// 	err = rawValue.Unmarshal(&out)
// 	if err != nil {
// 		return err
// 	}

// 	var outcome struct {
// 		Outcome bson.RawValue
// 	}

// 	err = rawValue.Unmarshal(&outcome)
// 	if err != nil {
// 		return err
// 	}

// 	switch out.Type {
// 	case "course":
// 		outType := primitive.ObjectID{}
// 		err = outcome.Outcome.Unmarshal(&outType)
// 		out.Outcome = outType
// 	case "credit":
// 		outType := Credit{}
// 		err = outcome.Outcome.Unmarshal(&outType)
// 		out.Outcome = outType
// 	default:
// 		return errors.Errorf("Unknown outcome type %s", out.Type)
// 	}

// 	return err
// }

// //*/
