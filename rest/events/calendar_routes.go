package events

import (
	"net/http"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func CalendarRoute(router *gin.Engine, repository shared.EventRepository[schema.Event]) {
	handler := calendarController{repository: repository}
	//All routes related to comet calendar events come here
	calendarGroup := router.Group("/calendar")

	calendarGroup.OPTIONS("", func(c *gin.Context) { c.JSON(http.StatusOK, struct{}{}) })
	// More specific routes must be defined first in Gin
	calendarGroup.GET(":date/:building/:room", handler.CometCalendarEventsByBuildingAndRoom)
	calendarGroup.GET(":date/:building", handler.CometCalendarEventsByBuilding)
	calendarGroup.GET(":date", handler.CometCalendarEvents)
}
