package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSectionCourseById_SectionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	context.Params = []gin.Param{{Key: "id", Value: primitive.NewObjectID().Hex()}}

	SectionCourseById(context)

	if responseRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, responseRecorder.Code)
	}

	var response schema.APIResponse[string]
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data != "No section with given ID" {
		t.Errorf("expected not-found message, got %q", response.Data)
	}
}
