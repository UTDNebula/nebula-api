package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"

	"cloud.google.com/go/storage"
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/UTDNebula/nebula-api/api/schema"

	"google.golang.org/api/option"
)

// Singleton client, not to be changed
var client *storage.Client

// To prevent changing above singleton
var clientOnce sync.Once

func initStorageClient() *storage.Client {
	// Only do once
	clientOnce.Do(func() {
		ctx := context.Background()
		var c *storage.Client
		var err error
		_, exist := os.LookupEnv("USE_CLOUD_CREDS")

		// If USE_CLOUD_CREDS env var set, assume we're running on cloud and don't need to set creds
		if exist {
			c, err = storage.NewClient(ctx)
		} else {
			// We're not running on the cloud, get JSON service account key from .env
			encodedCreds, exist := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
			if !exist {
				log.Println("Error loading 'GOOGLE_APPLICATION_CREDENTIALS' from the .env file, skipping cloud storage routes")
				return
			}
			c, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(encodedCreds)))
		}
		if err != nil {
			log.Printf("Failed to create GCS client: %v", err)
			return
		}
		client = c
	})
	return client
}

func StorageRoute(router *gin.Engine) {
	// Create client, don't procede if error
	storageClient := initStorageClient()
	if storageClient == nil {
		log.Println("GCS client not initialized; skipping cloud storage routes")
		return
	}

	//Rescrict with password
	authMiddleware := func(c *gin.Context) {
		secret := c.GetHeader("x-storage-key")
		expected, exist := os.LookupEnv("STORAGE_ROUTE_KEY")
		if !exist || secret != expected {
			c.AbortWithStatusJSON(http.StatusForbidden, schema.APIResponse[string]{Status: http.StatusForbidden, Message: "error", Data: "Forbidden"})
			return
		}
		c.Next()
	}

	// Pass to next layer
	router.Use(func(c *gin.Context) {
		c.Set("gcsClient", storageClient)
		c.Next()
	})

	// All routes related to storage come here
	storageGroup := router.Group("/storage")

	//Use auth
	storageGroup.Use(authMiddleware)

	storageGroup.OPTIONS("", controllers.Preflight)
	storageGroup.GET(":bucket", controllers.BucketInfo)
	storageGroup.DELETE(":bucket", controllers.DeleteBucket)
	storageGroup.GET(":bucket/:objectID", controllers.ObjectInfo)
	storageGroup.POST(":bucket/:objectID", controllers.PostObject)
	storageGroup.DELETE(":bucket/:objectID", controllers.DeleteObject)
}
