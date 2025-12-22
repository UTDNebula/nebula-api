package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func SectionRoute(router *gin.Engine) {
	// All routes related to sections come here
	sectionGroup := router.Group("/section")

	sectionGroup.OPTIONS("", controllers.Preflight)
	sectionGroup.GET("", controllers.SectionSearch)
	sectionGroup.GET(":id", controllers.SectionById)
	//sectionGroup.GET(":id/evaluation", controllers.EvalBySectionID)

	// Endpoints for aggregate
	sectionGroup.GET("/courses", controllers.SectionCourseSearch)
	sectionGroup.GET(":id/course", controllers.SectionCourseById)
	sectionGroup.GET("/professors", controllers.SectionProfessorSearch)
	sectionGroup.GET("/:id/professors", controllers.SectionProfessorById)

	// Route for section grades
	sectionGroup.GET(":id/grades", controllers.GradesBySectionID())
}
