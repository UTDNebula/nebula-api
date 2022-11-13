package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func SectionRoute(router *gin.Engine) {
	// All routes related to sections come here
	sectionGroup := router.Group("/section")

	sectionGroup.OPTIONS("", controllers.Preflight)
	sectionGroup.GET("", controllers.SectionSearch())
	sectionGroup.GET(":id", controllers.SectionById())
}
