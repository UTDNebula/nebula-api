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

// TestGetAggregateLimit now achieves regression testing
func TestGetAggregateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("DefaultValuesWhenNoParams", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)

		query := bson.M{}
		paginateMap, err := GetAggregateLimit(&query, c)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify default initialization (The core of your refactor)
		for _, key := range []string{"former_offset", "latter_offset"} {
			stage := paginateMap[key]
			if len(stage) == 0 || stage[0].Key != "$skip" || stage[0].Value != int64(0) {
				t.Errorf("Expected default $skip 0 for %s, got %v", key, stage)
			}
		}
	})

	t.Run("ValidOverrideFromQuery", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Provide only one offset to see if the other remains default
		c.Request = httptest.NewRequest("GET", "/?former_offset=50", nil)

		query := bson.M{"former_offset": "to-be-deleted"}
		paginateMap, _ := GetAggregateLimit(&query, c)

		// Check override
		former := paginateMap["former_offset"]
		if former[0].Value != int64(50) {
			t.Errorf("Expected former_offset to be 50, got %v", former[0].Value)
		}

		// Verify the query map was cleaned
		if _, exists := query["former_offset"]; exists {
			t.Error("Expected 'former_offset' to be deleted from query map")
		}
	})

	t.Run("InvalidParamReturnsErrorAndDefaults", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?former_offset=invalid", nil)

		query := bson.M{}
		_, err := GetAggregateLimit(&query, c)

		if err == nil {
			t.Error("Expected error for non-integer offset")
		}
	})
}

// the new version:  paginateMap[former_offset] is no longer an int so we are testing for Key($skip) and the Value
// we are also verifying delete(*query, field) logic 
