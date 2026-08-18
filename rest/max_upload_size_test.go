package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func TestMaxUploadSize(t *testing.T) {
	// Set the environment variable for max upload size (e.g., 100 bytes)
	os.Setenv("MAX_UPLOAD_SIZE", "100")
	defer os.Unsetenv("MAX_UPLOAD_SIZE")

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/storage/:bucket/:objectID", func(c *gin.Context) {

		c.Set("gcsClient", &storage.Client{})
		controllers.PostObject(c)
	})

	t.Run("Upload within limit", func(t *testing.T) {

		defer func() {
			if r := recover(); r != nil {

			}
		}()

		data := make([]byte, 50)
		req, _ := http.NewRequest("POST", "/storage/test-bucket/small-file", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/octet-stream")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status NOT 413, got %d", w.Code)
		}
	})

	t.Run("Upload exceeding limit via Content-Length", func(t *testing.T) {
		data := make([]byte, 150)
		req, _ := http.NewRequest("POST", "/storage/test-bucket/large-file", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/octet-stream")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status 413, got %d. Body: %s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "File too large") {
			t.Errorf("Expected error message 'File too large', got %s", w.Body.String())
		}
	})

	t.Run("Upload exceeding limit via Stream", func(t *testing.T) {
		pr, pw := io.Pipe()
		go func() {
			pw.Write(make([]byte, 150))
			pw.Close()
		}()

		req, _ := http.NewRequest("POST", "/storage/test-bucket/stream-file", pr)
		req.ContentLength = -1
		req.Header.Set("Content-Type", "application/octet-stream")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status 413, got %d. Body: %s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "File too large") {
			t.Errorf("Expected error message 'File too large', got %s", w.Body.String())
		}
	})
}
