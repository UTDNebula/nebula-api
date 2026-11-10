package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func ClubRoute(router *gin.Engine) {
	// All routes related to UTD Clubs come here
	clubGroup := router.Group("/clubs")

	clubGroup.OPTIONS("", controllers.Preflight)
	clubGroup.GET("", controllers.Clubs)
	clubGroup.GET("/:id", controllers.ClubById)
	clubGroup.GET("/:id/events", controllers.ClubEvents)
	clubGroup.GET("/events", controllers.ClubsEvents)
	clubGroup.GET("/events/:id", controllers.ClubsEventById)
}
