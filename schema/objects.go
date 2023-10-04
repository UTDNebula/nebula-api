package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// This wrapper is necessary to properly represent the ObjectID structure when serializing to JSON instead of BSON
type IdWrapper struct {
	Id primitive.ObjectID `json:"$oid"`
}

type Course struct {
	Id                       IdWrapper              `json:"_id"`
	Subject_prefix           string                 `json:"subject_prefix"`
	Course_number            string                 `json:"course_number"`
	Title                    string                 `json:"title"`
	Description              string                 `json:"description"`
	Enrollment_reqs          string                 `json:"enrollment_reqs"`
	School                   string                 `json:"school"`
	Credit_hours             string                 `json:"credit_hours"`
	Class_level              string                 `json:"class_level"`
	Activity_type            string                 `json:"activity_type"`
	Grading                  string                 `json:"grading"`
	Internal_course_number   string                 `json:"internal_course_number"`
	Prerequisites            *CollectionRequirement `json:"prerequisites"`
	Corequisites             *CollectionRequirement `json:"corequisites"`
	Co_or_pre_requisites     *CollectionRequirement `json:"co_or_pre_requisites"`
	Sections                 []IdWrapper            `json:"sections"`
	Lecture_contact_hours    string                 `json:"lecture_contact_hours"`
	Laboratory_contact_hours string                 `json:"laboratory_contact_hours"`
	Offering_frequency       string                 `json:"offering_frequency"`
	Catalog_year             string                 `json:"catalog_year"`
	Attributes               interface{}            `json:"attributes"`
}

type AcademicSession struct {
	Name       string    `json:"name"`
	Start_date time.Time `json:"start_date"`
	End_date   time.Time `json:"end_date"`
}

type Assistant struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Role       string `json:"role"`
	Email      string `json:"email"`
}

type Location struct {
	Building string `json:"building"`
	Room     string `json:"room"`
	Map_uri  string `json:"map_uri"`
}

type Meeting struct {
	Start_date   time.Time `json:"start_date"`
	End_date     time.Time `json:"end_date"`
	Meeting_days []string  `json:"meeting_days"`
	Start_time   time.Time `json:"start_time"`
	End_time     time.Time `json:"end_time"`
	Modality     string    `json:"modality"`
	Location     Location  `json:"location"`
}

type Section struct {
	Id                    IdWrapper              `json:"_id"`
	Section_number        string                 `json:"section_number"`
	Course_reference      IdWrapper              `json:"course_reference"`
	Section_corequisites  *CollectionRequirement `json:"section_corequisites"`
	Academic_session      AcademicSession        `json:"academic_session"`
	Professors            []IdWrapper            `json:"professors"`
	Teaching_assistants   []Assistant            `json:"teaching_assistants"`
	Internal_class_number string                 `json:"internal_class_number"`
	Instruction_mode      string                 `json:"instruction_mode"`
	Meetings              []Meeting              `json:"meetings"`
	Core_flags            []string               `json:"core_flags"`
	Syllabus_uri          string                 `json:"syllabus_uri"`
	Grade_distribution    []int                  `json:"grade_distribution"`
	Attributes            interface{}            `json:"attributes"`
}

type Professor struct {
	Id           IdWrapper   `json:"_id"`
	First_name   string      `json:"first_name"`
	Last_name    string      `json:"last_name"`
	Titles       []string    `json:"titles"`
	Email        string      `json:"email"`
	Phone_number string      `json:"phone_number"`
	Office       Location    `json:"office"`
	Profile_uri  string      `json:"profile_uri"`
	Image_uri    string      `json:"image_uri"`
	Office_hours []Meeting   `json:"office_hours"`
	Sections     []IdWrapper `json:"sections"`
}

type Organization struct {
	Id             IdWrapper `json:"_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Categories     []string  `json:"categories"`
	President_name string    `json:"president_name"`
	Emails         []string  `json:"emails"`
	Picture_data   string    `json:"picture_data"`
}
