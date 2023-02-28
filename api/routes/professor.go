package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func ProfessorRoute(router *gin.Engine) {
	// All routes related to professors come here
	professorGroup := router.Group("/professor")

	professorGroup.OPTIONS("", controllers.Preflight)
	professorGroup.GET("", controllers.ProfessorSearch())
	professorGroup.GET(":id", controllers.ProfessorById())
}
