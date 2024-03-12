package schema

import (
	"encoding/json"
	"time"
)

// Wrapper type for primitive.ObjectID to allow for custom mashalling below
type IdWrapper string

// Custom JSON marshalling for ObjectID to marshal ObjectIDs correctly
func (id IdWrapper) MarshalJSON() (data []byte, err error) {

	type tmp struct {
		Id string `json:"$oid"`
	}

	return json.Marshal(tmp{string(id)})
}

// Custom JSON unmarshalling for ObjectID to unmarshal ObjectIDs correctly
func (id *IdWrapper) UnmarshalJSON(data []byte) (e error) {

	type tmp struct {
		Id string `json:"$oid"`
	}

	var t tmp

	e = json.Unmarshal(data, &t)
	*id = IdWrapper(t.Id)
	return e
}

type Course struct {
	Id                       IdWrapper              `bson:"_id" json:"_id" schema:"-"`
	Subject_prefix           string                 `bson:"subject_prefix" json:"subject_prefix" schema:"subject_prefix"`
	Course_number            string                 `bson:"course_number" json:"course_number" schema:"course_number"`
	Title                    string                 `bson:"title" json:"title" schema:"title"`
	Description              string                 `bson:"description" json:"description" schema:"-"`
	Enrollment_reqs          string                 `bson:"enrollment_reqs" json:"enrollment_reqs" schema:"-"`
	School                   string                 `bson:"school" json:"school" schema:"school"`
	Credit_hours             string                 `bson:"credit_hours" json:"credit_hours" schema:"credit_hours"`
	Class_level              string                 `bson:"class_level" json:"class_level" schema:"class_level"`
	Activity_type            string                 `bson:"activity_type" json:"activity_type" schema:"activity_type"`
	Grading                  string                 `bson:"grading" json:"grading" schema:"grading"`
	Internal_course_number   string                 `bson:"internal_course_number" json:"internal_course_number" schema:"internal_course_number"`
	Prerequisites            *CollectionRequirement `bson:"prerequisites" json:"prerequisites" schema:"-"`
	Corequisites             *CollectionRequirement `bson:"corequisites" json:"corequisites" schema:"-"`
	Co_or_pre_requisites     *CollectionRequirement `bson:"co_or_pre_requisites" json:"co_or_pre_requisites" schema:"-"`
	Sections                 []IdWrapper            `bson:"sections" json:"sections" schema:"-"`
	Lecture_contact_hours    string                 `bson:"lecture_contact_hours" json:"lecture_contact_hours" schema:"lecture_contact_hours"`
	Laboratory_contact_hours string                 `bson:"laboratory_contact_hours" json:"laboratory_contact_hours" schema:"laboratory_contact_hours"`
	Offering_frequency       string                 `bson:"offering_frequency" json:"offering_frequency" schema:"offering_frequency"`
	Catalog_year             string                 `bson:"catalog_year" json:"catalog_year" schema:"catalog_year"`
	Attributes               interface{}            `bson:"attributes" json:"attributes" schema:"-"`
}

type AcademicSession struct {
	Name       string    `bson:"name" json:"name"`
	Start_date time.Time `bson:"start_date" json:"start_date"`
	End_date   time.Time `bson:"end_date" json:"end_date"`
}

type Assistant struct {
	First_name string `bson:"first_name" json:"first_name"`
	Last_name  string `bson:"last_name" json:"last_name"`
	Role       string `bson:"role" json:"role"`
	Email      string `bson:"email" json:"email"`
}

type Location struct {
	Building string `bson:"building" json:"building"`
	Room     string `bson:"room" json:"room"`
	Map_uri  string `bson:"map_uri" json:"map_uri"`
}

type Meeting struct {
	Start_date   time.Time `bson:"start_date" json:"start_date"`
	End_date     time.Time `bson:"end_date" json:"end_date"`
	Meeting_days []string  `bson:"meeting_days" json:"meeting_days"`
	Start_time   string    `bson:"start_time" json:"start_time"`
	End_time     string    `bson:"end_time" json:"end_time"`
	Modality     string    `bson:"modality" json:"modality"`
	Location     Location  `bson:"location" json:"location"`
}

type Section struct {
	Id                    IdWrapper              `bson:"_id" json:"_id" schema:"-"`
	Section_number        string                 `bson:"section_number" json:"section_number" schema:"section_number"`
	Course_reference      IdWrapper              `bson:"course_reference" json:"course_reference" schema:"course_reference"`
	Section_corequisites  *CollectionRequirement `bson:"section_corequisites" json:"section_corequisites" schema:"-"`
	Academic_session      AcademicSession        `bson:"academic_session" json:"academic_session" schema:"-"`
	Professors            []IdWrapper            `bson:"professors" json:"professors" schema:"-"`
	Teaching_assistants   []Assistant            `bson:"teaching_assistants" json:"teaching_assistants" schema:"-"`
	Internal_class_number string                 `bson:"internal_class_number" json:"internal_class_number" schema:"internal_class_number"`
	Instruction_mode      string                 `bson:"instruction_mode" json:"instruction_mode" schema:"instruction_mode"`
	Meetings              []Meeting              `bson:"meetings" json:"meetings" schema:"-"`
	Core_flags            []string               `bson:"core_flags" json:"core_flags" schema:"-"`
	Syllabus_uri          string                 `bson:"syllabus_uri" json:"syllabus_uri" schema:"-"`
	Grade_distribution    []int                  `bson:"grade_distribution" json:"grade_distribution" schema:"-"`
	Attributes            interface{}            `bson:"attributes" json:"attributes" schema:"-"`
}

type Professor struct {
	Id           IdWrapper   `bson:"_id" json:"_id" schema:"-"`
	First_name   string      `bson:"first_name" json:"first_name" schema:"first_name"`
	Last_name    string      `bson:"last_name" json:"last_name" schema:"last_name"`
	Titles       []string    `bson:"titles" json:"titles" schema:"titles"`
	Email        string      `bson:"email" json:"email" schema:"email"`
	Phone_number string      `bson:"phone_number" json:"phone_number" schema:"phone_number"`
	Office       Location    `bson:"office" json:"office" schema:"-"`
	Profile_uri  string      `bson:"profile_uri" json:"profile_uri" schema:"-"`
	Image_uri    string      `bson:"image_uri" json:"image_uri" schema:"-"`
	Office_hours []Meeting   `bson:"office_hours" json:"office_hours" schema:"-"`
	Sections     []IdWrapper `bson:"sections" json:"sections" schema:"-"`
}

type Organization struct {
	Id             IdWrapper `bson:"_id" json:"_id"`
	Title          string    `bson:"title" json:"title"`
	Description    string    `bson:"description" json:"description"`
	Categories     []string  `bson:"categories" json:"categories"`
	President_name string    `bson:"president_name" json:"president_name"`
	Emails         []string  `bson:"emails" json:"emails"`
	Picture_data   string    `bson:"picture_data" json:"picture_data"`
}

type Event struct {
	Id                 IdWrapper `bson:"_id" json:"_id"`
	Summary            string    `bson:"summary" json:"summary"`
	Location           string    `bson:"location" json:"location"`
	StartTime          time.Time `bson:"start_time" json:"start_time"`
	EndTime            time.Time `bson:"end_time" json:"end_time"`
	Description        string    `bson:"description" json:"description"`
	EventType          []string  `bson:"event_type" json:"event_type"`
	TargetAudience     []string  `bson:"target_audience" json:"target_audience"`
	Topic              []string  `bson:"topic" json:"topic"`
	EventTags          []string  `bson:"event_tags" json:"event_tags"`
	EventWebsite       string    `bson:"event_website" json:"event_website"`
	Department         []string  `bson:"department" json:"department"`
	ContactName        string    `bson:"contact_name" json:"contact_name"`
	ContactEmail       string    `bson:"contact_email" json:"contact_email"`
	ContactPhoneNumber string    `bson:"contact_phone_number" json:"contact_phone_number"`
}

// 5 Level Likert Item scale for evaluation responses
type EvaluationResponse int

const (
	STRONGLY_DISAGREE EvaluationResponse = 1 << iota
	DISAGREE
	NEUTRAL
	AGREE
	STRONGLY_AGREE
)

type EvaluationSummary struct {
	Median            float32 `bson:"median" json:"median"`
	Mean              float32 `bson:"mean" json:"mean"`
	StandardDeviation float32 `bson:"standard_deviation" json:"standard_deviation"`
	Responses         int     `bson:"responses" json:"responses"`
}

type EvaluationField struct {
	Description string                         `bson:"description" json:"description"`
	Percentages map[EvaluationResponse]float32 `bson:"percentages" json:"percentages"`
	Counts      map[EvaluationResponse]int     `bson:"counts" json:"counts"`
	Summary     EvaluationSummary              `bson:"summary" json:"summary"`
}

type Evaluation struct {
	Id                   IdWrapper         `bson:"_id" json:"_id"`
	CourseExperience     []EvaluationField `bson:"course_experience" json:"course_experience"`
	InstructorExperience []EvaluationField `bson:"instructor_experience" json:"instructor_experience"`
	StudentExperience    []EvaluationField `bson:"student_experience" json:"student_experience"`
}
