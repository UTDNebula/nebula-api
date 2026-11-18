package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func CombinedRoute(router *gin.Engine, repo shared.TrendsRepo) {
	handler := controllers.NewTrendsController(repo)
	combinedGroup := router.Group("/combined")
	combinedGroup.OPTIONS("", controllers.Preflight)
	combinedGroup.GET("/sections/trends", handler.TrendsCombinedSectionSearch)
}
