package routes

import (
	"log"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func ClubsRoute(router *gin.Engine) {
	db := configs.ConnectClubsDB()

	// Skip accessing Clubs DB if CLUBS_DB_URI is not provided.
	if db == nil {
		log.Println("Skipping club routes.")
		return
	}

	// All routes related to UTD Clubs come here
	clubGroup := router.Group("/clubs")

	clubGroup.OPTIONS("", controllers.Preflight)
	clubGroup.GET("", controllers.Clubs)
	clubGroup.GET("/:id", controllers.ClubById)
	clubGroup.GET("/:id/events", controllers.ClubEvents)
	clubGroup.GET("/events", controllers.ClubsEvents)
	clubGroup.GET("/events/:id", controllers.ClubsEventById)
	clubGroup.GET("/tags", controllers.ClubsTags)
}
