package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func GradesRoute(router *gin.Engine, repo shared.GradesRepo) {
	handler := controllers.NewGradesController(repo)
	// All routes related to sections come here
	gradesGroup := router.Group("/grades")

	gradesGroup.OPTIONS("", controllers.Preflight)

	//	@TODO:	Do we need this?
	// ---- gradesGroup.OPTIONS("semester", controllers.Preflight)
	// ---- gradesGroup.OPTIONS("overall", controllers.Preflight)

	gradesGroup.GET("semester", handler.GradeAggregationSemester)
	gradesGroup.GET("semester/sectionType", handler.GradesAggregationSectionType)
	gradesGroup.GET("overall", handler.GradesAggregationOverall)
}
