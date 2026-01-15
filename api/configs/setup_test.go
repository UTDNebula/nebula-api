package configs

import (
	"net/http/httptest"
	"strconv" // Added missing import
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// TestGetOptionLimit checks if the function correctly parses offset from query params
func TestGetOptionLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidOffset", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?offset=25", nil)

		query := bson.M{"offset": "should-be-deleted"}
		options, err := GetOptionLimit(&query, c)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err) // Use Fatalf to stop if options is nil
		}
		
		if _, exists := query["offset"]; exists {
			t.Error("Expected 'offset' to be deleted from the query map")
		}
		
		// Ensure we compare the same types (int64)
		if options.Skip == nil || *options.Skip != int64(25) {
			t.Errorf("Expected Skip to be 25, got %v", options.Skip)
		}
	})

	t.Run("EmptyOffset", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)

		query := bson.M{}
		options, _ := GetOptionLimit(&query, c)

		if options.Skip == nil || *options.Skip != int64(0) {
			t.Errorf("Expected default Skip to be 0, got %v", options.Skip)
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

// TestGetAggregateLimit achieves regression testing for the refactored logic
func TestGetAggregateLimit(t *testing.T) {
	gin.SetMode(gin.SetMode)

	t.Run("DefaultValuesWhenNoParams", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)

		query := bson.M{}
		paginateMap, err := GetAggregateLimit(&query, c)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check for the presence of keys and their BSON structure
		for _, key := range []string{"former_offset", "latter_offset"} {
			stage, ok := paginateMap[key]
			if !ok {
				t.Fatalf("Key %s missing from paginateMap", key)
			}
			// Use type assertion or direct indexing for bson.D
			if len(stage) == 0 || stage[0].Key != "$skip" || !isEqual(stage[0].Value, 0) {
				t.Errorf("Expected default $skip 0 for %s, got %v", key, stage)
			}
		}
	})

	t.Run("ValidOverrideFromQuery", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?former_offset=50", nil)

		query := bson.M{"former_offset": "to-be-deleted"}
		paginateMap, _ := GetAggregateLimit(&query, c)

		former := paginateMap["former_offset"]
		// Explicitly check the value as int64 to avoid architecture-based build fails
		if !isEqual(former[0].Value, 50) {
			t.Errorf("Expected former_offset to be 50, got %v", former[0].Value)
		}

		if _, exists := query["former_offset"]; exists {
			t.Error("Expected 'former_offset' to be deleted from query map")
		}
	})
}

// Helper function to handle cross-platform integer comparisons safely
func isEqual(actual interface{}, expected int) bool {
	switch v := actual.(type) {
	case int:
		return v == expected
	case int64:
		return v == int64(expected)
	case int32:
		return v == int32(expected)
	default:
		return false
	}
}
