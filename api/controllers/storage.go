package controllers

import (
	"context"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

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
	attrs, err := client.Bucket(bucket).Attrs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attrs)
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

	// First delete all objects (GCS requires an empty bucket before deletion)
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := client.Bucket(bucket).Object(objAttrs.Name).Delete(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := client.Bucket(bucket).Delete(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bucket deleted"})
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
	attrs, err := client.Bucket(bucket).Object(objectID).Attrs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attrs)
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
	fileReader := c.Request.Body
	if fileReader == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty body"})
		return
	}
	defer fileReader.Close()

	wc := client.Bucket(bucket).Object(objectID).NewWriter(ctx)
	wc.ContentType = c.ContentType()
	wc.CacheControl = "public, max-age=3600"
	wc.Metadata = map[string]string{
		"uploaded-at": time.Now().Format(time.RFC3339),
	}

	if _, err := io.Copy(wc, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object uploaded", "object": objectID})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object deleted"})
}
