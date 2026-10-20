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

func TestCourseById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                    string
		id                      string
		expectedStatus          int
		expectedResponseMessage string
		expectedData            string
	}{
		{
			name:                    "ReturnCourseSuccessfully",
			id:                      "6a8f623120b48b6efe790172",
			expectedStatus:          http.StatusOK,
			expectedResponseMessage: "success",
		},
		{
			name:                    "InvalidMongoId",
			id:                      "This is invalid MongoDB ID",
			expectedStatus:          http.StatusBadRequest,
			expectedResponseMessage: `Parameter "id" is not a valid ObjectID.`,
		},
		{
			name:                    "NoCourseFound",
			id:                      primitive.NewObjectID().Hex(),
			expectedStatus:          http.StatusNotFound,
			expectedResponseMessage: "error",
			expectedData:            "No courses with given ID",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(responseRecorder)
			context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			context.Params = []gin.Param{{Key: "id", Value: testCase.id}}

			CourseById(context)

			if responseRecorder.Code != testCase.expectedStatus {
				t.Fatalf("expected status %d, got %d", testCase.expectedStatus, responseRecorder.Code)
			}

			var response schema.APIResponse[json.RawMessage]
			if err := json.Unmarshal(responseRecorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if response.Status != testCase.expectedStatus {
				t.Errorf("expected response status %d, got %d", testCase.expectedStatus, response.Status)
			}
			if response.Message != testCase.expectedResponseMessage {
				t.Errorf("expected response message %q, got %q", testCase.expectedResponseMessage, response.Message)
			}

			if testCase.expectedData == "" {
				return
			}

			var responseData string
			if err := json.Unmarshal(response.Data, &responseData); err != nil {
				t.Fatalf("decode response data: %v", err)
			}
			if responseData != testCase.expectedData {
				t.Errorf("expected response data %q, got %q", testCase.expectedData, responseData)
			}
		})
	}
}
