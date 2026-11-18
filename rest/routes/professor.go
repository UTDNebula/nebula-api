package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func ProfessorRoute(router *gin.Engine, repo shared.ProfessorRepo) {
	handler := controllers.NewProfessorController(repo)
	// All routes related to professors come here
	professorGroup := router.Group("/professor")

	professorGroup.OPTIONS("", controllers.Preflight)
	professorGroup.GET("", handler.ProfessorSearch)
	professorGroup.GET(":id", handler.ProfessorById)
	professorGroup.GET("all", handler.ProfessorAll)

	// Endpoints to get the courses of the professors
	professorGroup.GET("courses", handler.ProfessorCourseSearch)
	professorGroup.GET(":id/courses", handler.ProfessorCourseById)

	// Endpoints to get the sections of the professors
	professorGroup.GET("sections", handler.ProfessorSectionSearch)
	professorGroup.GET(":id/sections", handler.ProfessorSectionById)
	professorGroup.GET("sections/trends", handler.TrendsProfessorSectionSearch)

	// Endpoints to get the grades for a professor
	professorGroup.GET(":id/grades", handler.GradesByProfessorID)
}
