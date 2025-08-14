package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func GamesRoute(router *gin.Engine) {
	// All routes related to games come here
	gamesGroup := router.Group("/games")

	gamesGroup.OPTIONS("", controllers.Preflight)
	gamesGroup.GET("/letters/:date", controllers.Letters)
}
