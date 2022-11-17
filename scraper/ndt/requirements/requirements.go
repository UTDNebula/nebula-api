package requirements

type Requirement struct {
	Type string `bson:"type" json:"type"`
}

type CourseRequirement struct {
	Requirement    `bson:"inline"`
	ClassReference string `bson:"class_reference" json:"class_reference"`
	MinimumGrade   string `bson:"minimum_grade" json:"minimum_grade"`
}

func NewCourseRequirement(classRef, minGrade string) *CourseRequirement {
	return &CourseRequirement{Requirement{"course"}, classRef, minGrade}
}

type SectionRequirement struct {
	Requirement
	SectionReference string `bson:"section_reference" json:"section_reference"`
}

func NewSectionRequirement(sectionRef string) *SectionRequirement {
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
	Requirement `bson:"inline"`
	Name        string        `bson:"name" json:"name"`
	Required    int           `bson:"required" json:"required"`
	Options     []interface{} `bson:"options" json:"options"`
}

func NewCollectionRequirement(name string, required int, options []interface{}) *CollectionRequirement {
	return &CollectionRequirement{Requirement{"collection"}, name, required, options}
}

type HoursRequirement struct {
	Requirement `bson:"inline"`
	MinHours    int           `bson:"min_hours" json:"min_hours"`
	MaxHours    int           `bson:"max_hours" json:"max_hours"`
	Options     []interface{} `bson:"options" json:"options"`
}

func NewHoursRequirement(minHours int, maxHours int, options []interface{}) *HoursRequirement {
	return &HoursRequirement{Requirement{"hours"}, minHours, maxHours, options}
}

type ChoiceRequirement struct {
	Requirement `bson:"inline"`
	Choices     []interface{} `bson:"choices" json:"choices"`
}

func NewChoiceRequirement(choices []interface{}) *ChoiceRequirement {
	return &ChoiceRequirement{Requirement{"choice"}, choices}
}

type LimitRequirement struct {
	Requirement `bson:"inline"`
	MinHours    int `bson:"min_hours" json:"min_hours"`
	MaxHours    int `bson:"max_hours" json:"max_hours"`
}

func NewLimitRequirement(minHours int, maxHours int) *LimitRequirement {
	return &LimitRequirement{Requirement{"limit"}, minHours, maxHours}
}

type CoreRequirement struct {
	Requirement `bson:"inline"`
	CoreFlag    string `bson:"core_flag" json:"core_flag"`
	Hours       int    `bson:"hours" json:"hours"`
}

func NewCoreRequirement(coreFlag string, hours int) *CoreRequirement {
	return &CoreRequirement{Requirement{"core"}, coreFlag, hours}
}

type ElectiveRequirement struct {
	Requirement  `bson:"inline"`
	ElectiveType string `bson:"elective_type" json:"elective_type"`
	Hours        int    `bson:"hours" json:"hours"`
	Level        string `bson:"level" json:"level"`
}

func NewElectiveRequirement(elective_t string, hours int, level string) *ElectiveRequirement {
	return &ElectiveRequirement{Requirement{"elective"}, elective_t, hours, level}
}

type NotRequirement struct {
	Requirement `bson:"inline"`
	Req         interface{} `bson:"requirement" json:"requirement"`
}

func NewNotRequirement(req interface{}) *NotRequirement {
	return &NotRequirement{Requirement{"not"}, req}
}

type Degree struct {
	Subtype              string                 `bson:"subtype" json:"subtype"`
	School               string                 `bson:"school" json:"school"`
	Name                 string                 `bson:"name" json:"name"`
	Year                 string                 `bson:"year" json:"year"`
	Abbreviation         string                 `bson:"abbreviation" json:"abbreviation"`
	Minimum_Credit_Hours int                    `bson:"minimum_credit_hours" json:"minimum_credit_hours"`
	Catalog_Uri          string                 `bson:"catalog_uri" json:"catalog_uri"`
	Requirements         *CollectionRequirement `bson:"requirements" json:"requirements"`
}
