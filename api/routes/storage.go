package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func StorageRoute(router *gin.Engine) {
	// All routes related to storage come here
	storageGroup := router.Group("/storage")

	storageGroup.OPTIONS("", controllers.Preflight)
	/*
		storageGroup.GET(":bucket", controllers.BucketInfo())
		storageGroup.GET(":bucket/info/:objectID", controllers.ObjectInfo())
		storageGroup.POST(":bucket/post/:objectID", controllers.PostObject())
		storageGroup.GET(":bucket/get/:objectID", controllers.GetObject())
	*/
}
