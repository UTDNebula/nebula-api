package routes

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"

	"cloud.google.com/go/storage"
	"github.com/UTDNebula/nebula-api/api/controllers"
	"google.golang.org/api/option"
)

var (
	client     *storage.Client
	clientOnce sync.Once
)

func initStorageClient() *storage.Client {
	clientOnce.Do(func() {
		encodedCreds, exist := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
		if !exist {
			log.Fatal("Error loading 'GOOGLE_APPLICATION_CREDENTIALS' from the .env file")
		}
		ctx := context.Background()
		var err error
		client, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(encodedCreds)))
		if err != nil {
			log.Fatalf("Failed to create GCS client: %v", err)
		}
	})
	return client
}

func StorageRoute(router *gin.Engine) {
	storageClient := initStorageClient()
	router.Use(func(c *gin.Context) {
		c.Set("gcsClient", storageClient)
		c.Next()
	})

	// All routes related to storage come here
	storageGroup := router.Group("/storage")

	storageGroup.OPTIONS("", controllers.Preflight)
	storageGroup.GET(":bucket", controllers.BucketInfo)
	storageGroup.DELETE(":bucket", controllers.DeleteBucket)
	storageGroup.GET(":bucket/:objectID", controllers.ObjectInfo)
	storageGroup.POST(":bucket/:objectID", controllers.PostObject)
	storageGroup.DELETE(":bucket/:objectID", controllers.DeleteObject)
}
