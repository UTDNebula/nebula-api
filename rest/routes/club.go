package routes

import (
	"log"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func ClubRoute(router *gin.Engine) {
	if configs.ConnectClubsDB() == nil {
		log.Println("Clubs database is not configured; skipping club routes")
		return
	}

	// All routes related to clubs come here
	clubGroup := router.Group("/club")

	clubGroup.OPTIONS("", controllers.Preflight)
	clubGroup.GET(":id", controllers.ClubDirectoryInfo)
	clubGroup.GET("/search", controllers.ClubSearch)
}
