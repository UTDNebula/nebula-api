package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UTDNebula/nebula-api/rest/routes"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSectionCourseById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	course := schema.Course{Id: bson.NewObjectID(), Title: "Example course"}
	for _, tc := range []struct {
		name, id      string
		courses       []schema.Course
		err           error
		status, calls int
		expected      any
	}{
		{"success", id.Hex(), []schema.Course{course}, nil, 200, 1, course},
		{"missing section", id.Hex(), []schema.Course{}, nil, 404, 1, "No section with given ID"},
		{"invalid ID", "invalid", nil, nil, 400, 0, nil},
		{"database error", id.Hex(), nil, errors.New("query failed"), 500, 1, "query failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &TestSectionRepo{&testRecordRepo[schema.Section]{CourseResults: tc.courses, Err: tc.err}}
			router := gin.New()
			routes.SectionRoute(router, repo)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/section/"+tc.id+"/course", nil))
			checkResponse(t, response, tc.status, tc.expected)
			if repo.Calls != tc.calls {
				t.Fatalf("expected %d calls, got %d", tc.calls, repo.Calls)
			}
		})
	}
}
