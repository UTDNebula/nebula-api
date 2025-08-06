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
// @Param			bucket			path		string									true	"Name of the bucket"
// @Param			x-storage-key	header		string									true	"The internal storage key"
// @Success		200				{object}	schema.APIResponse[schema.BucketInfo]	"The bucket's info"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
func BucketInfo(c *gin.Context) {
	bucket := c.Param("bucket")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Get attributes
	attrs, err := bucketHandle.Attrs(ctx)
	// Catch all from above
	if err != nil {
		respondWithInternalError(c, err)
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
			respondWithInternalError(c, err)
			return
		}
		contents = append(contents, objAttrs.Name)
	}

	bucketInfo := schema.BucketInfoFromAttrs(attrs)
	bucketInfo.Contents = contents

	respond(c, http.StatusOK, "success", bucketInfo)
}

// @Id				deleteBucket
// @Router			/storage/{bucket} [delete]
// @Description	"Delete a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket			path		string						true	"Name of the bucket"
// @Param			x-storage-key	header		string						true	"The internal storage key"
// @Success		200				{object}	schema.APIResponse[int]		"The number of objects that were in the deleted bucket"
// @Failure		500				{object}	schema.APIResponse[string]	"A string describing the error"
func DeleteBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		respondWithInternalError(c, err)
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
			respondWithInternalError(c, err)
			return
		}
		if err := bucketHandle.Object(objAttrs.Name).Delete(ctx); err != nil {
			respondWithInternalError(c, err)
			return
		}
		deletedCount++
	}

	// Delete bucket
	if err := bucketHandle.Delete(ctx); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", deletedCount)
}

// @Id				objectInfo
// @Router			/storage/{bucket}/{objectID} [get]
// @Description	"Get info on an object in a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket			path		string									true	"Name of the bucket"
// @Param			objectID		path		string									true	"ID of the object"
// @Param			x-storage-key	header		string									true	"The internal storage key"
// @Success		200				{object}	schema.APIResponse[schema.ObjectInfo]	"The object's info"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
func ObjectInfo(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	objectHandle := bucketHandle.Object(objectID)
	if objectHandle == nil {
		respondWithInternalError(c, err)
		return
	}

	// Get object attributes
	attrs, err := objectHandle.Attrs(ctx)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	objectInfo := schema.ObjectInfoFromAttrs(attrs)
	respond(c, http.StatusOK, "success", objectInfo)
}

// @Id				postObject
// @Router			/storage/{bucket}/{objectID} [post]
// @Description	"Upload an object to a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket			path		string									true	"Name of the bucket"
// @Param			objectID		path		string									true	"ID of the object"
// @Param			data			body		string									true	"The data to upload"
// @Param			x-storage-key	header		string									true	"The internal storage key"
// @Success		200				{object}	schema.APIResponse[schema.ObjectInfo]	"The object's info"
// @Failure		500				{object}	schema.APIResponse[string]				"A string describing the error"
func PostObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Read body as byte stream
	fileReader := c.Request.Body
	if fileReader == nil {
		respond(c, http.StatusBadRequest, "error", "Empty body")
		return
	}
	defer fileReader.Close()

	objectHandle := bucketHandle.Object(objectID)
	if objectHandle == nil {
		respondWithInternalError(c, err)
		return
	}

	wc := objectHandle.NewWriter(ctx)
	// Makes object public
	// Set metadata
	wc.CacheControl = "public, max-age=3600"

	// Upload
	if _, err := io.Copy(wc, fileReader); err != nil {
		respondWithInternalError(c, err)
		return
	}

	if err := wc.Close(); err != nil {
		respondWithInternalError(c, err)
		return
	}

	attrs, err := objectHandle.Attrs(ctx)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	objectInfo := schema.ObjectInfoFromAttrs(attrs)
	respond(c, http.StatusOK, "success", objectInfo)
}

// @Id				deleteObject
// @Router			/storage/{bucket}/{objectID} [delete]
// @Description	"Delete an object from a bucket. This route is restricted to only Nebula Labs internal Projects."
// @Param			bucket			path		string						true	"Name of the bucket"
// @Param			objectID		path		string						true	"ID of the object"
// @Param			x-storage-key	header		string						true	"The internal storage key"
// @Success		200				{object}	schema.APIResponse[int]		"Placeholder response, always set to 1"
// @Failure		500				{object}	schema.APIResponse[string]	"A string describing the error"
func DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")
	client := getClient(c)
	ctx := context.Background()

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	objectHandle := bucketHandle.Object(objectID)
	if objectHandle == nil {
		respond(c, http.StatusInternalServerError, "error", "invalid object id")
		return
	}
	err = objectHandle.Delete(ctx)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", 1)
}

// @Id				objectUploadURL
// @Router			/storage/{bucket}/{objectID}/url [put]
// @Accept			json
// @Description	"Create's a new signed URL for target object"
// @Param			bucket			path		string						true	"Name of the bucket"
// @Param			objectID		path		string						true	"ID of the object"
// @Param			body			body		schema.ObjectSignedURLBody	true	"Request body"
// @Param			x-storage-key	header		string						true	"The internal storage key"
// @Success		200				{object}	schema.APIResponse[string]	"Presigned url for the target Object"
// @Failure		500				{object}	schema.APIResponse[string]	"A string describing the error"
func ObjectSignedURL(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	var body schema.ObjectSignedURLBody
	client := getClient(c)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		respond(c, http.StatusBadRequest, "error", "Bad Request Syntax")
		return
	}

	expirationTime, err := time.Parse(time.RFC3339, body.Expiration)
	if err != nil {
		respond(c, http.StatusBadRequest, "error", "Malformatted expiration time")
		return
	}

	bucketHandle, err := getOrCreateBucket(client, bucket)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  body.Method,
		Headers: body.Headers,
		Expires: expirationTime,
	}

	url, err := bucketHandle.SignedURL(objectID, opts)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", url)
}
