package controllers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/UTDNebula/nebula-api/rest/routes"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TestGradesRepo struct {
	Result  shared.GradeResult
	Err     error
	Calls   int
	Query   shared.GradeQuery
	Context context.Context
}

func (r *TestGradesRepo) Aggregate(ctx context.Context, query shared.GradeQuery) (shared.GradeResult, error) {
	r.Calls++
	r.Context, r.Query = ctx, query
	return r.Result, r.Err
}

type TestTrendsRepo struct {
	SectionsResult []schema.Section
	Err            error
	Calls          int
	Filter         bson.M
	Context        context.Context
}

func (r *TestTrendsRepo) Sections(ctx context.Context, filter bson.M) ([]schema.Section, error) {
	r.Calls++
	r.Context, r.Filter = ctx, filter
	return r.SectionsResult, r.Err
}

var _ shared.GradesRepo = (*TestGradesRepo)(nil)
var _ shared.TrendsRepo = (*TestTrendsRepo)(nil)

func TestCourseGrades(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	for _, tc := range []struct {
		name, id      string
		grades        [14]int
		err           error
		status, calls int
	}{
		{"success", id.Hex(), [14]int{3, 4}, nil, 200, 1},
		{"empty grades", id.Hex(), [14]int{}, nil, 200, 1},
		{"invalid ID", "invalid", [14]int{}, nil, 400, 0},
		{"database error", id.Hex(), [14]int{}, errors.New("query failed"), 500, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &TestCourseRepo{GradeResults: tc.grades, Err: tc.err}
			var expected any
			if tc.status == 200 {
				expected = tc.grades
			}
			checkResponse(t, courseRequest(repo, "/course/"+tc.id+"/grades"), tc.status, expected)
			if repo.Calls != tc.calls {
				t.Fatalf("wrong call count: %d", repo.Calls)
			}
			if repo.Calls > 0 && (repo.Operation != "Grades" || repo.ID != id) {
				t.Fatal("incorrect grade lookup")
			}
		})
	}
}

func TestTrendsRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, target := range []string{"course", "professor", "combined"} {
		for _, scenario := range []string{"success", "empty", "database error", "canceled request"} {
			t.Run(target+"/"+scenario, func(t *testing.T) {
				sections := []schema.Section{{Id: bson.NewObjectID(), Section_number: "001"}}
				if scenario == "empty" {
					sections = []schema.Section{}
				}
				var repoErr error
				status := 200
				if scenario == "database error" {
					repoErr = errors.New("aggregation failed")
					status = 500
				}
				course := &TestCourseRepo{SectionResults: sections, Err: repoErr}
				professor := &TestProfessorRepo{&testRecordRepo[schema.Professor]{SectionResults: sections, Err: repoErr}}
				combined := &TestTrendsRepo{SectionsResult: sections, Err: repoErr}
				router := gin.New()
				var path string
				var want bson.M
				switch target {
				case "course":
					routes.CourseRoute(router, course)
					path = "/course/sections/trends?subject_prefix=CS&course_number=3345"
					want = bson.M{"_id": "CS3345"}
				case "professor":
					routes.ProfessorRoute(router, professor)
					path = "/professor/sections/trends?first_name=Jane&last_name=Doe"
					want = bson.M{"first_name": "Jane", "last_name": "Doe"}
				case "combined":
					routes.CombinedRoute(router, combined)
					path = "/combined/sections/trends?subject_prefix=CS&course_number=3345&first_name=Jane&last_name=Doe"
					want = bson.M{"_id": bson.M{"course": "CS3345", "prof_first": "Jane", "prof_last": "Doe"}}
				}
				request := httptest.NewRequest(http.MethodGet, path, nil)
				if scenario == "canceled request" {
					ctx, cancel := context.WithCancel(request.Context())
					cancel()
					request = request.WithContext(ctx)
				}
				response := httptest.NewRecorder()
				router.ServeHTTP(response, request)
				var expected any
				if status == 200 {
					expected = sections
				}
				checkResponse(t, response, status, expected)
				calls, filter, ctx := course.Calls, course.Filter, course.Context
				if target == "professor" {
					calls, filter, ctx = professor.Calls, professor.Filter, professor.Context
				}
				if target == "combined" {
					calls, filter, ctx = combined.Calls, combined.Filter, combined.Context
				}
				if calls != 1 || !reflect.DeepEqual(filter, want) {
					t.Fatalf("incorrect repository call: %d, %#v", calls, filter)
				}
				if _, ok := ctx.Deadline(); !ok {
					t.Fatal("missing deadline")
				}
				if ctx.Err() != context.Canceled {
					t.Fatal("context not canceled")
				}
			})
		}
	}
	t.Run("invalid professor filter", func(t *testing.T) {
		repo := &TestProfessorRepo{&testRecordRepo[schema.Professor]{}}
		router := gin.New()
		routes.ProfessorRoute(router, repo)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/professor/sections/trends?unknown=value", nil))
		checkResponse(t, response, 400, nil)
		if repo.Calls != 0 {
			t.Fatal("invalid filter reached repository")
		}
	})
}

func TestGradeAggregationRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	result := shared.GradeResult{Overall: [14]int{3, 4}, Semesters: []schema.GradeData{{Id: "26f", GradeDistribution: [14]int{3, 4}}}, SectionTypes: []schema.TypedGradeData{{Id: "26f"}}}
	for _, mode := range []struct {
		path, name string
		expected   any
	}{
		{"overall", "overall", result.Overall},
		{"semester", "semester", result.Semesters},
		{"semester/sectionType", "section_type", result.SectionTypes},
	} {
		for _, scenario := range []struct {
			name   string
			err    error
			status int
		}{
			{"success", nil, 200}, {"invalid query", shared.ErrInvalidGradeQuery, 400}, {"database error", errors.New("query failed"), 500},
		} {
			t.Run(mode.name+"/"+scenario.name, func(t *testing.T) {
				repo := &TestGradesRepo{Result: result, Err: scenario.err}
				router := gin.New()
				routes.GradesRoute(router, repo)
				request := httptest.NewRequest(http.MethodGet, "/grades/"+mode.path+"?prefix=CS&number=3345&section_number=001&first_name=Jane&last_name=Doe", nil)
				response := httptest.NewRecorder()
				router.ServeHTTP(response, request)
				var expected any
				if scenario.status == 200 {
					expected = mode.expected
				}
				checkResponse(t, response, scenario.status, expected)
				want := shared.GradeQuery{Mode: mode.name, Prefix: "CS", Number: "3345", SectionNumber: "001", FirstName: "Jane", LastName: "Doe"}
				if repo.Calls != 1 || !reflect.DeepEqual(repo.Query, want) {
					t.Fatalf("wrong query: %+v", repo.Query)
				}
				if _, ok := repo.Context.Deadline(); !ok {
					t.Fatal("missing deadline")
				}
				if repo.Context.Err() != context.Canceled {
					t.Fatal("context not canceled")
				}
			})
		}
	}
}
