package shared

import "go.mongodb.org/mongo-driver/v2/bson"

func lookupSectionsStage() bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}},
	}
}

func unwindSectionsStage() bson.D {
	return bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$sections"}}}}
}

func projectGradeDistributionStage(flag string, withSection bool) bson.D {
	// Arrays of regexes and section types
	typeRegexes := [14]string{"0[0-9][0-9]", "0W[0-9]", "0H[0-9]", "0L[0-9]", "5H[0-9]", "1[0-9][0-9]", "2[0-9][0-9]", "3[0-9][0-9]", "5[0-9][0-9]", "6[0-9][0-9]", "7[0-9][0-9]", "HN[0-9]", "HON", "[0-9]U[0-9]"}
	typeStrings := [14]string{"0xx", "0Wx", "0Hx", "0Lx", "5Hx", "1xx", "2xx", "3xx", "5xx", "6xx", "7xx", "HNx", "HON", "xUx"}

	// Build the branches for the regex stage, each branch representing a section type
	var regexBranches [14]bson.D
	regexInput := "$sections.section_number"
	if withSection {
		regexInput = "$section_number"
	}
	for i, typeRegex := range typeRegexes {
		regexBranches[i] = bson.D{
			{Key: "case", Value: bson.D{
				{Key: "$regexMatch", Value: bson.D{
					{Key: "input", Value: regexInput},
					{Key: "regex", Value: typeRegex},
				}},
			}},
			{Key: "then", Value: typeStrings[i]},
		}
	}

	id, gradeDist := "$sections.academic_session.name", "$sections.grade_distribution"
	if withSection {
		id, gradeDist = "$academic_session.name", "$grade_distribution"
	}
	projectCrit := bson.D{
		{Key: "_id", Value: id},
		{Key: "grade_distribution", Value: gradeDist},
	}
	// Add the section_type for each section
	if flag == "section_type" {
		projectCrit = append(projectCrit, bson.E{Key: "section_type", Value: bson.D{
			{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: regexBranches},
				{Key: "default", Value: "OTHERS"}, // Cases in which section doesn't have the listed type
			}},
		}})
	}
	return bson.D{{Key: "$project", Value: projectCrit}}
}

func unwindGradeDistributionStage() bson.D {
	return bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$grade_distribution"},
			{Key: "includeArrayIndex", Value: "ix"},
		}},
	}
}

func groupGradesStage(flag string) bson.D {
	groupID := bson.D{
		{Key: "academic_session", Value: "$_id"},
		{Key: "ix", Value: "$ix"},
	}
	// Add section_type to _id to group grades by both academic_session and section_type
	if flag == "section_type" {
		groupID = append(groupID, bson.E{Key: "section_type", Value: "$section_type"})
	}
	return bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupID},
			{Key: "grades", Value: bson.D{{Key: "$push", Value: "$grade_distribution"}}},
		}},
	}
}

func sortGradesStage(flag string) bson.D {
	sortCrit := bson.D{
		{Key: "_id.ix", Value: 1},
		{Key: "_id", Value: 1},
	}
	// Insert the section-type criteria to index 1
	if flag == "section_type" {
		sortCrit = append(sortCrit[:1], append(bson.D{{Key: "_id.section_type", Value: 1}}, sortCrit[1:]...)...)
	}
	return bson.D{{Key: "$sort", Value: sortCrit}}
}

func sumGradesStage() bson.D {
	return bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "grades", Value: bson.D{
				{Key: "$sum", Value: "$grades"},
			}},
		}},
	}
}

// Stage to group grade distribution
func groupGradeDistributionStage(flag string) bson.D {
	var groupDistributionID any = "$_id.academic_session"
	// Add the section-type criteria
	if flag == "section_type" {
		groupDistributionID = bson.D{
			{Key: "academic_section", Value: "$_id.academic_session"},
			{Key: "section_type", Value: "$_id.section_type"},
		}
	}
	return bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupDistributionID},
			{Key: "grade_distribution", Value: bson.D{
				{Key: "$push", Value: "$grades"},
			}},
		}},
	}
}

// Additional stages for "section-type" pipeline
// Stage to sort the section-type-specific grade distributions before grouping
func sortGradeDistributionsStage() bson.D {
	return bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "_id.section_type", Value: 1},
			{Key: "_id", Value: 1},
		}},
	}
}

// Stage to group section-type-specific grade distributions together based on semester
func groupSemesterGradeDistributionsStage() bson.D {
	return bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.academic_section"},
			{Key: "data", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "type", Value: "$_id.section_type"},
				{Key: "grade_distribution", Value: "$grade_distribution"},
			}}}},
		}},
	}
}
