package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func RoomsRoute(router *gin.Engine) {
	// All routes related to sections come here
	roomsGroup := router.Group("/rooms")

	roomsGroup.OPTIONS("", controllers.Preflight)
	roomsGroup.GET("", controllers.Rooms)
}
