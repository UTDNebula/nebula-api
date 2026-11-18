package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/UTDNebula/nebula-api/rest/routes"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestCourseRepo supplies deterministic results and records the controller's inputs.
type TestCourseRepo struct {
	GradeResults     [14]int
	Course           schema.Course
	Courses          []schema.Course
	SectionResults   []schema.Section
	ProfessorResults []schema.Professor
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

func (r *TestCourseRepo) record(ctx context.Context, operation string) {
	r.Calls++
	r.Context = ctx
	r.Operation = operation
}

func (r *TestCourseRepo) ByID(ctx context.Context, id bson.ObjectID) (schema.Course, error) {
	r.record(ctx, "ByID")
	r.ID = id
	return r.Course, r.Err
}

func (r *TestCourseRepo) Search(ctx context.Context, filter bson.M, offset, limit int64) ([]schema.Course, error) {
	r.record(ctx, "Search")
	r.Filter, r.Offset, r.Limit = filter, offset, limit
	return r.Courses, r.Err
}

func (r *TestCourseRepo) All(ctx context.Context) ([]schema.Course, error) {
	r.record(ctx, "All")
	return r.Courses, r.Err
}

func (r *TestCourseRepo) Sections(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Section, error) {
	r.record(ctx, "Sections")
	r.Filter, r.Pagination = filter, pagination
	return r.SectionResults, r.Err
}

func (r *TestCourseRepo) Professors(ctx context.Context, filter bson.M, pagination map[string]bson.D) ([]schema.Professor, error) {
	r.record(ctx, "Professors")
	r.Filter, r.Pagination = filter, pagination
	return r.ProfessorResults, r.Err
}

func (r *TestCourseRepo) Grades(ctx context.Context, id bson.ObjectID) ([14]int, error) {
	r.record(ctx, "Grades")
	r.ID = id
	return r.GradeResults, r.Err
}
func (r *TestCourseRepo) TrendsSections(ctx context.Context, filter bson.M) ([]schema.Section, error) {
	r.record(ctx, "TrendsSections")
	r.Filter = filter
	return r.SectionResults, r.Err
}

var _ shared.CourseRepo = (*TestCourseRepo)(nil)

func courseRequest(repo *TestCourseRepo, path string) *httptest.ResponseRecorder {
	router := gin.New()
	routes.CourseRoute(router, repo)
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func checkResponse(t *testing.T, response *httptest.ResponseRecorder, status int, expected any) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("expected HTTP %d, got %d: %s", status, response.Code, response.Body.String())
	}
	var body schema.APIResponse[json.RawMessage]
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != status {
		t.Fatalf("expected envelope status %d, got %d", status, body.Status)
	}
	if status == http.StatusOK && body.Message != "success" {
		t.Fatalf("unexpected message %q", body.Message)
	}
	if expected != nil {
		encoded, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}
		var want, got any
		if err := json.Unmarshal(encoded, &want); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(body.Data, &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("expected data %s, got %s", encoded, body.Data)
		}
	}
}

func TestCourseById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	course := schema.Course{Id: id, Title: "Example course", Subject_prefix: "CS", Course_number: "3345"}
	cases := []struct {
		name     string
		id       string
		err      error
		status   int
		calls    int
		expected any
	}{
		{"success", id.Hex(), nil, http.StatusOK, 1, course},
		{"invalid ID", "invalid", nil, http.StatusBadRequest, 0, nil},
		{"missing record", id.Hex(), shared.ErrNotFound, http.StatusNotFound, 1, "No courses with given ID"},
		{"wrapped missing record", id.Hex(), errors.Join(errors.New("lookup"), shared.ErrNotFound), http.StatusNotFound, 1, "No courses with given ID"},
		{"database error", id.Hex(), errors.New("query failed"), http.StatusInternalServerError, 1, "query failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &TestCourseRepo{Course: course, Err: tc.err}
			response := courseRequest(repo, "/course/"+tc.id)
			checkResponse(t, response, tc.status, tc.expected)
			if repo.Calls != tc.calls {
				t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
			}
			if repo.Calls > 0 {
				if repo.ID != id {
					t.Fatalf("wrong ID: %s", repo.ID.Hex())
				}
				if _, ok := repo.Context.Deadline(); !ok {
					t.Fatal("repository context has no deadline")
				}
				if repo.Context.Err() != context.Canceled {
					t.Fatal("handler did not cancel its timeout context")
				}
			}
		})
	}
}

func TestCourseSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("LIMIT", "7")
	courses := []schema.Course{{Title: "Example course", Subject_prefix: "CS"}}
	cases := []struct {
		name   string
		query  string
		err    error
		status int
		calls  int
		offset int64
	}{
		{"filtered search", "?subject_prefix=CS&offset=3", nil, 200, 1, 3},
		{"default offset", "?subject_prefix=CS", nil, 200, 1, 0},
		{"invalid offset", "?offset=invalid", nil, 400, 0, 0},
		{"unknown filter", "?unknown=value", nil, 400, 0, 0},
		{"database error", "?subject_prefix=CS", errors.New("query failed"), 500, 1, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &TestCourseRepo{Courses: courses, Err: tc.err}
			response := courseRequest(repo, "/course"+tc.query)
			var expected any
			if tc.status == http.StatusOK {
				expected = courses
			}
			checkResponse(t, response, tc.status, expected)
			if repo.Calls != tc.calls {
				t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
			}
			if repo.Calls > 0 {
				if repo.Operation != "Search" {
					t.Fatalf("wrong operation: %s", repo.Operation)
				}
				if !reflect.DeepEqual(repo.Filter, bson.M{"subject_prefix": "CS"}) {
					t.Fatalf("wrong filter: %#v", repo.Filter)
				}
				if repo.Offset != tc.offset || repo.Limit != 7 {
					t.Fatalf("wrong pagination: %d, %d", repo.Offset, repo.Limit)
				}
			}
		})
	}
}

func TestCourseAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name    string
		courses []schema.Course
		err     error
		status  int
	}{
		{"success", []schema.Course{{Title: "Example course"}}, nil, 200},
		{"empty result", []schema.Course{}, nil, 200},
		{"database error", nil, errors.New("query failed"), 500},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &TestCourseRepo{Courses: tc.courses, Err: tc.err}
			response := courseRequest(repo, "/course/all")
			var expected any
			if tc.status == 200 {
				expected = tc.courses
			}
			checkResponse(t, response, tc.status, expected)
			if repo.Operation != "All" || repo.Calls != 1 {
				t.Fatalf("wrong repository call: %+v", repo)
			}
		})
	}
}

func TestCourseRelationships(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("LIMIT", "7")
	id := bson.NewObjectID()
	for _, target := range []string{"sections", "professors"} {
		for _, byID := range []bool{false, true} {
			prefix := "/course/" + target
			query := "subject_prefix=CS&"
			if byID {
				prefix = "/course/" + id.Hex() + "/" + target
				query = ""
			}
			type testCase struct {
				name   string
				path   string
				err    error
				status int
				calls  int
				empty  bool
			}
			cases := []testCase{
				{"success", prefix + "?" + query + "former_offset=2&latter_offset=4", nil, 200, 1, false},
				{"empty result", prefix, nil, 200, 1, true},
				{"database error", prefix, errors.New("aggregation failed"), 500, 1, false},
				{"invalid former offset", prefix + "?former_offset=invalid", nil, 400, 0, false},
				{"invalid latter offset", prefix + "?latter_offset=invalid", nil, 400, 0, false},
			}
			if byID {
				cases = append(cases, testCase{"invalid ID", "/course/invalid/" + target, nil, 400, 0, false})
			} else {
				cases = append(cases, testCase{"unknown filter", prefix + "?unknown=value", nil, 400, 0, false})
			}
			for _, tc := range cases {
				t.Run(target+"/"+map[bool]string{true: "by ID", false: "search"}[byID]+"/"+tc.name, func(t *testing.T) {
					repo := &TestCourseRepo{SectionResults: []schema.Section{{}}, ProfessorResults: []schema.Professor{{}}, Err: tc.err}
					if tc.empty {
						repo.SectionResults = []schema.Section{}
						repo.ProfessorResults = []schema.Professor{}
					}
					response := courseRequest(repo, tc.path)
					var expected any
					if tc.status == 200 {
						if target == "sections" {
							expected = repo.SectionResults
						} else {
							expected = repo.ProfessorResults
						}
					}
					checkResponse(t, response, tc.status, expected)
					if repo.Calls != tc.calls {
						t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
					}
					if repo.Calls > 0 {
						if !strings.EqualFold(repo.Operation, target) {
							t.Fatalf("wrong operation: %s", repo.Operation)
						}
						if byID {
							got, ok := repo.Filter["_id"].(*bson.ObjectID)
							if !ok || *got != id {
								t.Fatalf("wrong ID filter: %#v", repo.Filter)
							}
						} else if tc.name == "success" && !reflect.DeepEqual(repo.Filter, bson.M{"subject_prefix": "CS"}) {
							t.Fatalf("wrong search filter: %#v", repo.Filter)
						}
						if tc.name == "success" {
							expectedPagination := map[string]bson.D{
								"former_offset": {{Key: "$skip", Value: int64(2)}},
								"latter_offset": {{Key: "$skip", Value: int64(4)}},
								"limit":         {{Key: "$limit", Value: int64(7)}},
							}
							if !reflect.DeepEqual(repo.Pagination, expectedPagination) {
								t.Fatalf("wrong pagination: %#v", repo.Pagination)
							}
						}
					}
				})
			}
		}
	}
}

func TestCourseRepositoryContextCancellation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	repo := &TestCourseRepo{}
	routes.CourseRoute(router, repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodGet, "/course/"+bson.NewObjectID().Hex(), nil).WithContext(ctx)
	router.ServeHTTP(httptest.NewRecorder(), request)
	if repo.Context.Err() != context.Canceled {
		t.Fatal("request cancellation was not propagated")
	}
}
