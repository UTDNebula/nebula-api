package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func EventsRoute(router *gin.Engine) {
	// All routes related to sections come here
	eventsGroup := router.Group("/events/coursebook")

	eventsGroup.OPTIONS("", controllers.Preflight)
	eventsGroup.GET(":date", controllers.Events)
	eventsGroup.GET(":date/:building", controllers.EventsByBuilding)
}
