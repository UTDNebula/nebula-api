package requirements

type Requirement struct {
	Type string `bson:"type"`
}

type CourseRequirement struct {
	Requirement    `bson:"inline"`
	ClassReference string `bson:"class_reference"`
	MinimumGrade   string `bson:"minimum_grade"`
}

func NewCourseRequirement(classRef, minGrade string) *CourseRequirement {
	return &CourseRequirement{Requirement{"course"}, classRef, minGrade}
}

type SectionRequirement struct {
	Requirement
	SectionReference string `bson:"section_reference"`
}

func NewSectionRequirement(sectionRef string) *SectionRequirement {
	return &SectionRequirement{Requirement{"section"}, sectionRef}
}

type ExamRequirement struct {
	Requirement
	ExamReference string  `bson:"exam_reference"`
	MinimumScore  float64 `bson:"minimum_score"`
}

func NewExamRequirement(examRef string, minScore float64) *ExamRequirement {
	return &ExamRequirement{Requirement{"exam"}, examRef, minScore}
}

type MajorRequirement struct {
	Requirement
	Major string `bson:"major"`
}

func NewMajorRequirement(major string) *MajorRequirement {
	return &MajorRequirement{Requirement{"major"}, major}
}

type MinorRequirement struct {
	Requirement
	Minor string `bson:"minor"`
}

func NewMinorRequirement(minor string) *MinorRequirement {
	return &MinorRequirement{Requirement{"minor"}, minor}
}

type GPARequirement struct {
	Requirement
	Minimum float64 `bson:"minimum"`
	Subset  string  `bson:"subset"`
}

func NewGPARequirement(min float64, subset string) *GPARequirement {
	return &GPARequirement{Requirement{"gpa"}, min, subset}
}

type ConsentRequirement struct {
	Requirement
	Granter string `bson:"granter"`
}

func NewConsentRequirement(granter string) *ConsentRequirement {
	return &ConsentRequirement{Requirement{"consent"}, granter}
}

type OtherRequirement struct {
	Requirement
	Description string `bson:"description"`
	Condition   string `bson:"condition"`
}

func NewOtherRequirement(description, condition string) *OtherRequirement {
	return &OtherRequirement{Requirement{"other"}, description, condition}
}

type CollectionRequirement struct {
	Requirement `bson:"inline"`
	Name        string        `bson:"name"`
	Required    int           `bson:"required"`
	Options     []interface{} `bson:"options"`
}

func NewCollectionRequirement(name string, required int, options []interface{}) *CollectionRequirement {
	return &CollectionRequirement{Requirement{"collection"}, name, required, options}
}

type HoursRequirement struct {
	Requirement
	Required int                  `bson:"required"`
	Options  []*CourseRequirement `bson:"options"`
}

func NewHoursRequirement(required int, options []*CourseRequirement) *HoursRequirement {
	return &HoursRequirement{Requirement{"hours"}, required, options}
}

type ChoiceRequirement struct {
	Requirement
	Choices *CollectionRequirement `bson:"choices"`
}

func NewChoiceRequirement(choices *CollectionRequirement) *ChoiceRequirement {
	return &ChoiceRequirement{Requirement{"choice"}, choices}
}

type LimitRequirement struct {
	Requirement
	MaxHours int `bson:"max_hours"`
}

func NewLimitRequirement(maxHours int) *LimitRequirement {
	return &LimitRequirement{Requirement{"limit"}, maxHours}
}

type CoreRequirement struct {
	Requirement
	CoreFlag string `bson:"core_flag"`
	Hours    int    `bson:"hours"`
}

func NewCoreRequirement(coreFlag string, hours int) *CoreRequirement {
	return &CoreRequirement{Requirement{"core"}, coreFlag, hours}
}

type Degree struct {
	Subtype            string                 `bson:"subtype"`
	School             string                 `bson:"school"`
	Name               string                 `bson:"name"`
	Year               string                 `bson:"year"`
	Abbreviation       string                 `bson:"abbreviation"`
	MinimumCreditHours int                    `bson:"minimum_credit_hours"`
	CatalogUri         string                 `bson:"catalog_uri"`
	Requirements       *CollectionRequirement `bson:"requirements"`
}
