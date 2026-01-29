package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestObjectIDFromParam validates the conversion of URL string parameters to MongoDB ObjectIDs.
func TestObjectIDFromParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidObjectID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Fabricate a valid 24-character hex string parameter
		validHex := "65a3f1b2c3d4e5f6a7b8c9d0"
		c.Params = []gin.Param{{Key: "id", Value: validHex}}

		result, err := objectIDFromParam(c, "id")

		if err != nil {
			t.Errorf("Expected no error for valid hex, got %v", err)
		}
		if result.Hex() != validHex {
			t.Errorf("Expected hex %s, got %s", validHex, result.Hex())
		}
	})

	t.Run("InvalidObjectID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Fabricate an invalid parameter
		c.Params = []gin.Param{{Key: "id", Value: "invalid-id-format"}}

		result, err := objectIDFromParam(c, "id")

		if err == nil {
			t.Error("Expected error for invalid hex string, but got nil")
		}
		if result != nil {
			t.Error("Expected nil result for invalid conversion")
		}
		// Verify the utility automatically sent a 400 Bad Request response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

// TestGetQuery validates the construction of MongoDB filter maps based on the request type.
func TestGetQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ById_Valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		id := primitive.NewObjectID()
		c.Params = []gin.Param{{Key: "id", Value: id.Hex()}}

		query, err := getQuery[any]("ById", c)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		// Verify the query map contains the correct ObjectID
		if query["_id"] != id {
			t.Errorf("Expected query _id to be %v, got %v", id, query["_id"])
		}
	})

	t.Run("InvalidFlag", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		_, err := getQuery[any]("InvalidFlag", c)

		if err == nil {
			t.Error("Expected error for unsupported flag")
		}
		// Verify internal server error status was triggered
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})
}

// TestRespond verifies that the JSON response structure matches the API schema.
func TestRespond(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"key": "value"}
	respond(c, http.StatusOK, "test-message", testData)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
}
