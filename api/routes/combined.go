package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func CombinedRoute(router *gin.Engine) {
	combinedGroup := router.Group("/combined")
	combinedGroup.OPTIONS("", controllers.Preflight)
	combinedGroup.GET("/sections/trends", controllers.TrendsCombinedSectionSearch)
}
