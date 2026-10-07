package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransformProfessor(t *testing.T) {
	input := &DBProfessor{
		ID:          primitive.ObjectID{1},
		FirstName:   "Ada",
		LastName:    "Example",
		Titles:      []string{"Professor", "Department Chair"},
		Email:       "ada@example.com",
		PhoneNumber: "555-0100",
		Office:      &DBOffice{Building: "ECSS", MapURI: "https://example.com/office", Room: "3.100"},
		ProfileURI:  "https://example.com/profile",
		ImageURI:    "https://example.com/image",
		OfficeHours: []string{"Monday 10:00-11:00"},
		Sections:    []primitive.ObjectID{{2}, {3}},
	}
	want := &Professor{
		ID:          "010000000000000000000000",
		FirstName:   "Ada",
		LastName:    "Example",
		Titles:      []string{"Professor", "Department Chair"},
		Email:       "ada@example.com",
		PhoneNumber: "555-0100",
		Office:      &Office{Building: "ECSS", MapURI: "https://example.com/office", Room: "3.100"},
		ProfileURI:  "https://example.com/profile",
		ImageURI:    "https://example.com/image",
		OfficeHours: []string{"Monday 10:00-11:00"},
		Sections: []*Section{
			{ID: "020000000000000000000000"},
			{ID: "030000000000000000000000"},
		},
	}

	if diff := cmp.Diff(want, TransformProfessor(input)); diff != "" {
		t.Fatalf("TransformProfessor() mismatch (-want +got):\n%s", diff)
	}
}

func TestTransformProfessorNil(t *testing.T) {
	if got := TransformProfessor(nil); got != nil {
		t.Fatalf("TransformProfessor(nil) = %v, want nil", got)
	}
}

func TestTransformSection(t *testing.T) {
	startDate := time.Date(2026, time.August, 17, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, time.December, 11, 0, 0, 0, 0, time.UTC)
	input := &DBSection{
		ID:              primitive.ObjectID{2},
		SectionNumber:   "001",
		CourseReference: "course-example",
		SectionCorequisites: &DBCollectionRequirement{
			Type: "all", Name: "Lab", Required: 1, Options: []string{"CS 1100"},
		},
		AcademicSession: DBAcademicSession{Name: "Fall 2026", StartDate: startDate, EndDate: endDate},
		Professors:      []primitive.ObjectID{{1}, {3}},
		TeachingAssistants: []DBAssistant{
			{FirstName: "Grace", LastName: "Example", Role: "TA", Email: "grace@example.com"},
			{FirstName: "Alan", LastName: "Sample", Role: "Tutor", Email: "alan@example.com"},
		},
		InternalClassNumber: "12345",
		InstructionMode:     "In Person",
		Meetings: []DBMeeting{
			{
				StartDate: startDate, EndDate: endDate,
				MeetingDays: []string{"Monday", "Wednesday"}, StartTime: "10:00", EndTime: "11:15",
				Modality: "Lecture", Location: DBLocation{Building: "ECSS", Room: "2.100", MapURI: "https://example.com/lecture"},
			},
			{
				StartDate: startDate.AddDate(0, 0, 1), EndDate: endDate.AddDate(0, 0, -1),
				MeetingDays: []string{"Tuesday"}, StartTime: "14:00", EndTime: "15:00",
				Modality: "Lab", Location: DBLocation{Building: "ECSW", Room: "1.200", MapURI: "https://example.com/lab"},
			},
		},
		CoreFlags:         []string{"030", "090"},
		SyllabusURI:       "https://example.com/syllabus",
		GradeDistribution: []int32{10, 8, 3, 1, 0},
		Attributes:        map[string]string{"topic": "Computing"},
	}
	want := &Section{
		ID:              "020000000000000000000000",
		SectionNumber:   "001",
		CourseReference: &Course{ID: "course-example"},
		SectionCorequisites: &CollectionRequirement{
			Type: "all", Name: "Lab", Required: 1, Options: []string{"CS 1100"},
		},
		AcademicSession: &AcademicSession{Name: "Fall 2026", StartDate: startDate, EndDate: endDate},
		Professors: []*Professor{
			{ID: "010000000000000000000000"},
			{ID: "030000000000000000000000"},
		},
		TeachingAssistants: []*Assistant{
			{FirstName: "Grace", LastName: "Example", Role: "TA", Email: "grace@example.com"},
			{FirstName: "Alan", LastName: "Sample", Role: "Tutor", Email: "alan@example.com"},
		},
		InternalClassNumber: "12345",
		InstructionMode:     "In Person",
		Meetings: []*Meeting{
			{
				StartDate: startDate, EndDate: endDate,
				MeetingDays: []string{"Monday", "Wednesday"}, StartTime: "10:00", EndTime: "11:15",
				Modality: "Lecture", Location: &Location{Building: "ECSS", Room: "2.100", MapURI: "https://example.com/lecture"},
			},
			{
				StartDate: startDate.AddDate(0, 0, 1), EndDate: endDate.AddDate(0, 0, -1),
				MeetingDays: []string{"Tuesday"}, StartTime: "14:00", EndTime: "15:00",
				Modality: "Lab", Location: &Location{Building: "ECSW", Room: "1.200", MapURI: "https://example.com/lab"},
			},
		},
		CoreFlags:         []string{"030", "090"},
		SyllabusURI:       "https://example.com/syllabus",
		GradeDistribution: []int32{10, 8, 3, 1, 0},
		Attributes:        map[string]string{"topic": "Computing"},
	}

	if diff := cmp.Diff(want, TransformSection(input)); diff != "" {
		t.Fatalf("TransformSection() mismatch (-want +got):\n%s", diff)
	}
}

func TestTransformSectionNil(t *testing.T) {
	if got := TransformSection(nil); got != nil {
		t.Fatalf("TransformSection(nil) = %v, want nil", got)
	}
}
