package schema

import (
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	Id                       primitive.ObjectID     `bson:"_id" json:"_id"`
	Subject_prefix           string                 `bson:"subject_prefix" json:"subject_prefix" queryable:""`
	Course_number            string                 `bson:"course_number" json:"course_number" queryable:""`
	Title                    string                 `bson:"title" json:"title" queryable:""`
	Description              string                 `bson:"description" json:"description"`
	Enrollment_reqs          string                 `bson:"enrollment_reqs" json:"enrollment_reqs"`
	School                   string                 `bson:"school" json:"school" queryable:""`
	Credit_hours             string                 `bson:"credit_hours" json:"credit_hours" queryable:""`
	Class_level              string                 `bson:"class_level" json:"class_level" queryable:""`
	Activity_type            string                 `bson:"activity_type" json:"activity_type" queryable:""`
	Grading                  string                 `bson:"grading" json:"grading" queryable:""`
	Internal_course_number   string                 `bson:"internal_course_number" json:"internal_course_number" queryable:""`
	Prerequisites            *CollectionRequirement `bson:"prerequisites" json:"prerequisites"`
	Corequisites             *CollectionRequirement `bson:"corequisites" json:"corequisites"`
	Co_or_pre_requisites     *CollectionRequirement `bson:"co_or_pre_requisites" json:"co_or_pre_requisites"`
	Sections                 []primitive.ObjectID   `bson:"sections" json:"sections"`
	Lecture_contact_hours    string                 `bson:"lecture_contact_hours" json:"lecture_contact_hours" queryable:""`
	Laboratory_contact_hours string                 `bson:"laboratory_contact_hours" json:"laboratory_contact_hours" queryable:""`
	Offering_frequency       string                 `bson:"offering_frequency" json:"offering_frequency" queryable:""`
	Catalog_year             string                 `bson:"catalog_year" json:"catalog_year" queryable:""`
	Attributes               interface{}            `bson:"attributes" json:"attributes"`
}

type BasicCourse struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Subject_prefix string             `bson:"subject_prefix" json:"subject_prefix" queryable:""`
	Course_number  string             `bson:"course_number" json:"course_number" queryable:""`
	Title          string             `bson:"title" json:"title" queryable:""`
	Credit_hours   string             `bson:"credit_hours" json:"credit_hours" queryable:""`
	Class_level    string             `bson:"class_level" json:"class_level" queryable:""`
	Activity_type  string             `bson:"activity_type" json:"activity_type" queryable:""`
	Catalog_year   string             `bson:"catalog_year" json:"catalog_year" queryable:""`
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
	Id                    primitive.ObjectID     `bson:"_id" json:"_id"`
	Section_number        string                 `bson:"section_number" json:"section_number" queryable:""`
	Course_reference      primitive.ObjectID     `bson:"course_reference" json:"course_reference" queryable:""`
	Section_corequisites  *CollectionRequirement `bson:"section_corequisites" json:"section_corequisites"`
	Academic_session      AcademicSession        `bson:"academic_session" json:"academic_session"`
	Professors            []primitive.ObjectID   `bson:"professors" json:"professors"`
	Teaching_assistants   []Assistant            `bson:"teaching_assistants" json:"teaching_assistants"`
	Internal_class_number string                 `bson:"internal_class_number" json:"internal_class_number" queryable:""`
	Instruction_mode      string                 `bson:"instruction_mode" json:"instruction_mode" queryable:""`
	Meetings              []Meeting              `bson:"meetings" json:"meetings"`
	Core_flags            []string               `bson:"core_flags" json:"core_flags"`
	Syllabus_uri          string                 `bson:"syllabus_uri" json:"syllabus_uri"`
	Grade_distribution    []int                  `bson:"grade_distribution" json:"grade_distribution"`
	Attributes            interface{}            `bson:"attributes" json:"attributes"`
	Professor_details     *[]BasicProfessor      `bson:"professor_details,omitempty" json:"professor_details,omitempty"` // only shows if professor_details was set by the endpoint
	Course_details        *[]BasicCourse         `bson:"course_details,omitempty" json:"course_details,omitempty"`       // only shows if course_details was set by the endpoint
}

type Professor struct {
	Id           primitive.ObjectID   `bson:"_id" json:"_id"`
	First_name   string               `bson:"first_name" json:"first_name" queryable:""`
	Last_name    string               `bson:"last_name" json:"last_name" queryable:""`
	Titles       []string             `bson:"titles" json:"titles" queryable:""`
	Email        string               `bson:"email" json:"email" queryable:""`
	Phone_number string               `bson:"phone_number" json:"phone_number" queryable:""`
	Office       Location             `bson:"office" json:"office"`
	Profile_uri  string               `bson:"profile_uri" json:"profile_uri"`
	Image_uri    string               `bson:"image_uri" json:"image_uri"`
	Office_hours []Meeting            `bson:"office_hours" json:"office_hours"`
	Sections     []primitive.ObjectID `bson:"sections" json:"sections"`
}

type BasicProfessor struct {
	Id           primitive.ObjectID `bson:"_id" json:"_id"`
	First_name   string             `bson:"first_name" json:"first_name" queryable:""`
	Last_name    string             `bson:"last_name" json:"last_name" queryable:""`
	Email        string             `bson:"email" json:"email" queryable:""`
	Phone_number string             `bson:"phone_number" json:"phone_number" queryable:""`
	Office       Location           `bson:"office" json:"office"`
	Office_hours []Meeting          `bson:"office_hours" json:"office_hours"`
}

type Organization struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Title          string             `bson:"title" json:"title"`
	Description    string             `bson:"description" json:"description"`
	Categories     []string           `bson:"categories" json:"categories"`
	President_name string             `bson:"president_name" json:"president_name"`
	Emails         []string           `bson:"emails" json:"emails"`
	Picture_data   string             `bson:"picture_data" json:"picture_data"`
}

type DiscountProgram struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Category string             `bson:"category" json:"category"`
	Business string             `bson:"business" json:"business"`
	Address  string             `bson:"address" json:"address"`
	Phone    string             `bson:"phone" json:"phone"`
	Email    string             `bson:"email" json:"email"`
	Website  string             `bson:"website" json:"website"`
	Discount string             `bson:"discount" json:"discount"`
}

type Event struct {
	Id                 primitive.ObjectID `bson:"_id" json:"_id"`
	Summary            string             `bson:"summary" json:"summary"`
	Location           string             `bson:"location" json:"location"`
	StartTime          time.Time          `bson:"start_time" json:"start_time"`
	EndTime            time.Time          `bson:"end_time" json:"end_time"`
	Description        string             `bson:"description" json:"description"`
	EventType          []string           `bson:"event_type" json:"event_type"`
	TargetAudience     []string           `bson:"target_audience" json:"target_audience"`
	Topic              []string           `bson:"topic" json:"topic"`
	EventTags          []string           `bson:"event_tags" json:"event_tags"`
	EventWebsite       string             `bson:"event_website" json:"event_website"`
	Department         []string           `bson:"department" json:"department"`
	ContactName        string             `bson:"contact_name" json:"contact_name"`
	ContactEmail       string             `bson:"contact_email" json:"contact_email"`
	ContactPhoneNumber string             `bson:"contact_phone_number" json:"contact_phone_number"`
}

// Event hierarchy
type MultiBuildingEvents[T any] struct {
	Date      string                    `bson:"date" json:"date"`
	Buildings []SingleBuildingEvents[T] `bson:"buildings" json:"buildings"`
}

type SingleBuildingEvents[T any] struct {
	Building string          `bson:"building" json:"building"`
	Rooms    []RoomEvents[T] `bson:"rooms" json:"rooms"`
}

type RoomEvents[T any] struct {
	Room   string `bson:"room" json:"room"`
	Events []T    `bson:"events" json:"events"`
}

// Event types
type SectionWithTime struct {
	Section   primitive.ObjectID `bson:"section" json:"section"`
	StartTime string             `bson:"start_time" json:"start_time"`
	EndTime   string             `bson:"end_time" json:"end_time"`
}

type AstraEvent struct {
	ActivityName        *string  `bson:"activity_name" json:"activity_name"`
	MeetingType         *string  `bson:"meeting_type" json:"meeting_type"`
	StartDate           *string  `bson:"start_date" json:"start_date"`
	EndDate             *string  `bson:"end_date" json:"end_date"`
	CurrentState        *string  `bson:"current_state" json:"current_state"`
	NotAllowedUsageMask *float64 `bson:"not_allowed_usage_mask" json:"not_allowed_usage_mask"`
	UsageColor          *string  `bson:"usage_color" json:"usage_color"`
	Capacity            *float64 `bson:"capacity" json:"capacity"`
}
type MazevoEvent struct {
	EventName         *string  `bson:"eventName" json:"eventName"`
	OrganizationName  *string  `bson:"organizationName" json:"organizationName"`
	ContactName       *string  `bson:"contactName" json:"contactName"`
	SetupMinutes      *float64 `bson:"setupMinutes" json:"setupMinutes"`
	DateTimeStart     *string  `bson:"dateTimeStart" json:"dateTimeStart"`
	DateTimeEnd       *string  `bson:"dateTimeEnd" json:"dateTimeEnd"`
	TeardownMinutes   *float64 `bson:"teardownMinutes" json:"teardownMinutes"`
	StatusDescription *string  `bson:"statusDescription" json:"statusDescription"`
	StatusColor       *string  `bson:"statusColor" json:"statusColor"`
}

// Rooms type
type BuildingRooms struct {
	Building string  `bson:"building" json:"building"`
	Rooms    []Room  `bson:"rooms" json:"rooms"`
	Lat      float64 `bson:"lat" json:"lat"`
	Lng      float64 `bson:"lng" json:"lng"`
}
type Room struct {
	Room     string `bson:"room" json:"room"`
	Capacity int    `bson:"capacity" json:"capacity"`
}

// Map location type
type MapBuilding struct {
	Name    *string  `bson:"name" json:"name"`
	Acronym *string  `bson:"acronym" json:"acronym"`
	Lat     *float64 `bson:"lat" json:"lat"`
	Lng     *float64 `bson:"lng" json:"lng"`
}

type GradeData struct {
	Id                string  `bson:"_id" json:"_id"`
	GradeDistribution [14]int `bson:"grade_distribution" json:"grade_distribution"`
}

type TypedGradeData struct {
	Id   string `bson:"_id" json:"_id"`
	Data []struct {
		Type              string  `bson:"type" json:"type"`
		GradeDistribution [14]int `bson:"grade_distribution" json:"grade_distribution"`
	} `bson:"data" json:"data"`
}

// Prefix used for cloud storage bucket names
const BUCKET_PREFIX = "utdnebula_"

// Minimized form of storage.BucketAttrs for cloud storage
type BucketInfo struct {
	Name     string    `bson:"name" json:"name"`
	Created  time.Time `bson:"created" json:"created"`
	Updated  time.Time `bson:"updated" json:"updated"`
	Contents []string  `bson:"contents" json:"contents"`
}

func BucketInfoFromAttrs(attrs *storage.BucketAttrs) BucketInfo {
	// Don't show the bucket prefix externally
	bucketName, _ := strings.CutPrefix(attrs.Name, BUCKET_PREFIX)
	return BucketInfo{bucketName, attrs.Created, attrs.Updated, []string{}}
}

// Minimized form of storage.ObjectAttrs for cloud storage
type ObjectInfo struct {
	Bucket          string    `bson:"bucket" json:"bucket"`
	Name            string    `bson:"name" json:"name"`
	ContentType     string    `bson:"content_type" json:"content_type"`
	Size            int64     `bson:"size" json:"size"`
	ContentEncoding string    `bson:"content_encoding" json:"content_encoding"`
	MD5             []byte    `bson:"md5" json:"md5"`
	MediaLink       string    `bson:"media_link" json:"media_link"`
	Created         time.Time `bson:"created" json:"created"`
	Updated         time.Time `bson:"updated" json:"updated"`
	PublicUrl       string    `bson:"public_url" json:"public_url"`
}

func ObjectInfoFromAttrs(attrs *storage.ObjectAttrs, url string) ObjectInfo {
	// Don't show the bucket prefix externally
	bucketName, _ := strings.CutPrefix(attrs.Bucket, BUCKET_PREFIX)
	return ObjectInfo{
		bucketName,
		attrs.Name,
		attrs.ContentType,
		attrs.Size,
		attrs.ContentEncoding,
		attrs.MD5,
		attrs.MediaLink,
		attrs.Created,
		attrs.Updated,
		url,
	}
}

// Signed URL request body
type ObjectSignedURLBody struct {
	Method     string   `json:"method"`     // method to be used with signed URL. For example, PUT
	Headers    []string `json:"headers"`    // headers for signed URL
	Expiration string   `json:"expiration"` // timestamp for when the signed URL will expire
}

// Letters type
type Letters struct {
	Date    string `bson:"date" json:"date"`
	Letters string `bson:"letters" json:"letters"`
}

// Academic Calendar type
type AcademicCalendar struct {
	Id                       string                    `bson:"_id" json:"_id"`
	Timeline                 string                    `bson:"timeline" json:"timeline"`
	Sessions                 []AcademicCalendarSession `bson:"sessions" json:"sessions"`
	EnrollmentOpens          string                    `bson:"enrollment_opens" json:"enrollment_opens"`
	SchedulePlannerAvailable string                    `bson:"schedule_planner_available" json:"schedule_planner_available"`
	OnlineAddSwapEnds        string                    `bson:"online_add_swap_ends" json:"online_add_swap_ends"`
	LastReadmission          string                    `bson:"last_readmission" json:"last_readmission"`
	LastFromWaitlist         string                    `bson:"last_from_waitlist" json:"last_from_waitlist"`
	MidtermsDue              string                    `bson:"midterms_due" json:"midterms_due"`
	UniversityClosings       [][]string                `bson:"university_closings" json:"university_closings"`
	NoClasses                [][]string                `bson:"no_classes" json:"no_classes"`
}
type AcademicCalendarSession struct {
	Name               string                        `bson:"name" json:"name"`
	LastRegistration   string                        `bson:"last_registration" json:"last_registration"`
	Begin              string                        `bson:"begin" json:"begin"`
	LateRegistration   []string                      `bson:"late_registration" json:"late_registration"`
	CensusDay          string                        `bson:"census_day" json:"census_day"`
	DropDeadlines      AcademicCalendarDropDeadlines `bson:"drop_deadlines" json:"drop_deadlines"`
	End                string                        `bson:"end" json:"end"`
	ReadingDays        []string                      `bson:"reading_days" json:"reading_days"`
	Exams              []string                      `bson:"exams" json:"exams"`
	FinalGradingPeriod []string                      `bson:"final_grading_period" json:"final_grading_period"`
}
type AcademicCalendarDropDeadlines struct {
	WithoutW                  string `bson:"without_w" json:"without_w"`
	UndergradApprovalRequired string `bson:"undergrad_approval_required" json:"undergrad_approval_required"`
	GraduateWithdrawlEnds     string `bson:"graduate_withdrawl_ends" json:"graduate_withdrawl_ends"`
}

// Type for all API responses
type APIResponse[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}


// SearchResponse handles the 1:1 mapping of the discovery/search endpoint
type SearchResponse struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	TotalFound int               `json:"totalFound"` // Total matching programs (e.g., 27)
	Start      int               `json:"start"`      // Offset for pagination
	Limit      int               `json:"limit"`      // Results per page
	// This slice contains the summary objects used to find the IDs
	Items      []ProgramSummary  `json:"items"`
}

// Program Schema
type Program struct {
	ProgramID        int               `json:"programId"`
	DateCreated      string            `json:"dateCreated"`
	DateModified     string            `json:"dateModified"`
	BannerImageID    int               `json:"bannerImageFileId"`
	BannerImageTxt   string            `json:"bannerImageAltText"`
	BrochureSections []BrochureSection `json:"sections"`
	Costs            Cost              `json:"-"`
}

type BrochureSection struct {
	ID          int      `json:"sectionId"`
	DisplayName string   `json:"sectionDisplayName"`
	Ordinal     int      `json:"sectionOrdinal"`
	Default     int      `json:"sectionIsDefault"`
	Widgets     []Widget `json:"sectionWidgets"`
}
type Widget struct {
	Cols                    int               `json:"cols"`
	Rows                    int               `json:"rows"`
	X                       int               `json:"x"`
	Y                       int               `json:"y"`
	MinItemCols             int               `json:"minItemCols"`
	ContentID               int               `json:"contentId"`
	ContentType             string            `json:"contentType"`
	ContentHTML             string            `json:"contentHTML"`
	HeaderText              string            `json:"headerText"`
	ContentInformationSheet *InformationSheet `json:"contentInformationSheet,omitempty"` // Often "". Making it a ptr will help w/ null handling
	ContentMedia            *Media            `json:"contentMedia,omitempty"`
}

// Widget structs
type InformationSheet struct {
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	ParameterId       int      `json:"parameterId"`
	ParameterName     string   `json:"parameterName"`
	ParameterType     string   `json:"parameterType"`
	ParameterOrdinal  int      `json:"parameterOrdinal"`
	ParameterGlossary string   `json:"parameterGlossary"`
	AssignedValues    []string `json:"assignedValues"`
}

type Media struct {
	Type          string        `json:"type"`
	ImageId       int           `json:"imageId"`
	ImageAltText  string        `json:"imageAltText"`
	EmbedCode     string        `json:"embedCode"`
	MapProperties MapProperties `json:"mapProperties"`
}

type MapProperties struct {
	ZoomLevel       int     `json:"zoomLevel"`
	CenterLatitude  float64 `json:"centerLatitude"`
	CenterLongitude float64 `json:"centerLongitude"`
}

// Cost Structs
type Cost struct {
	CostSheets            []CostSheetObject `json:"costSheets"`
	NextAppCycleCostSheet NextAppCycle      `json:"nextAppCycleCostSheet"`
}

type CostSheetObject struct {
	CostSheetID         int    `json:"costSheetId"`
	Term                string `json:"term"`
	Year                int    `json:"year"`
	ProgramID           int    `json:"programId"`
	Public              bool   `json:"public"`
	DualYear            bool   `json:"dualYear"`
	ApplicationDeadline string `json:"applicationDeadline"`
}

type NextAppCycle struct {
	ID                        int        `json:"costSheetId"`
	Term                      string     `json:"term"`
	Year                      int        `json:"year"`
	ProgramID                 int        `json:"programId"`
	CostSheetNotes            string     `json:"costSheetNotes"`
	Public                    bool       `json:"public"`
	DtCreated                 string     `json:"dtCreated"`
	DtModified                string     `json:"dtModified"`
	UserID                    int        `json:"userId"`
	IsAdmin                   bool       `json:"isAdmin"`
	BillableCostSheetItems    []CostItem `json:"billableCostSheetItems"`
	NonBillableCostSheetItems []CostItem `json:"nonBillableCostSheetItems"`
	CreditCostSheetItems      []CostItem `json:"creditCostSheetItems"`
	LocaleCurrencyData        LocaleData `json:"localeCurrencyData"`
}

// NextAppCycle helper struct
type CostItem struct {
	CostSheetItemID       int          `json:"costSheetItemId"`
	CostSheetItemName     string       `json:"costSheetItemName"`
	CostSheetItemCategory string       `json:"costSheetItemCategory"`
	CostSheetItemType     string       `json:"costSheetItemType"`
	CostSheetItemGlossary string       `json:"costSheetItemGlossary"`
	CostSheetItemHint     string       `json:"costSheetItemHint"`
	Ordinal               int          `json:"ordinal"`
	Costs                 []CostDetail `json:"costs"`
}

type CostDetail struct {
	CostKey      string  `json:"costKey"`
	CostValue    float64 `json:"costValue"`
	CostCurrency string  `json:"costCurrency"`
}

type LocaleData struct {
	Locale                          string `json:"locale"`
	ThreeLetterCurrencyAbbreviation string `json:"threeLetterCurrencyAbbreviation"`
}

/* Can uncomment these if we ever get evals

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
	Id                   primitive.ObjectID `bson:"_id" json:"_id"`
	CourseExperience     []EvaluationField  `bson:"course_experience" json:"course_experience"`
	InstructorExperience []EvaluationField  `bson:"instructor_experience" json:"instructor_experience"`
	StudentExperience    []EvaluationField  `bson:"student_experience" json:"student_experience"`
}

*/
