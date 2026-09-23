package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func AcademicCalendarRoute(router *gin.Engine) {
	//All routes related to academic calendars route through here
	academicCalendarsGroup := router.Group("/academicCalendars")

	academicCalendarsGroup.OPTIONS("", controllers.Preflight)
	academicCalendarsGroup.GET(":id", controllers.AcademicCalendarsById)
	academicCalendarsGroup.GET("/current", controllers.AcadenemicCalendarsCurrent)
}
