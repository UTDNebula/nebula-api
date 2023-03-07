package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func GradesRoute(router *gin.Engine) {
	// All routes related to sections come here
	gradesGroup := router.Group("/grades")

	gradesGroup.OPTIONS("", controllers.Preflight)
	gradesGroup.GET("", controllers.GradesSearch())
	gradesGroup.GET("all", controllers.GradesAll())
}
