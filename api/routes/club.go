package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func ClubRoute(router *gin.Engine) {
	// All routes related to courses come here
	clubGroup := router.Group("/club")

	clubGroup.OPTIONS("", controllers.Preflight)
	clubGroup.GET(":id", controllers.ClubDirectoryInfo)
}
