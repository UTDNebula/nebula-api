package shared

import (
	"context"
	"errors"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *MongoGradesRepo) Aggregate(ctx context.Context, query GradeQuery) (GradeResult, error) {
	flag := query.Mode
	switch flag {
	case "overall", "semester", "section_type":
	case "course_endpoint", "section_endpoint", "professor_endpoint":
		if query.ID == nil {
			return GradeResult{}, ErrInvalidGradeQuery
		}
	default:
		return GradeResult{}, ErrInvalidGradeQuery
	}
	var grades []schema.GradeData
	var results []map[string]any

	var sectionTypeGrades []schema.TypedGradeData // used to parse the response to section-type endpoints

	var cursor *mongo.Cursor
	var collection *mongo.Collection
	var pipeline mongo.Pipeline

	var sectionMatch bson.D
	var courseMatch bson.D
	var courseFind bson.D
	var professorMatch bson.D
	var professorFind bson.D

	var sampleCourse schema.Course // the sample course with the given prefix and course number parameter
	var sampleCourseFind bson.D    // the filter using prefix and course number to get sample course

	objId := query.ID

	var err error

	//	@TODO:	Recommend forcing using first_name and last_name to ensure single professors per query.
	// All professors sharing the name will be aggregated together in the current implementation
	prefix := query.Prefix
	number := query.Number
	section_number := query.SectionNumber
	first_name := query.FirstName
	last_name := query.LastName

	professor := (first_name != "" || last_name != "")

	// Find internal_course_number associated with subject_prefix and course_number, which will be used later on
	sampleCourseFind = bson.D{
		{Key: "subject_prefix", Value: prefix},
		{Key: "course_number", Value: number},
	}

	err = r.courses.FindOne(ctx, sampleCourseFind).Decode(&sampleCourse)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		// If the error is not that there is no matching documents,
		// throw an internal server error
		return GradeResult{}, err
	}
	internalCourseNumber := sampleCourse.Internal_course_number

	switch {
	case flag == "course_endpoint":
		// Filter on course ID, from the course endpoint
		collection = r.courses

		courseMatch := bson.D{{Key: "$match", Value: bson.M{"_id": objId}}}
		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage(), unwindSectionsStage(), projectGradeDistributionStage(flag, false), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case flag == "section_endpoint":
		// Filter on section ID, from section endpoint
		collection = r.sections

		sectionMatch := bson.D{{Key: "$match", Value: bson.M{"_id": objId}}}
		pipeline = mongo.Pipeline{sectionMatch, projectGradeDistributionStage(flag, true), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case flag == "professor_endpoint":
		// Filter on Professor from professor endpoint
		collection = r.professors

		professorMatch := bson.D{{Key: "$match", Value: bson.M{"_id": objId}}}
		pipeline = mongo.Pipeline{professorMatch, lookupSectionsStage(), unwindSectionsStage(), projectGradeDistributionStage(flag, false), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case prefix != "" && number == "" && section_number == "" && !professor:
		// Filter on Course
		collection = r.courses

		courseMatch = bson.D{{Key: "$match", Value: bson.M{"subject_prefix": prefix}}}
		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage(), unwindSectionsStage(), projectGradeDistributionStage(flag, false), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case prefix != "" && number != "" && section_number == "" && !professor:
		// Filter on Course
		collection = r.courses

		// Query using internal_course_number of the documents
		courseMatch := bson.D{{Key: "$match", Value: bson.M{"internal_course_number": internalCourseNumber}}}
		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage(), unwindSectionsStage(), projectGradeDistributionStage(flag, false), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case prefix != "" && number != "" && section_number != "" && !professor:
		// Filter on Course then Section
		collection = r.courses

		// Here we query all the courses with the given internal_couse_number,
		// and then filter on the section_number of those courses
		courseMatch := bson.D{{Key: "$match", Value: bson.M{"internal_course_number": internalCourseNumber}}}
		sectionMatch := bson.D{{Key: "$match", Value: bson.M{"sections.section_number": section_number}}}

		pipeline = mongo.Pipeline{courseMatch, lookupSectionsStage(), unwindSectionsStage(), sectionMatch, projectGradeDistributionStage(flag, false), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case prefix == "" && number == "" && section_number == "" && professor:
		// Filter on Professor
		collection = r.professors

		// Build professorMatch
		if last_name == "" {
			professorMatch = bson.D{{Key: "$match", Value: bson.M{"first_name": first_name}}}
		} else if first_name == "" {
			professorMatch = bson.D{{Key: "$match", Value: bson.M{"last_name": last_name}}}
		} else {
			professorMatch = bson.D{{Key: "$match", Value: bson.M{"first_name": first_name, "last_name": last_name}}}
		}

		pipeline = mongo.Pipeline{professorMatch, lookupSectionsStage(), unwindSectionsStage(), projectGradeDistributionStage(flag, false), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	case prefix != "" && professor:
		// Filter on Section by Matching Course and Professor IDs

		// Here we get the valid course ids and professor ids
		// and then we perform the grades aggregation against the sections collection,
		// matching on the course_reference and professor

		collection = r.sections

		// Find valid professor ids
		if last_name == "" {
			professorFind = bson.D{{Key: "first_name", Value: first_name}}
		} else if first_name == "" {
			professorFind = bson.D{{Key: "last_name", Value: last_name}}
		} else {
			professorFind = bson.D{{Key: "first_name", Value: first_name}, {Key: "last_name", Value: last_name}}
		}

		cursor, err = r.professors.Find(ctx, professorFind)
		if err != nil {
			return GradeResult{}, err
		}
		defer cursor.Close(ctx)
		if err = cursor.All(ctx, &results); err != nil {
			return GradeResult{}, err
		}

		profIDs := make([]bson.ObjectID, 0, len(results))
		for _, prof := range results {
			profID := prof["_id"].(bson.ObjectID)
			profIDs = append(profIDs, profID)
		}

		// Get valid course ids
		if number == "" {
			// if only the prefix is provided, filter only on the prefix
			courseFind = bson.D{{Key: "subject_prefix", Value: prefix}}
		} else {
			// if both prefix and course_number are provided, filter on internal_course_number
			courseFind = bson.D{{Key: "internal_course_number", Value: internalCourseNumber}}
		}

		cursor, err = r.courses.Find(ctx, courseFind)
		if err != nil {
			return GradeResult{}, err
		}
		defer cursor.Close(ctx)
		if err = cursor.All(ctx, &results); err != nil {
			return GradeResult{}, err
		}

		courseIDs := make([]bson.ObjectID, 0, len(results))
		for _, course := range results {
			courseID := course["_id"].(bson.ObjectID)
			courseIDs = append(courseIDs, courseID)
		}

		// Build sectionMatch
		if section_number == "" {
			sectionMatch =
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "course_reference", Value: bson.D{{Key: "$in", Value: courseIDs}}},
					{Key: "professors", Value: bson.D{{Key: "$in", Value: profIDs}}},
				}}}
		} else {
			sectionMatch =
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "course_reference", Value: bson.D{{Key: "$in", Value: courseIDs}}},
					{Key: "professors", Value: bson.D{{Key: "$in", Value: profIDs}}},
					{Key: "section_number", Value: section_number},
				}}}
		}

		pipeline = mongo.Pipeline{sectionMatch, projectGradeDistributionStage(flag, true), unwindGradeDistributionStage(), groupGradesStage(flag), sortGradesStage(flag), sumGradesStage(), groupGradeDistributionStage(flag)}

	default:
		return GradeResult{}, ErrInvalidGradeQuery
	}

	// if this is for section type, add the 2 additional stages to the pipeline
	if flag == "section_type" {
		pipeline = append(pipeline, sortGradeDistributionsStage(), groupSemesterGradeDistributionsStage())
	}

	// peform aggregation
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		return GradeResult{}, err
	}
	defer cursor.Close(ctx)

	// retrieve and parse all valid documents to appropriate type
	if flag != "section_type" {
		if err = cursor.All(ctx, &grades); err != nil {
			return GradeResult{}, err
		}
	} else {
		if err = cursor.All(ctx, &sectionTypeGrades); err != nil {
			return GradeResult{}, err
		}
	}

	switch flag {
	case "overall", "course_endpoint", "section_endpoint", "professor_endpoint":
		// combine all semester grade_distributions
		overallResponse := [14]int{}
		for _, sem := range grades {
			for i, grade := range sem.GradeDistribution {
				overallResponse[i] += grade
			}
		}
		return GradeResult{Overall: overallResponse}, nil
	case "semester":
		return GradeResult{Semesters: grades}, nil
	case "section_type":
		return GradeResult{SectionTypes: sectionTypeGrades}, nil
	}

	return GradeResult{}, ErrInvalidGradeQuery
}
