package controllers

import (
	"github.com/gin-gonic/gin"
)

// @Id bucketInfo
// @Router /{bucket} [get]
// @Description "Get info on a bucket"
// @Param bucket path string true "Name of the bucket"
// @Success 200
func BucketInfo(c *gin.Context) {

}

// @Id deleteBucket
// @Router /{bucket} [get]
// @Description "Delete a bucket"
// @Param bucket path string true "Name of the bucket"
// @Success 200
func DeleteBucket(c *gin.Context) {

}

// @Id objectInfo
// @Router /storage/{bucket}/info/{objectID} [get]
// @Description "Get info on an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func ObjectInfo(c *gin.Context) {

}

// @Id postObject
// @Router /{bucket} [get]
// @Router /storage/{bucket}/info/{objectID} [post]
// @Description "Get info on an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func PostObject(c *gin.Context) {

}

// @Id deleteObject
// @Router /storage/{bucket}/info/{objectID} [delete]
// @Description "Get info on an object"
// @Param bucket path string true "Name of the bucket"
// @Param objectID path string true "ID of the object"
// @Success 200
func DeleteObject(c *gin.Context) {

}
