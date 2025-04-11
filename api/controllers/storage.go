package controllers

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
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

}

// @Id objectInfo
// @Router /storage/{bucket}/{objectID} [get]
// @Description "Get info on an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func ObjectInfo(c *gin.Context) {

}

// @Id postObject
// @Router /storage/{bucket}/{objectID} [post]
// @Description "Get info on an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func PostObject(c *gin.Context) {

}

// @Id deleteObject
// @Router /storage/{bucket}/{objectID} [delete]
// @Description "Delete an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func DeleteObject(c *gin.Context) {

}
