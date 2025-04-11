package controllers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"

	"github.com/UTDNebula/nebula-api/api/responses"
)

// Get client from routes
func getClient(c *gin.Context) *storage.Client {
	val, exists := c.Get("gcsClient")
	if !exists {
		panic("storage client not set in context")
	}
	return val.(*storage.Client)
}

// @Id bucketInfo
// @Router /storage/{bucket} [get]
// @Description "Get info on a bucket"
// @Param bucket path string true "Name of the bucket"
// @Success 200
func BucketInfo(c *gin.Context) {
	bucket := c.Param("bucket")
	client := getClient(c)
	ctx := context.Background()

	// Get attributes
	attrs, err := client.Bucket(bucket).Attrs(ctx)
	// Create bucket if it does not exist
	if errors.Is(err, storage.ErrBucketNotExist) {
		err = client.Bucket(bucket).Create(ctx, "nebula-api-368223", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to create bucket: " + err.Error()})
			return
		}
		attrs, err = client.Bucket(bucket).Attrs(ctx)
	}
	// Catch all from above
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to get bucket attributes: " + err.Error()})
		return
	}

	// Loop through objects and add names
	contents := []string{}
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		contents = append(contents, objAttrs.Name)
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": gin.H{"attrs": attrs, "contents": contents}})
}

// @Id deleteBucket
// @Router /storage/{bucket} [delete]
// @Description "Delete a bucket"
// @Param bucket path string true "Name of the bucket"
// @Success 200
func DeleteBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	client := getClient(c)
	ctx := context.Background()

	// Delete all objects (GCS requires an empty bucket before deletion)
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err := client.Bucket(bucket).Object(objAttrs.Name).Delete(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to delete object: " + err.Error()})
			return
		}
	}

	// Delete bucket
	if err := client.Bucket(bucket).Delete(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to delete bucket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// @Id objectInfo
// @Router /storage/{bucket}/{objectID} [get]
// @Description "Get info on an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func ObjectInfo(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	// Get object attreibutes
	attrs, err := client.Bucket(bucket).Object(objectID).Attrs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": attrs})
}

// @Id postObject
// @Router /storage/{bucket}/{objectID} [post]
// @Description "Upload an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func PostObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	// Read body as byte stream
	fileReader := c.Request.Body
	if fileReader == nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: "Empty body"})
		return
	}
	defer fileReader.Close()

	// Set metadata
	wc := client.Bucket(bucket).Object(objectID).NewWriter(ctx)
	wc.ContentType = c.ContentType()
	wc.CacheControl = "public, max-age=3600"
	wc.Metadata = map[string]string{
		"uploaded-at": time.Now().Format(time.RFC3339),
	}

	// Upload
	if _, err := io.Copy(wc, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// @Id deleteObject
// @Router /storage/{bucket}/{objectID} [delete]
// @Description "Delete an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	err := client.Bucket(bucket).Object(objectID).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
