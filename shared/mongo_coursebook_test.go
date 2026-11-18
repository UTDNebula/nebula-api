package shared

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestInvalidGradeQueriesDoNotAccessMongo(t *testing.T) {
	repo := &MongoGradesRepo{}
	for _, mode := range []string{"unsupported", "course_endpoint", "section_endpoint", "professor_endpoint"} {
		t.Run(mode, func(t *testing.T) {
			_, err := repo.Aggregate(context.Background(), GradeQuery{Mode: mode})
			if !errors.Is(err, ErrInvalidGradeQuery) {
				t.Fatalf("expected invalid query, got %v", err)
			}
		})
	}
}

func TestMongoCoursebookIntegration(t *testing.T) {
	uri := os.Getenv("NEBULA_TEST_MONGODB_URI")
	if uri == "" {
		t.Skip("set NEBULA_TEST_MONGODB_URI to run against a test MongoDB")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	database := client.Database("nebula_repo_test_" + bson.NewObjectID().Hex())
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if err := database.Drop(cleanupCtx); err != nil {
			t.Error(err)
		}
		if err := client.Disconnect(cleanupCtx); err != nil {
			t.Error(err)
		}
	})
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatal(err)
	}
	courseID, professorID := bson.NewObjectID(), bson.NewObjectID()
	sectionIDs := []bson.ObjectID{bson.NewObjectID(), bson.NewObjectID()}
	course := schema.Course{Id: courseID, Subject_prefix: "CS", Course_number: "3345", Internal_course_number: "fixture-course", Title: "Example course", Sections: sectionIDs}
	professor := schema.Professor{Id: professorID, First_name: "Jane", Last_name: "Doe", Sections: sectionIDs}
	sections := []schema.Section{
		{Id: sectionIDs[0], Section_number: "001", Course_reference: courseID, Professors: []bson.ObjectID{professorID}, Academic_session: schema.AcademicSession{Name: "26f"}, Grade_distribution: append([]int{3}, make([]int, 13)...)},
		{Id: sectionIDs[1], Section_number: "002", Course_reference: courseID, Professors: []bson.ObjectID{professorID}, Academic_session: schema.AcademicSession{Name: "26f"}, Grade_distribution: append([]int{4}, make([]int, 13)...)},
	}
	for _, fixture := range []struct {
		collection string
		records    []any
	}{
		{"courses", []any{course}}, {"professors", []any{professor}}, {"sections", []any{sections[0], sections[1]}},
		{"trends_course_sections", []any{bson.M{"_id": "CS3345", "sections": sections}}},
		{"trends_prof_sections", []any{bson.M{"first_name": "Jane", "last_name": "Doe", "sections": sections}}},
		{"trends_course_and_prof_sections", []any{bson.M{"marker": "combined", "sections": sections}}},
	} {
		if _, err := database.Collection(fixture.collection).InsertMany(ctx, fixture.records); err != nil {
			t.Fatal(err)
		}
	}
	courses := NewMongoCourseRepo(database.Collection("courses"))
	sectionRepo := NewMongoSectionRepo(database.Collection("sections"))
	professors := NewMongoProfessorRepo(database.Collection("professors"))
	t.Run("lookups", func(t *testing.T) {
		if got, err := courses.ByID(ctx, courseID); err != nil || got.Id != courseID {
			t.Fatalf("course: %+v, %v", got, err)
		}
		if got, err := sectionRepo.ByID(ctx, sectionIDs[0]); err != nil || got.Id != sectionIDs[0] {
			t.Fatalf("section: %+v, %v", got, err)
		}
		if got, err := professors.ByID(ctx, professorID); err != nil || got.Id != professorID {
			t.Fatalf("professor: %+v, %v", got, err)
		}
		for _, lookup := range []func(context.Context, bson.ObjectID) error{
			func(ctx context.Context, id bson.ObjectID) error { _, err := courses.ByID(ctx, id); return err },
			func(ctx context.Context, id bson.ObjectID) error { _, err := sectionRepo.ByID(ctx, id); return err },
			func(ctx context.Context, id bson.ObjectID) error { _, err := professors.ByID(ctx, id); return err },
		} {
			if err := lookup(ctx, bson.NewObjectID()); !errors.Is(err, ErrNotFound) {
				t.Fatalf("missing record: %v", err)
			}
			canceled, cancel := context.WithCancel(ctx)
			cancel()
			if err := lookup(canceled, courseID); !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled query: %v", err)
			}
		}
	})
	t.Run("search and all", func(t *testing.T) {
		if got, err := courses.Search(ctx, bson.M{"subject_prefix": "CS"}, 0, 1); err != nil || len(got) != 1 {
			t.Fatalf("courses: %v, %v", got, err)
		}
		if got, err := sectionRepo.Search(ctx, bson.M{}, 1, 1); err != nil || len(got) != 1 {
			t.Fatalf("sections: %v, %v", got, err)
		}
		if got, err := professors.Search(ctx, bson.M{"first_name": "Jane"}, 0, 1); err != nil || len(got) != 1 {
			t.Fatalf("professors: %v, %v", got, err)
		}
		if got, err := courses.All(ctx); err != nil || len(got) != 1 {
			t.Fatalf("all courses: %v, %v", got, err)
		}
		if got, err := professors.All(ctx); err != nil || len(got) != 1 {
			t.Fatalf("all professors: %v, %v", got, err)
		}
		if _, err := courses.Search(ctx, bson.M{"$invalid": true}, 0, 1); err == nil {
			t.Fatal("invalid MongoDB query succeeded")
		}
	})
	pagination := map[string]bson.D{"former_offset": {{Key: "$skip", Value: 0}}, "latter_offset": {{Key: "$skip", Value: 0}}, "limit": {{Key: "$limit", Value: 10}}}
	t.Run("relationships", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			want  int
			fetch func() (int, error)
		}{
			{"course sections", 2, func() (int, error) { return recordCount(courses.Sections(ctx, bson.M{"_id": courseID}, pagination)) }},
			{"course professors", 1, func() (int, error) { return recordCount(courses.Professors(ctx, bson.M{"_id": courseID}, pagination)) }},
			{"section courses", 1, func() (int, error) { return recordCount(sectionRepo.Courses(ctx, bson.M{}, pagination)) }},
			{"section professors", 1, func() (int, error) { return recordCount(sectionRepo.Professors(ctx, bson.M{}, pagination)) }},
			{"professor courses", 1, func() (int, error) { return recordCount(professors.Courses(ctx, bson.M{}, pagination)) }},
			{"professor sections", 2, func() (int, error) { return recordCount(professors.Sections(ctx, bson.M{}, pagination)) }},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if count, err := tc.fetch(); err != nil || count != tc.want {
					t.Fatalf("want %d, got %d, %v", tc.want, count, err)
				}
			})
		}
	})
	t.Run("grades", func(t *testing.T) {
		for _, tc := range []struct {
			fetch func(context.Context, bson.ObjectID) ([14]int, error)
			id    bson.ObjectID
			want  [14]int
		}{
			{courses.Grades, courseID, [14]int{7}}, {sectionRepo.Grades, sectionIDs[0], [14]int{3}}, {professors.Grades, professorID, [14]int{7}},
		} {
			if got, err := tc.fetch(ctx, tc.id); err != nil || got != tc.want {
				t.Fatalf("grades: %v, %v", got, err)
			}
		}
		repo := NewMongoGradesRepo(database)
		for _, mode := range []string{"overall", "semester", "section_type"} {
			got, err := repo.Aggregate(ctx, GradeQuery{Mode: mode, Prefix: "CS", Number: "3345", FirstName: "Jane", LastName: "Doe"})
			if err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "overall":
				if got.Overall != [14]int{7} {
					t.Fatalf("overall: %+v", got)
				}
			case "semester":
				if !reflect.DeepEqual(got.Semesters, []schema.GradeData{{Id: "26f", GradeDistribution: [14]int{7}}}) {
					t.Fatalf("semesters: %+v", got)
				}
			case "section_type":
				if len(got.SectionTypes) != 1 || len(got.SectionTypes[0].Data) != 1 || got.SectionTypes[0].Data[0].GradeDistribution != [14]int{7} {
					t.Fatalf("section types: %+v", got)
				}
			}
		}
	})
	t.Run("trends details", func(t *testing.T) {
		for _, fetch := range []func() ([]schema.Section, error){
			func() ([]schema.Section, error) { return courses.TrendsSections(ctx, bson.M{"_id": "CS3345"}) },
			func() ([]schema.Section, error) { return professors.TrendsSections(ctx, bson.M{"first_name": "Jane"}) },
			func() ([]schema.Section, error) {
				return NewMongoTrendsRepo(database.Collection("trends_course_and_prof_sections")).Sections(ctx, bson.M{"marker": "combined"})
			},
		} {
			got, err := fetch()
			if err != nil || len(got) != 2 {
				t.Fatalf("trends: %+v, %v", got, err)
			}
			for _, section := range got {
				if section.Course_details == nil || section.Professor_details == nil || len(*section.Course_details) != 1 || len(*section.Professor_details) != 1 {
					t.Fatalf("missing details: %+v", section)
				}
			}
		}
	})
}

func recordCount[T any](records []T, err error) (int, error) { return len(records), err }
