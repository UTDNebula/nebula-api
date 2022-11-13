package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func ProfessorRoute(router *gin.Engine) {
	// All routes related to professors come here
	router.OPTIONS("/professor", controllers.Preflight)
	professorGroup := router.Group("/professor")

	professorGroup.GET("", controllers.ProfessorSearch())
	professorGroup.GET(":id", controllers.ProfessorById())
}
