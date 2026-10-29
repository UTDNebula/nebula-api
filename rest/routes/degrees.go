package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func DegreesRoute(router *gin.Engine) {
	// All routes related to degrees come here
	degreesGroup := router.Group("/degrees")

	degreesGroup.OPTIONS("", controllers.Preflight)
	degreesGroup.GET("", controllers.Degrees)
}
