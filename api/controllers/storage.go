package controllers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"

	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"
)

const (
	PROJECT_ID = "nebula-api-368223"
)

// Get client from routes
func getClient(c *gin.Context) *storage.Client {
	val, exists := c.Get("gcsClient")
	if !exists {
		panic("storage client not set in context")
	}
	return val.(*storage.Client)
}

// Get bucket or create it if it doesn't already exist
func getOrCreateBucket(client *storage.Client, bucket string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	// Get bucket, or create it if it does not exist
	// NOTE: We automatically prefix bucket names with "utdnebula_" here since bucket names need to be GLOBALLY unique
	bucketHandle := client.Bucket(schema.BUCKET_PREFIX + bucket)
	_, err := bucketHandle.Attrs(ctx)
	if err != nil {
		err = bucketHandle.Create(ctx, PROJECT_ID, nil)
		if err != nil {
			return nil, errors.New("failed to create bucket: " + err.Error())
		}
	}
	return bucketHandle, nil
}

// @Id				bucketInfo
// @Router			/storage/{bucket} [get]
// @Description	"Get info on a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket	path		string				true	"Name of the bucket"
// @Success		200		{object}	schema.BucketInfo	"The bucket's info"
// @security		storage_key
func BucketInfo(c *gin.Context) {
	bucket := c.Param("bucket")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// Get attributes
	attrs, err := bucketHandle.Attrs(ctx)

	// Catch all from above
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to get bucket attributes: " + err.Error()})
		return
	}

	// Loop through objects and add names
	contents := []string{}
	it := bucketHandle.Objects(ctx, nil)
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

	bucketInfo := schema.BucketInfoFromAttrs(attrs)
	bucketInfo.Contents = contents

	c.JSON(http.StatusOK, responses.BucketResponse{Status: http.StatusOK, Message: "success", Data: bucketInfo})
}

// @Id				deleteBucket
// @Router			/storage/{bucket} [delete]
// @Description	"Delete a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket	path	string	true	"Name of the bucket"
// @Success		200
// @security		storage_key
func DeleteBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// Delete all objects (GCS requires an empty bucket before deletion)
	it := bucketHandle.Objects(ctx, nil)
	deletedCount := 0
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err := bucketHandle.Object(objAttrs.Name).Delete(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to delete object: " + err.Error()})
			return
		}
		deletedCount++
	}

	// Delete bucket
	if err := bucketHandle.Delete(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "Failed to delete bucket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.DeleteResponse{Status: http.StatusOK, Message: "success", Data: deletedCount})
}

// @Id				objectInfo
// @Router			/storage/{bucket}/{objectID} [get]
// @Description	"Get info on an object in a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket		path		string				true	"Name of the bucket"
// @Param			objectID	path		string				true	"ID of the object"
// @Success		200			{object}	schema.ObjectInfo	"The object's info"
// @security		storage_key
func ObjectInfo(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	objectHandle := bucketHandle.Object(objectID)
	if objectHandle == nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "invalid object id"})
		return
	}

	// Get object attributes
	attrs, err := objectHandle.Attrs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	objectInfo := schema.ObjectInfoFromAttrs(attrs)
	c.JSON(http.StatusOK, responses.ObjectResponse{Status: http.StatusOK, Message: "success", Data: objectInfo})
}

// @Id				postObject
// @Router			/storage/{bucket}/{objectID} [post]
// @Description	"Upload an object to a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket		path		string				true	"Name of the bucket"
// @Param			objectID	path		string				true	"ID of the object"
// @Param			data		formData	file				true	"The data to upload"
// @Success		200			{object}	schema.ObjectInfo	"The object's info"
// @security		storage_key
func PostObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// Read body as byte stream
	fileReader := c.Request.Body
	if fileReader == nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: "Empty body"})
		return
	}
	defer fileReader.Close()

	objectHandle := bucketHandle.Object(objectID)
	if objectHandle == nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "invalid object id"})
		return
	}

	wc := objectHandle.NewWriter(ctx)
	// Makes object public
	wc.ACL = []storage.ACLRule{
		{Entity: storage.AllUsers, EntityID: "", Role: storage.RoleReader, Domain: "", Email: "", ProjectTeam: nil},
	}
	// Set metadata
	wc.CacheControl = "public, max-age=3600"

	// Upload
	if _, err := io.Copy(wc, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	attrs, err := objectHandle.Attrs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	objectInfo := schema.ObjectInfoFromAttrs(attrs)
	c.JSON(http.StatusOK, responses.ObjectResponse{Status: http.StatusOK, Message: "success", Data: objectInfo})
}

// @Id				deleteObject
// @Router			/storage/{bucket}/{objectID} [delete]
// @Description	"Delete an object from a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket		path	string	true	"Name of the bucket"
// @Param			objectID	path	string	true	"ID of the object"
// @Success		200
// @security		storage_key
func DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	objectHandle := bucketHandle.Object(objectID)
	if objectHandle == nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "invalid object id"})
		return
	}

	err = objectHandle.Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.DeleteResponse{Status: http.StatusOK, Message: "success", Data: 1})
}
