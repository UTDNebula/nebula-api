package model

type Attributes struct {
	RawAttributes []string `json:"raw_attributes"`
}

type Section struct {
	ID                  string                 `json:"_id"`
	SectionNumber       string                 `json:"section_number"`
	CourseReference     *Course                `json:"course_reference"`
	SectionCorequisites *CollectionRequirement `json:"section_corequisites"`
	AcademicSession     *AcademicSession       `json:"academic_session"`
	Professors          []*Professor           `json:"professors"`
	TeachingAssistants  []*Assistant           `json:"teaching_assistants"`
	InternalClassNumber string                 `json:"internal_class_number"`
	InstructionMode     string                 `json:"instruction_mode"`
	Meetings            []*Meeting             `json:"meetings"`
	CoreFlags           []string               `json:"core_flags"`
	SyllabusURI         string                 `json:"syllabus_uri"`
	GradeDistribution   []int                  `json:"grade_distribution"`
	Attributes          *Attributes            `json:"attributes"`
}
