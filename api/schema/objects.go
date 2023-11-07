package schema

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wrapper type for primitive.ObjectID to allow for custom mashalling below
type IdWrapper struct {
	primitive.ObjectID
}

// Custom JSON marshalling for ObjectID to marshal ObjectIDs correctly
func (id IdWrapper) MarshalJSON() (data []byte, err error) {

	type tmp struct {
		Id primitive.ObjectID `json:"$oid"`
	}

	return json.Marshal(tmp{id.ObjectID})
}

type Course struct {
	Id                       IdWrapper              `bson:"_id" json:"_id"`
	Subject_prefix           string                 `bson:"subject_prefix" json:"subject_prefix"`
	Course_number            string                 `bson:"course_number" json:"course_number"`
	Title                    string                 `bson:"title" json:"title"`
	Description              string                 `bson:"description" json:"description"`
	Enrollment_reqs          string                 `bson:"enrollment_reqs" json:"enrollment_reqs"`
	School                   string                 `bson:"school" json:"school"`
	Credit_hours             string                 `bson:"credit_hours" json:"credit_hours"`
	Class_level              string                 `bson:"class_level" json:"class_level"`
	Activity_type            string                 `bson:"activity_type" json:"activity_type"`
	Grading                  string                 `bson:"grading" json:"grading"`
	Internal_course_number   string                 `bson:"internal_course_number" json:"internal_course_number"`
	Prerequisites            *CollectionRequirement `bson:"prerequisites" json:"prerequisites"`
	Corequisites             *CollectionRequirement `bson:"corequisites" json:"corequisites"`
	Co_or_pre_requisites     *CollectionRequirement `bson:"co_or_pre_requisites" json:"co_or_pre_requisites"`
	Sections                 []IdWrapper            `bson:"sections" json:"sections"`
	Lecture_contact_hours    string                 `bson:"lecture_contact_hours" json:"lecture_contact_hours"`
	Laboratory_contact_hours string                 `bson:"laboratory_contact_hours" json:"laboratory_contact_hours"`
	Offering_frequency       string                 `bson:"offering_frequency" json:"offering_frequency"`
	Catalog_year             string                 `bson:"catalog_year" json:"catalog_year"`
	Attributes               interface{}            `bson:"attributes" json:"attributes"`
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
	Start_time   time.Time `bson:"start_time" json:"start_time"`
	End_time     time.Time `bson:"end_time" json:"end_time"`
	Modality     string    `bson:"modality" json:"modality"`
	Location     Location  `bson:"location" json:"location"`
}

type Section struct {
	Id                    IdWrapper              `bson:"_id" json:"_id"`
	Section_number        string                 `bson:"section_number" json:"section_number"`
	Course_reference      IdWrapper              `bson:"course_reference" json:"course_reference"`
	Section_corequisites  *CollectionRequirement `bson:"section_corequisites" json:"section_corequisites"`
	Academic_session      AcademicSession        `bson:"academic_session" json:"academic_session"`
	Professors            []IdWrapper            `bson:"professors" json:"professors"`
	Teaching_assistants   []Assistant            `bson:"teaching_assistants" json:"teaching_assistants"`
	Internal_class_number string                 `bson:"internal_class_number" json:"internal_class_number"`
	Instruction_mode      string                 `bson:"instruction_mode" json:"instruction_mode"`
	Meetings              []Meeting              `bson:"meetings" json:"meetings"`
	Core_flags            []string               `bson:"core_flags" json:"core_flags"`
	Syllabus_uri          string                 `bson:"syllabus_uri" json:"syllabus_uri"`
	Grade_distribution    []int                  `bson:"grade_distribution" json:"grade_distribution"`
	Attributes            interface{}            `bson:"attributes" json:"attributes"`
}

type Professor struct {
	Id           IdWrapper   `bson:"_id" json:"_id"`
	First_name   string      `bson:"first_name" json:"first_name"`
	Last_name    string      `bson:"last_name" json:"last_name"`
	Titles       []string    `bson:"titles" json:"titles"`
	Email        string      `bson:"email" json:"email"`
	Phone_number string      `bson:"phone_number" json:"phone_number"`
	Office       Location    `bson:"office" json:"office"`
	Profile_uri  string      `bson:"profile_uri" json:"profile_uri"`
	Image_uri    string      `bson:"image_uri" json:"image_uri"`
	Office_hours []Meeting   `bson:"office_hours" json:"office_hours"`
	Sections     []IdWrapper `bson:"sections" json:"sections"`
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
