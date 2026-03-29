package model

import "time"

/*
These are the models that interact with the MongoDB
*/

type DBCollectionRequirement struct {
	Type     string `bson:"type"`
	Name     string `bson:"name"`
	Required int32  `bson:"required"`
	Options  any    `bson:"options"`
}

func transformColReq(dbColReq *DBCollectionRequirement) *CollectionRequirement {
	if dbColReq == nil {
		return nil
	}
	return &CollectionRequirement{
		dbColReq.Type,
		dbColReq.Name,
		dbColReq.Required,
		dbColReq.Options,
	}
}

// NOTE: For now, course on Mongo side is identical to course on GraphQL side
// However, in the near future, when we implement the REFERENCE RESOLVER, there will be diffs
type DBCourse struct {
	ID                     string                   `bson:"_id"`
	SubjectPrefix          string                   `bson:"subject_prefix"`
	CourseNumber           string                   `bson:"course_number"`
	Title                  string                   `bson:"title"`
	Description            string                   `bson:"description"`
	EnrollmentReqs         string                   `bson:"enrollment_reqs"`
	School                 string                   `bson:"school"`
	CreditHours            string                   `bson:"credit_hours"`
	ClassLevel             string                   `bson:"class_level"`
	ActivityType           string                   `bson:"activity_type"`
	Grading                string                   `bson:"grading"`
	InternalCourseNumber   string                   `bson:"internal_course_number"`
	Prerequisites          *DBCollectionRequirement `bson:"prerequisites"`
	Corequisites           *DBCollectionRequirement `bson:"corequisites"`
	CoOrPreRequisites      *DBCollectionRequirement `bson:"co_or_pre_requisites"`
	Sections               []string                 `bson:"sections"`
	LectureContactHours    string                   `bson:"lecture_contact_hours"`
	LaboratoryContactHours string                   `bson:"laboratory_contact_hours"`
	OfferingFrequency      string                   `bson:"offering_frequency"`
	CatalogYear            string                   `bson:"catalog_year"`
	Attributes             any                      `bson:"attributes"`
}

// Transform the course object that interacts with Mongo to the course object that interacts with GraphQL
func TransformCourse(dbCourse *DBCourse) *Course {
	if dbCourse == nil {
		return nil
	}
	return &Course{
		dbCourse.ID,
		dbCourse.SubjectPrefix,
		dbCourse.CourseNumber,
		dbCourse.Title,
		dbCourse.Description,
		dbCourse.EnrollmentReqs,
		dbCourse.School,
		dbCourse.CreditHours,
		dbCourse.ClassLevel,
		dbCourse.ActivityType,
		dbCourse.Grading,
		dbCourse.InternalCourseNumber,
		transformColReq(dbCourse.Prerequisites),
		transformColReq(dbCourse.Corequisites),
		transformColReq(dbCourse.Corequisites),
		dbCourse.Sections,
		dbCourse.LectureContactHours,
		dbCourse.LaboratoryContactHours,
		dbCourse.OfferingFrequency,
		dbCourse.CatalogYear,
		dbCourse.Attributes,
	}
}

// DBEvent represents the database model for an event.
type DBEvent struct {
	ID                 string    `bson:"_id"`
	Summary            string    `bson:"summary"`
	Location           string    `bson:"location"`
	StartTime          time.Time `bson:"start_time"`
	EndTime            time.Time `bson:"end_time"`
	Description        string    `bson:"description"`
	EventType          []string  `bson:"event_type"`
	TargetAudience     []string  `bson:"target_audience"`
	Topic              []string  `bson:"topic"`
	EventTags          []string  `bson:"event_tags"`
	EventWebsite       string    `bson:"event_website"`
	Department         []string  `bson:"department"`
	ContactName        string    `bson:"contact_name"`
	ContactEmail       string    `bson:"contact_email"`
	ContactPhoneNumber string    `bson:"contact_phone_number"`
}

// DBCometCalendarRoom represents the database model for a room with its events.
type DBCometCalendarRoom struct {
	Room   string    `bson:"room"`
	Events []DBEvent `bson:"events"`
}

// DBCometCalendarBuilding represents the database model for a building with its rooms.
type DBCometCalendarBuilding struct {
	Building string                `bson:"building"`
	Rooms    []DBCometCalendarRoom `bson:"rooms"`
}

// DBCometCalendar represents the database model for a calendar day.
type DBCometCalendar struct {
	ID        string                    `bson:"_id"`
	Date      string                    `bson:"date"`
	Buildings []DBCometCalendarBuilding `bson:"buildings"`
}

// TransformEvent converts a database event model into a GraphQL Event type.
func TransformEvent(dbEvent *DBEvent) *Event {
	return &Event{
		ID:                 dbEvent.ID,
		Summary:            dbEvent.Summary,
		Location:           dbEvent.Location,
		StartTime:          dbEvent.StartTime,
		EndTime:            dbEvent.EndTime,
		Description:        dbEvent.Description,
		EventType:          dbEvent.EventType,
		TargetAudience:     dbEvent.TargetAudience,
		Topic:              dbEvent.Topic,
		EventTags:          dbEvent.EventTags,
		EventWebsite:       dbEvent.EventWebsite,
		Department:         dbEvent.Department,
		ContactName:        dbEvent.ContactName,
		ContactEmail:       dbEvent.ContactEmail,
		ContactPhoneNumber: dbEvent.ContactPhoneNumber,
	}
}

// TransformRoom converts a database room model into a GraphQL CometCalendarRoom type.
func TransformRoom(dbRoom *DBCometCalendarRoom) *CometCalendarRoom {
	events := make([]*Event, len(dbRoom.Events))
	for i, e := range dbRoom.Events {
		// Note: e is a value, so we take its address for the slice
		events[i] = TransformEvent(&e)
	}
	return &CometCalendarRoom{
		Room:   dbRoom.Room,
		Events: events,
	}
}

// TransformBuilding converts a database building model into a GraphQL CometCalendarBuilding type.
func TransformBuilding(dbBuilding *DBCometCalendarBuilding) *CometCalendarBuilding {
	rooms := make([]*CometCalendarRoom, len(dbBuilding.Rooms))
	for i, r := range dbBuilding.Rooms {
		rooms[i] = TransformRoom(&r)
	}
	return &CometCalendarBuilding{
		Building: dbBuilding.Building,
		Rooms:    rooms,
	}
}

// TransformCalendar converts a database calendar model into a GraphQL CometCalendar type.
func TransformCalendar(dbCometCalendar *DBCometCalendar) *CometCalendar {
	buildings := make([]*CometCalendarBuilding, len(dbCometCalendar.Buildings))
	for i, b := range dbCometCalendar.Buildings {
		buildings[i] = TransformBuilding(&b)
	}
	return &CometCalendar{
		ID:        dbCometCalendar.ID,
		Date:      dbCometCalendar.Date,
		Buildings: buildings,
	}

// --- Section DB Models and Transformers ---

type DBAcademicSession struct {
	Name      string    `bson:"name"`
	StartDate time.Time `bson:"start_date"`
	EndDate   time.Time `bson:"end_date"`
}

type DBAssistant struct {
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Role      string `bson:"role"`
	Email     string `bson:"email"`
}

type DBLocation struct {
	Building string `bson:"building"`
	Room     string `bson:"room"`
	MapUri   string `bson:"map_uri"`
}

type DBMeeting struct {
	StartDate   time.Time  `bson:"start_date"`
	EndDate     time.Time  `bson:"end_date"`
	MeetingDays []string   `bson:"meeting_days"`
	StartTime   string     `bson:"start_time"`
	EndTime     string     `bson:"end_time"`
	Modality    string     `bson:"modality"`
	Location    DBLocation `bson:"location"`
}

type DBSection struct {
	ID                  string                   `bson:"_id"`
	SectionNumber       string                   `bson:"section_number"`
	CourseReference     string                   `bson:"course_reference"`
	SectionCorequisites *DBCollectionRequirement `bson:"section_corequisites"`
	AcademicSession     DBAcademicSession        `bson:"academic_session"`
	Professors          []string                 `bson:"professors"`
	TeachingAssistants  []DBAssistant            `bson:"teaching_assistants"`
	InternalClassNumber string                   `bson:"internal_class_number"`
	InstructionMode     string                   `bson:"instruction_mode"`
	Meetings            []DBMeeting              `bson:"meetings"`
	CoreFlags           []string                 `bson:"core_flags"`
	SyllabusUri         string                   `bson:"syllabus_uri"`
	GradeDistribution   []int32                  `bson:"grade_distribution"`
	Attributes          any                      `bson:"attributes"`
}

func transformAcademicSession(dbSession *DBAcademicSession) *AcademicSession {
	if dbSession == nil {
		return nil
	}
	return &AcademicSession{
		Name:      dbSession.Name,
		StartDate: dbSession.StartDate,
		EndDate:   dbSession.EndDate,
	}
}

func transformAssistant(dbAssistant *DBAssistant) *Assistant {
	if dbAssistant == nil {
		return nil
	}
	return &Assistant{
		FirstName: dbAssistant.FirstName,
		LastName:  dbAssistant.LastName,
		Role:      dbAssistant.Role,
		Email:     dbAssistant.Email,
	}
}

func transformLocation(dbLocation *DBLocation) *Location {
	if dbLocation == nil {
		return nil
	}
	return &Location{
		Building: dbLocation.Building,
		Room:     dbLocation.Room,
		MapUri:   dbLocation.MapUri,
	}
}

func transformMeeting(dbMeeting *DBMeeting) *Meeting {
	if dbMeeting == nil {
		return nil
	}
	return &Meeting{
		StartDate:   dbMeeting.StartDate,
		EndDate:     dbMeeting.EndDate,
		MeetingDays: dbMeeting.MeetingDays,
		StartTime:   dbMeeting.StartTime,
		EndTime:     dbMeeting.EndTime,
		Modality:    dbMeeting.Modality,
		Location:    transformLocation(&dbMeeting.Location),
	}
}

// TransformSection converts a database section model into a GraphQL Section type.
func TransformSection(dbSection *DBSection) *Section {
	if dbSection == nil {
		return nil
	}

	assistants := make([]*Assistant, len(dbSection.TeachingAssistants))
	for i, a := range dbSection.TeachingAssistants {
		// Taking address since `a` is a value in the range loop
		assistants[i] = transformAssistant(&a)
	}

	meetings := make([]*Meeting, len(dbSection.Meetings))
	for i, m := range dbSection.Meetings {
		meetings[i] = transformMeeting(&m)
	}

	return &Section{
		ID:                  dbSection.ID,
		SectionNumber:       dbSection.SectionNumber,
		CourseReference:     dbSection.CourseReference,
		SectionCorequisites: transformColReq(dbSection.SectionCorequisites),
		AcademicSession:     transformAcademicSession(&dbSection.AcademicSession),
		Professors:          dbSection.Professors,
		TeachingAssistants:  assistants,
		InternalClassNumber: dbSection.InternalClassNumber,
		InstructionMode:     dbSection.InstructionMode,
		Meetings:            meetings,
		CoreFlags:           dbSection.CoreFlags,
		SyllabusUri:         dbSection.SyllabusUri,
		GradeDistribution:   dbSection.GradeDistribution,
		Attributes:          dbSection.Attributes,
	}
}
