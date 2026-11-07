package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func SemestersRoute(router *gin.Engine) {
	// All routes related to semesters come here
	semestersGroup := router.Group("/semesters")

	semestersGroup.OPTIONS("", controllers.Preflight)
	semestersGroup.GET("", controllers.Semesters)
}
