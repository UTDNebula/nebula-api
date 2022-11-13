package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func DegreeRoute(router *gin.Engine) {
	// All routes related to degrees come here
	degreeGroup := router.Group("/degree")

	degreeGroup.OPTIONS("", controllers.Preflight)
	degreeGroup.GET("", controllers.DegreeSearch())
	degreeGroup.GET(":id", controllers.DegreeById())
}
