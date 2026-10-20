package schema

type Autocomplete struct {
	Course_numbers []CourseNumberAcademicSessions `bson:"course_numbers" json:"course_numbers" schema:"course_numbers"`
	Subject_prefix string                         `bson:"subject_prefix" json:"subject_prefix" schema:"subject_prefix"`
}

type CourseNumberAcademicSessions struct {
	Academic_sessions []AcademicSessionSections `bson:"academic_sessions" json:"academic_sessions" schema:"academic_sessions"`
	Course_number     string                    `bson:"course_number" json:"course_number" schema:"course_number"`
	Title             string                    `bson:"title" json:"title" schema:"title"`
}

type AcademicSessionSections struct {
	Academic_session SimpleAcademicSession     `bson:"academic_session" json:"academic_session" schema:"academic_session"`
	Sections         []SectionNumberProfessors `bson:"sections" json:"sections" schema:"sections"`
}

type SectionNumberProfessors struct {
	Professors     []SimpleProfessor `bson:"professors" json:"professors" schema:"professors"`
	Section_number string            `bson:"section_number" json:"section_number" schema:"section_number"`
	Total_students int               `bson:"total_students" json:"total_students" schema:"total_students"`
}

type SimpleAcademicSession struct {
	Name string `bson:"name" json:"name"`
}

type SimpleProfessor struct {
	First_name string `bson:"first_name" json:"first_name" schema:"first_name"`
	Last_name  string `bson:"last_name" json:"last_name" schema:"last_name"`
}
