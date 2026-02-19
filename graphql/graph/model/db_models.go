package model

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
