package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
)

func TestCourseById_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a response recorder to capture the result
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Test sending a completely invalid ID format
	c.Params = []gin.Param{{Key: "id", Value: "not-an-object-id"}}
	
	CourseById(c)
	
	// Verify that the controller utility correctly returns a 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}
}
