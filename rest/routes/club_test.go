package routes

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestClubRouteSkipsRoutesWithoutDatabaseURI(t *testing.T) {
	for _, test := range []struct {
		name  string
		unset bool
	}{
		{name: "unset", unset: true},
		{name: "empty"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("CLUBS_DB_URI", "")
			if test.unset {
				if err := os.Unsetenv("CLUBS_DB_URI"); err != nil {
					t.Fatal(err)
				}
			}

			router := gin.New()
			router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
			ClubRoute(router)

			for _, route := range router.Routes() {
				if strings.HasPrefix(route.Path, "/club") {
					t.Errorf("unexpected club route %s", route.Path)
				}
			}

			for _, request := range []struct {
				path string
				want int
			}{
				{path: "/health", want: http.StatusOK},
				{path: "/club/search", want: http.StatusNotFound},
				{path: "/club/123", want: http.StatusNotFound},
			} {
				response := httptest.NewRecorder()
				router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, request.path, nil))
				if response.Code != request.want {
					t.Errorf("GET %s returned %d, want %d", request.path, response.Code, request.want)
				}
			}
		})
	}
}
