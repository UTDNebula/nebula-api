package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func GradesRoute(router *gin.Engine) {
	// All routes related to sections come here
	gradesGroup := router.Group("/grades")

	gradesGroup.OPTIONS("", controllers.Preflight)

	// @TODO: Do we need this?
	// ---- gradesGroup.OPTIONS("semester", controllers.Preflight)
	// ---- gradesGroup.OPTIONS("overall", controllers.Preflight)

	gradesGroup.GET("semester", controllers.GradeAggregationSemester())
	gradesGroup.GET("overall", controllers.GradesAggregationOverall())
}
