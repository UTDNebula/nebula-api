package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func SectionRoute(router *gin.Engine, repo shared.SectionRepo) {
	handler := controllers.NewSectionController(repo)
	// All routes related to sections come here
	sectionGroup := router.Group("/section")

	sectionGroup.OPTIONS("", controllers.Preflight)
	sectionGroup.GET("", handler.SectionSearch)
	sectionGroup.GET(":id", handler.SectionById)
	//sectionGroup.GET(":id/evaluation", controllers.EvalBySectionID)

	// Endpoints for aggregate
	sectionGroup.GET("/courses", handler.SectionCourseSearch)
	sectionGroup.GET(":id/course", handler.SectionCourseById)
	sectionGroup.GET("/professors", handler.SectionProfessorSearch)
	sectionGroup.GET("/:id/professors", handler.SectionProfessorById)

	// Route for section grades
	sectionGroup.GET(":id/grades", handler.GradesBySectionID)
}
