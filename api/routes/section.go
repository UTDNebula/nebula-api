package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func SectionRoute(router *gin.Engine) {
	// All routes related to sections come here
	router.OPTIONS("/section", controllers.Preflight)
	sectionGroup := router.Group("/section")

	sectionGroup.GET("", controllers.SectionSearch())
	sectionGroup.GET(":id", controllers.SectionById())
}
