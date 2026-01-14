package configs

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// TestGetOptionLimit checks if the function correctly parses offset from query params
func TestGetOptionLimit(t *testing.T) {
	// Set Gin to test mode to keep logs clean
	gin.SetMode(gin.TestMode)

	t.Run("ValidOffset", func(t *testing.T) {
		// Create a mock context with a query parameter --> offset=25
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?offset=25", nil)

		query := bson.M{"offset": "should-be-deleted"}
		options, err := GetOptionLimit(&query, c)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		// Verify offset was deleted from the query map
		if _, exists := query["offset"]; exists {
			t.Error("Expected 'offset' to be deleted from the query map")
		}
		// Verify Skip was set correctly in Mongo options
		if *options.Skip != 25 {
			t.Errorf("Expected Skip to be 25, got %d", *options.Skip)
		}
	})

	t.Run("EmptyOffset", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil) // No offset param

		query := bson.M{}
		options, _ := GetOptionLimit(&query, c)

		// Default offset = 0
		if *options.Skip != 0 {
			t.Errorf("Expected default Skip to be 0, got %d", *options.Skip)
		}
	})

	t.Run("InvalidOffsetType", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?offset=not-a-number", nil)

		query := bson.M{}
		_, err := GetOptionLimit(&query, c)

		if err == nil {
			t.Error("Expected an error when parsing a non-integer offset")
		}
	})
}

// TestGetAggregateLimit checks if multi-offset pagination works for aggregation pipelines
func TestGetAggregateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("MultipleOffsets", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?former_offset=10&latter_offset=5", nil)

		query := bson.M{"former_offset": 10, "latter_offset": 5}
		paginateMap, err := GetAggregateLimit(&query, c)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		// Check if former_offset was parsed correctly
		if paginateMap["former_offset"] != 10 {
			t.Errorf("Expected former_offset 10, got %d", paginateMap["former_offset"])
		}
		// Check if latter_offset was parsed correctly
		if paginateMap["latter_offset"] != 5 {
			t.Errorf("Expected latter_offset 5, got %d", paginateMap["latter_offset"])
		}
	})
}
