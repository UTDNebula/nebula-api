package schema

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Requirement struct {
	Type string `bson:"type" json:"type"`
}

type CourseRequirement struct {
	Requirement    `bson:",inline" json:",inline"`
	ClassReference string `bson:"class_reference" json:"class_reference"`
	MinimumGrade   string `bson:"minimum_grade" json:"minimum_grade"`
}

func NewCourseRequirement(classRef string, minGrade string) *CourseRequirement {
	return &CourseRequirement{Requirement{"course"}, classRef, minGrade}
}

type SectionRequirement struct {
	Requirement      `bson:",inline" json:",inline"`
	SectionReference primitive.ObjectID `bson:"section_reference" json:"section_reference"`
}

func NewSectionRequirement(sectionRef primitive.ObjectID) *SectionRequirement {
	return &SectionRequirement{Requirement{"section"}, sectionRef}
}

type ExamRequirement struct {
	Requirement   `bson:",inline" json:",inline"`
	ExamReference string  `bson:"exam_reference" json:"exam_reference"`
	MinimumScore  float64 `bson:"minimum_score" json:"minimum_score"`
}

func NewExamRequirement(examRef string, minScore float64) *ExamRequirement {
	return &ExamRequirement{Requirement{"exam"}, examRef, minScore}
}

type MajorRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Major       string `bson:"major" json:"major"`
}

func NewMajorRequirement(major string) *MajorRequirement {
	return &MajorRequirement{Requirement{"major"}, major}
}

type MinorRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Minor       string `bson:"minor" json:"minor"`
}

func NewMinorRequirement(minor string) *MinorRequirement {
	return &MinorRequirement{Requirement{"minor"}, minor}
}

type GPARequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Minimum     float64 `bson:"minimum" json:"minimum"`
	Subset      string  `bson:"subset" json:"subset"`
}

func NewGPARequirement(min float64, subset string) *GPARequirement {
	return &GPARequirement{Requirement{"gpa"}, min, subset}
}

type ConsentRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Granter     string `bson:"granter" json:"granter"`
}

func NewConsentRequirement(granter string) *ConsentRequirement {
	return &ConsentRequirement{Requirement{"consent"}, granter}
}

type OtherRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Description string `bson:"description" json:"description"`
	Condition   string `bson:"condition" json:"condition"`
}

func NewOtherRequirement(description, condition string) *OtherRequirement {
	return &OtherRequirement{Requirement{"other"}, description, condition}
}

type CollectionRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Name        string        `bson:"name" json:"name"`
	Required    int           `bson:"required" json:"required"`
	Options     []interface{} `bson:"options" json:"options"`
}

type CollectionRequirementIntermediate struct {
	Name     string     `bson:"name"`
	Required int        `bson:"required"`
	Options  []bson.Raw `bson:"options"`
}

func NewCollectionRequirement(name string, required int, options []interface{}) *CollectionRequirement {
	return &CollectionRequirement{Requirement{"collection"}, name, required, options}
}

func (cr *CollectionRequirement) UnmarshalBSON(data []byte) error {
	var dummyCollection CollectionRequirementIntermediate
	err := bson.Unmarshal(data, &dummyCollection)
	if err != nil {
		return err
	}

	var out []interface{}
	for _, v := range dummyCollection.Options {

		bytes, err := bson.Marshal(v)
		if err != nil {
			return err
		}

		optionType := v.Lookup("type").StringValue()

		switch optionType {
		case "course":
			var t CourseRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "section":
			var t SectionRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "exam":
			var t ExamRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "major":
			var t MajorRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "minor":
			var t MinorRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "gpa":
			var t GPARequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "consent":
			var t ConsentRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "collection":
			var t CollectionRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "hours":
			var t HoursRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "other":
			var t OtherRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "choice":
			var t ChoiceRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "limit":
			var t LimitRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		case "core":
			var t CoreRequirement
			bson.Unmarshal(bytes, &t)
			out = append(out, t)
		default:
			return fmt.Errorf("unknown option type: %v", err)
		}
	}
	cr.Name = dummyCollection.Name
	cr.Required = dummyCollection.Required
	cr.Options = out
	cr.Type = "collection"
	return nil
}

type HoursRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Required    int                  `bson:"required" json:"required"`
	Options     []*CourseRequirement `bson:"options" json:"options"`
}

func NewHoursRequirement(required int, options []*CourseRequirement) *HoursRequirement {
	return &HoursRequirement{Requirement{"hours"}, required, options}
}

type ChoiceRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	Choices     *CollectionRequirement `bson:"choices" json:"choices"`
}

func NewChoiceRequirement(choices *CollectionRequirement) *ChoiceRequirement {
	return &ChoiceRequirement{Requirement{"choice"}, choices}
}

type LimitRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	MaxHours    int `bson:"max_hours" json:"max_hours"`
}

func NewLimitRequirement(maxHours int) *LimitRequirement {
	return &LimitRequirement{Requirement{"limit"}, maxHours}
}

type CoreRequirement struct {
	Requirement `bson:",inline" json:",inline"`
	CoreFlag    string `bson:"core_flag" json:"core_flag"`
	Hours       int    `bson:"hours" json:"hours"`
}

func NewCoreRequirement(coreFlag string, hours int) *CoreRequirement {
	return &CoreRequirement{Requirement{"core"}, coreFlag, hours}
}

type Degree struct {
	Subtype            string                 `bson:"subtype" json:"subtype" schema:"subtype"`
	School             string                 `bson:"school" json:"school" schema:"school"`
	Name               string                 `bson:"name" json:"name" schema:"name"`
	Year               string                 `bson:"year" json:"year" schema:"year"`
	Abbreviation       string                 `bson:"abbreviation" json:"abbreviation" schema:"abbreviation"`
	MinimumCreditHours int                    `bson:"minimum_credit_hours" json:"minimum_credit_hours" schema:"minimum_credit_hours"`
	CatalogUri         string                 `bson:"catalog_uri" json:"catalog_uri" schema:"-"`
	Requirements       *CollectionRequirement `bson:"requirements" json:"requirements" schema:"-"`
}
