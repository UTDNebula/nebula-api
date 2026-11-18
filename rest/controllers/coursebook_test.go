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

// testRecordRepo records requests and returns controlled results for both entities.
type testRecordRepo[T any] struct {
	Record           T
	Records          []T
	CourseResults    []schema.Course
	SectionResults   []schema.Section
	ProfessorResults []schema.Professor
	GradeResults     [14]int
	Err              error
	Calls            int
	Operation        string
	ID               bson.ObjectID
	Filter           bson.M
	Offset           int64
	Limit            int64
	Pagination       map[string]bson.D
	Context          context.Context
}

func (r *testRecordRepo[T]) record(ctx context.Context, operation string) {
	r.Calls++
	r.Operation, r.Context = operation, ctx
}
func (r *testRecordRepo[T]) ByID(ctx context.Context, id bson.ObjectID) (T, error) {
	r.record(ctx, "ByID")
	r.ID = id
	return r.Record, r.Err
}
func (r *testRecordRepo[T]) Search(ctx context.Context, filter bson.M, offset, limit int64) ([]T, error) {
	r.record(ctx, "Search")
	r.Filter, r.Offset, r.Limit = filter, offset, limit
	return r.Records, r.Err
}
func (r *testRecordRepo[T]) All(ctx context.Context) ([]T, error) {
	r.record(ctx, "All")
	return r.Records, r.Err
}
func (r *testRecordRepo[T]) Courses(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Course, error) {
	r.record(ctx, "Courses")
	r.Filter, r.Pagination = filter, pagination
	return r.CourseResults, r.Err
}
func (r *testRecordRepo[T]) Sections(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Section, error) {
	r.record(ctx, "Sections")
	r.Filter, r.Pagination = filter, pagination
	return r.SectionResults, r.Err
}
func (r *testRecordRepo[T]) Professors(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Professor, error) {
	r.record(ctx, "Professors")
	r.Filter, r.Pagination = filter, pagination
	return r.ProfessorResults, r.Err
}
func (r *testRecordRepo[T]) Grades(ctx context.Context, id bson.ObjectID) ([14]int, error) {
	r.record(ctx, "Grades")
	r.ID = id
	return r.GradeResults, r.Err
}
func (r *testRecordRepo[T]) TrendsSections(ctx context.Context, filter bson.M) ([]schema.Section, error) {
	r.record(ctx, "TrendsSections")
	r.Filter = filter
	return r.SectionResults, r.Err
}

type TestSectionRepo struct {
	*testRecordRepo[schema.Section]
}
type TestProfessorRepo struct {
	*testRecordRepo[schema.Professor]
}

var _ shared.SectionRepo = (*TestSectionRepo)(nil)
var _ shared.ProfessorRepo = (*TestProfessorRepo)(nil)

type relationshipEndpoint struct {
	searchPath string
	idPath     string
	operation  string
	single     bool
}

type entityTest[T any] struct {
	name          string
	id            bson.ObjectID
	record        T
	filterKey     string
	filterValue   string
	all           bool
	relationships []relationshipEndpoint
	register      func(*gin.Engine, *testRecordRepo[T])
}

func TestSectionRoutes(t *testing.T) {
	id := bson.NewObjectID()
	runEntityTests(t, entityTest[schema.Section]{
		name: "section", id: id, record: schema.Section{Id: id, Section_number: "001"},
		filterKey: "section_number", filterValue: "001",
		relationships: []relationshipEndpoint{
			{searchPath: "courses", idPath: "course", operation: "Courses", single: true},
			{searchPath: "professors", idPath: "professors", operation: "Professors"},
		},
		register: func(router *gin.Engine, repo *testRecordRepo[schema.Section]) {
			routes.SectionRoute(router, &TestSectionRepo{repo})
		},
	})
}
func TestProfessorRoutes(t *testing.T) {
	id := bson.NewObjectID()
	runEntityTests(t, entityTest[schema.Professor]{
		name: "professor", id: id, record: schema.Professor{Id: id, First_name: "Jane"},
		filterKey: "first_name", filterValue: "Jane", all: true,
		relationships: []relationshipEndpoint{
			{searchPath: "courses", idPath: "courses", operation: "Courses"},
			{searchPath: "sections", idPath: "sections", operation: "Sections"},
		},
		register: func(router *gin.Engine, repo *testRecordRepo[schema.Professor]) {
			routes.ProfessorRoute(router, &TestProfessorRepo{repo})
		},
	})
}

func runEntityTests[T any](t *testing.T, entity entityTest[T]) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Setenv("LIMIT", "7")
	request := func(repo *testRecordRepo[T], path string) *httptest.ResponseRecorder {
		router := gin.New()
		entity.register(router, repo)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		return response
	}
	newRepo := func() *testRecordRepo[T] {
		return &testRecordRepo[T]{Record: entity.record, Records: []T{entity.record},
			CourseResults:    []schema.Course{{Title: "Example course"}},
			SectionResults:   []schema.Section{{Section_number: "001"}},
			ProfessorResults: []schema.Professor{{First_name: "Jane"}},
			GradeResults:     [14]int{1, 2, 3},
		}
	}
	t.Run("by ID", func(t *testing.T) {
		cases := []struct {
			name, id      string
			err           error
			status, calls int
			data          any
		}{
			{"success", entity.id.Hex(), nil, 200, 1, entity.record},
			{"invalid ID", "invalid", nil, 400, 0, nil},
			{"missing record", entity.id.Hex(), shared.ErrNotFound, 404, 1, "No " + entity.name + "s with given ID"},
			{"wrapped missing record", entity.id.Hex(), errors.Join(errors.New("lookup"), shared.ErrNotFound), 404, 1, "No " + entity.name + "s with given ID"},
			{"database error", entity.id.Hex(), errors.New("query failed"), 500, 1, "query failed"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				repo := newRepo()
				repo.Err = tc.err
				checkResponse(t, request(repo, "/"+entity.name+"/"+tc.id), tc.status, tc.data)
				if repo.Calls != tc.calls {
					t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
				}
				if repo.Calls > 0 && (repo.Operation != "ByID" || repo.ID != entity.id) {
					t.Fatal("incorrect repository operation or ID")
				}
				if repo.Calls > 0 && repo.Context.Err() != context.Canceled {
					t.Fatal("timeout context was not canceled")
				}
			})
		}
	})
	t.Run("search", func(t *testing.T) {
		for _, tc := range []struct {
			name, query   string
			err           error
			status, calls int
			offset        int64
		}{
			{"success", "?" + entity.filterKey + "=" + entity.filterValue + "&offset=3", nil, 200, 1, 3},
			{"default offset", "?" + entity.filterKey + "=" + entity.filterValue, nil, 200, 1, 0},
			{"invalid offset", "?offset=invalid", nil, 400, 0, 0},
			{"unknown filter", "?unknown=value", nil, 400, 0, 0},
			{"database error", "?" + entity.filterKey + "=" + entity.filterValue, errors.New("query failed"), 500, 1, 0},
		} {
			t.Run(tc.name, func(t *testing.T) {
				repo := newRepo()
				repo.Err = tc.err
				var expected any
				if tc.status == 200 {
					expected = repo.Records
				}
				checkResponse(t, request(repo, "/"+entity.name+tc.query), tc.status, expected)
				if repo.Calls != tc.calls {
					t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
				}
				if repo.Calls > 0 && (repo.Operation != "Search" || repo.Offset != tc.offset || repo.Limit != 7 || !reflect.DeepEqual(repo.Filter, bson.M{entity.filterKey: entity.filterValue})) {
					t.Fatalf("incorrect search inputs: %+v", repo)
				}
			})
		}
	})
	if entity.all {
		t.Run("all", func(t *testing.T) {
			for _, fail := range []bool{false, true} {
				t.Run(map[bool]string{false: "success", true: "database error"}[fail], func(t *testing.T) {
					repo := newRepo()
					status := 200
					var expected any = repo.Records
					if fail {
						repo.Err = errors.New("query failed")
						status = 500
						expected = nil
					}
					checkResponse(t, request(repo, "/"+entity.name+"/all"), status, expected)
					if repo.Operation != "All" || repo.Calls != 1 {
						t.Fatal("incorrect repository operation")
					}
				})
			}
		})
	}
	for _, endpoint := range entity.relationships {
		for _, byID := range []bool{false, true} {
			path := "/" + entity.name + "/" + endpoint.searchPath
			if byID {
				path = "/" + entity.name + "/" + entity.id.Hex() + "/" + endpoint.idPath
			}
			type testCase struct {
				name, path    string
				err           error
				status, calls int
				empty         bool
			}
			cases := []testCase{
				{"success", path + "?former_offset=2&latter_offset=4", nil, 200, 1, false},
				{"empty result", path, nil, 200, 1, true},
				{"database error", path, errors.New("aggregation failed"), 500, 1, false},
				{"invalid former offset", path + "?former_offset=invalid", nil, 400, 0, false},
				{"invalid latter offset", path + "?latter_offset=invalid", nil, 400, 0, false},
			}
			if byID {
				cases = append(cases, testCase{"invalid ID", "/" + entity.name + "/invalid/" + endpoint.idPath, nil, 400, 0, false})
			} else {
				cases = append(cases, testCase{"unknown filter", path + "?unknown=value", nil, 400, 0, false})
			}
			for _, tc := range cases {
				t.Run(endpoint.operation+"/"+map[bool]string{true: "by ID", false: "search"}[byID]+"/"+tc.name, func(t *testing.T) {
					repo := newRepo()
					repo.Err = tc.err
					if tc.empty {
						repo.CourseResults = []schema.Course{}
						repo.SectionResults = []schema.Section{}
						repo.ProfessorResults = []schema.Professor{}
					}
					var expected any
					status := tc.status
					switch endpoint.operation {
					case "Courses":
						expected = repo.CourseResults
					case "Sections":
						expected = repo.SectionResults
					case "Professors":
						expected = repo.ProfessorResults
					}
					if endpoint.single && byID {
						if tc.empty {
							status = 404
							expected = "No section with given ID"
						} else {
							expected = repo.CourseResults[0]
						}
					}
					if status != 200 && status != 404 {
						expected = nil
					}
					checkResponse(t, request(repo, tc.path), status, expected)
					if repo.Calls != tc.calls {
						t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
					}
					if repo.Calls > 0 {
						if repo.Operation != endpoint.operation {
							t.Fatal("incorrect relationship method")
						}
						if byID {
							id, ok := repo.Filter["_id"].(*bson.ObjectID)
							if !ok || *id != entity.id {
								t.Fatal("incorrect ID filter")
							}
						}
						if tc.name == "success" {
							expectedPagination := map[string]bson.D{"former_offset": {{Key: "$skip", Value: int64(2)}}, "latter_offset": {{Key: "$skip", Value: int64(4)}}, "limit": {{Key: "$limit", Value: int64(7)}}}
							if !reflect.DeepEqual(repo.Pagination, expectedPagination) {
								t.Fatalf("incorrect pagination: %#v", repo.Pagination)
							}
						}
					}
				})
			}
		}
	}
	t.Run("grades", func(t *testing.T) {
		for _, tc := range []struct {
			name, id      string
			err           error
			status, calls int
		}{
			{"success", entity.id.Hex(), nil, 200, 1},
			{"invalid ID", "invalid", nil, 400, 0},
			{"database error", entity.id.Hex(), errors.New("aggregation failed"), 500, 1},
		} {
			t.Run(tc.name, func(t *testing.T) {
				repo := newRepo()
				repo.Err = tc.err
				var expected any
				if tc.status == 200 {
					expected = repo.GradeResults
				}
				checkResponse(t, request(repo, "/"+entity.name+"/"+tc.id+"/grades"), tc.status, expected)
				if repo.Calls != tc.calls {
					t.Fatal("incorrect repository call count")
				}
				if repo.Calls > 0 && (repo.Operation != "Grades" || repo.ID != entity.id) {
					t.Fatal("incorrect grades lookup")
				}
			})
		}
	})
	t.Run("request cancellation", func(t *testing.T) {
		repo := newRepo()
		router := gin.New()
		entity.register(router, repo)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		req := httptest.NewRequest(http.MethodGet, "/"+entity.name+"/"+entity.id.Hex(), nil).WithContext(ctx)
		router.ServeHTTP(httptest.NewRecorder(), req)
		if repo.Context.Err() != context.Canceled {
			t.Fatal("request cancellation was not propagated")
		}
	})
}
