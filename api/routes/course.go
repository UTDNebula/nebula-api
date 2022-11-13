package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func CourseRoute(router *gin.Engine) {
	// All routes related to courses come here
	router.OPTIONS("/course", controllers.Preflight)
	courseGroup := router.Group("/course")

	courseGroup.GET("", controllers.CourseSearch())
	courseGroup.GET(":id", controllers.CourseById())
}
