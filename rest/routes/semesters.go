package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func SemestersRoute(router *gin.Engine) {
	semestersGroup := router.Group("/semesters")

	semestersGroup.OPTIONS("", controllers.Preflight)
	semestersGroup.GET("", controllers.Semesters)
}
