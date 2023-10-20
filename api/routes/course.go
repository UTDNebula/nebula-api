package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func CourseRoute(router *gin.Engine) {
	// All routes related to courses come here
	courseGroup := router.Group("/course")

	courseGroup.OPTIONS("", controllers.Preflight)
	courseGroup.GET("", controllers.CourseSearch())
	courseGroup.GET("all", controllers.CourseAll())
	courseGroup.GET(":id", controllers.CourseById())
}
