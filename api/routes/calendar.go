package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func CalendarRoute(router *gin.Engine) {
	//All routes related to comet calendar events come here
	calendarGroup := router.Group("/calendar")

	calendarGroup.OPTIONS("", controllers.Preflight)
	calendarGroup.GET(":date", controllers.CometCalendarEvents)
	calendarGroup.GET(":date/:building", controllers.CometCalendarEventsByBuilding)
	calendarGroup.GET(":date/:building/:room", controllers.CometCalendarEventsByBuildingAndRoom)
}
