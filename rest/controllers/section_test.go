package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSectionCourseById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	// Select a section with an existing course without changing database records.
	cursor, err := sectionCollection.Aggregate(ctx, bson.A{
		bson.M{"$lookup": bson.M{"from": "courses", "localField": "course_reference", "foreignField": "_id", "as": "course"}},
		bson.M{"$unwind": "$course"},
		bson.M{"$limit": 1},
	})
	if err != nil {
		t.Fatalf("select a section with a course for the success case: %v", err)
	}
	defer cursor.Close(ctx)
	if !cursor.Next(ctx) {
		t.Fatalf("no section with a course available for the success case (cursor error: %v)", cursor.Err())
	}
	var section struct {
		ID     bson.ObjectID `bson:"_id"`
		Course schema.Course `bson:"course"`
	}
	if err := cursor.Decode(&section); err != nil {
		t.Fatalf("decode the section and course for the success case: %v", err)
	}

	testCases := []struct {
		name                    string
		id                      string
		expectedStatus          int
		expectedResponseMessage string
		expectedData            string
	}{
		{
			name:                    "ReturnCourseSuccessfully",
			id:                      section.ID.Hex(),
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
			name:                    "NoSectionFound",
			id:                      bson.NewObjectID().Hex(),
			expectedStatus:          http.StatusNotFound,
			expectedResponseMessage: "error",
			expectedData:            "No section with given ID",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(responseRecorder)
			context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			context.Params = []gin.Param{{Key: "id", Value: testCase.id}}

			SectionCourseById(context)

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

			if testCase.expectedStatus == http.StatusOK {
				var returnedCourse schema.Course
				if err := json.Unmarshal(response.Data, &returnedCourse); err != nil {
					t.Fatalf("decode returned course: %v", err)
				}
				if returnedCourse.Id != section.Course.Id {
					t.Errorf("returned course ID = %s, want %s", returnedCourse.Id.Hex(), section.Course.Id.Hex())
				}
				return
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
