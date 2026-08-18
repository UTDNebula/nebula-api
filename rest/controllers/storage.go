package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
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

		if errors.Is(err, storage.ErrBucketNotExist) {
			err = bucketHandle.Create(ctx, PROJECT_ID, nil)
			if err != nil {
				return nil, errors.New("failed to create bucket: " + err.Error())
			}
		} else {
			return nil, err
		}
	}
	return bucketHandle, nil
}

// @Id				bucketInfo
// @Router			/storage/{bucket} [get]
// @Tags			Internal
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
// @Tags			Internal
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
// @Tags			Internal
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
		if errors.Is(err, storage.ErrObjectNotExist) {
			respond(c, http.StatusNotFound, "error", "Object with given ID not found")
		} else {
			respondWithInternalError(c, err)
		}
		return
	}

	// Generate public URL
	escapedObject := url.PathEscape(objectID)
	url := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		schema.BUCKET_PREFIX+bucket,
		escapedObject,
	)

	objectInfo := schema.ObjectInfoFromAttrs(attrs, url)
	respond(c, http.StatusOK, "success", objectInfo)
}

// @Id				postObject
// @Router			/storage/{bucket}/{objectID} [post]
// @Tags			Internal
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

	maxUploadSize := configs.GetEnvMaxUploadSize()

	// Force early 413 check via Content-Length if present
	if c.Request.ContentLength > maxUploadSize {
		respond(c, http.StatusRequestEntityTooLarge, "error", fmt.Sprintf("File too large. Maximum allowed size is %d bytes (%dMB)", maxUploadSize, maxUploadSize/(1024*1024)))
		return
	}

	// Use MaxBytesReader to limit the body
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// Read and validate the entire (capped) request body before touching GCS.

	fileBytes, readErr := io.ReadAll(c.Request.Body)
	if readErr != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(readErr, &maxBytesErr) {
			respond(c, http.StatusRequestEntityTooLarge, "error", fmt.Sprintf("File too large. Maximum allowed size is %d bytes (%dMB)", maxUploadSize, maxUploadSize/(1024*1024)))
			return
		}
		respondWithInternalError(c, readErr)
		return
	}

	fileReader := bytes.NewReader(fileBytes)

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

	// Get object attributes
	attrs, err := objectHandle.Attrs(ctx)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Generate public URL
	escapedObject := url.PathEscape(objectID)
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", schema.BUCKET_PREFIX+bucket, escapedObject)

	objectInfo := schema.ObjectInfoFromAttrs(attrs, url)
	respond(c, http.StatusOK, "success", objectInfo)
}

// @Id				deleteObject
// @Router			/storage/{bucket}/{objectID} [delete]
// @Tags			Internal
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
	if err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", 1)
}

// @Id				objectUploadURL
// @Router			/storage/{bucket}/{objectID}/url [put]
// @Tags			Internal
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

	headers := append([]string{}, body.Headers...)
	// Upload size limits for signed URL uploads.
	if strings.EqualFold(body.Method, http.MethodPut) || strings.EqualFold(body.Method, http.MethodPost) {
		maxUploadSize := configs.GetEnvMaxUploadSize()
		hasContentLengthRange := false
		for _, header := range headers {
			if strings.HasPrefix(strings.ToLower(header), "x-goog-content-length-range:") {
				hasContentLengthRange = true
				break
			}
		}

		if !hasContentLengthRange {
			headers = append(headers, fmt.Sprintf("x-goog-content-length-range:0,%d", maxUploadSize))
		}
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  body.Method,
		Headers: headers,
		Expires: expirationTime,
	}

	url, err := bucketHandle.SignedURL(objectID, opts)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", url)
}
