package schema

type Requirement struct {
	Type string `json:"type"`
}

type CourseRequirement struct {
	Requirement
	ClassReference string `json:"class_reference"`
	MinimumGrade   string `json:"minimum_grade"`
}

func NewCourseRequirement(classRef string, minGrade string) *CourseRequirement {
	return &CourseRequirement{Requirement{"course"}, classRef, minGrade}
}

type SectionRequirement struct {
	Requirement
	SectionReference IdWrapper `json:"section_reference"`
}

func NewSectionRequirement(sectionRef IdWrapper) *SectionRequirement {
	return &SectionRequirement{Requirement{"section"}, sectionRef}
}

type ExamRequirement struct {
	Requirement
	ExamReference string  `json:"exam_reference"`
	MinimumScore  float64 `json:"minimum_score"`
}

func NewExamRequirement(examRef string, minScore float64) *ExamRequirement {
	return &ExamRequirement{Requirement{"exam"}, examRef, minScore}
}

type MajorRequirement struct {
	Requirement
	Major string `json:"major"`
}

func NewMajorRequirement(major string) *MajorRequirement {
	return &MajorRequirement{Requirement{"major"}, major}
}

type MinorRequirement struct {
	Requirement
	Minor string `json:"minor"`
}

func NewMinorRequirement(minor string) *MinorRequirement {
	return &MinorRequirement{Requirement{"minor"}, minor}
}

type GPARequirement struct {
	Requirement
	Minimum float64 `json:"minimum"`
	Subset  string  `json:"subset"`
}

func NewGPARequirement(min float64, subset string) *GPARequirement {
	return &GPARequirement{Requirement{"gpa"}, min, subset}
}

type ConsentRequirement struct {
	Requirement
	Granter string `json:"granter"`
}

func NewConsentRequirement(granter string) *ConsentRequirement {
	return &ConsentRequirement{Requirement{"consent"}, granter}
}

type OtherRequirement struct {
	Requirement
	Description string `json:"description"`
	Condition   string `json:"condition"`
}

func NewOtherRequirement(description, condition string) *OtherRequirement {
	return &OtherRequirement{Requirement{"other"}, description, condition}
}

type CollectionRequirement struct {
	Requirement
	Name     string        `json:"name"`
	Required int           `json:"required"`
	Options  []interface{} `json:"options"`
}

func NewCollectionRequirement(name string, required int, options []interface{}) *CollectionRequirement {
	return &CollectionRequirement{Requirement{"collection"}, name, required, options}
}

type HoursRequirement struct {
	Requirement
	Required int                  `json:"required"`
	Options  []*CourseRequirement `json:"options"`
}

func NewHoursRequirement(required int, options []*CourseRequirement) *HoursRequirement {
	return &HoursRequirement{Requirement{"hours"}, required, options}
}

type ChoiceRequirement struct {
	Requirement
	Choices *CollectionRequirement `json:"choices"`
}

func NewChoiceRequirement(choices *CollectionRequirement) *ChoiceRequirement {
	return &ChoiceRequirement{Requirement{"choice"}, choices}
}

type LimitRequirement struct {
	Requirement
	MaxHours int `json:"max_hours"`
}

func NewLimitRequirement(maxHours int) *LimitRequirement {
	return &LimitRequirement{Requirement{"limit"}, maxHours}
}

type CoreRequirement struct {
	Requirement
	CoreFlag string `json:"core_flag"`
	Hours    int    `json:"hours"`
}

func NewCoreRequirement(coreFlag string, hours int) *CoreRequirement {
	return &CoreRequirement{Requirement{"core"}, coreFlag, hours}
}

type Degree struct {
	Subtype            string                 `json:"subtype"`
	School             string                 `json:"school"`
	Name               string                 `json:"name"`
	Year               string                 `json:"year"`
	Abbreviation       string                 `json:"abbreviation"`
	MinimumCreditHours int                    `json:"minimum_credit_hours"`
	CatalogUri         string                 `json:"catalog_uri"`
	Requirements       *CollectionRequirement `json:"requirements"`
}
