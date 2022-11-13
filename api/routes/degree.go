package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func DegreeRoute(router *gin.Engine) {
	// All routes related to degrees come here
	router.OPTIONS("/degree", controllers.Preflight)
	degreeGroup := router.Group("/degree")

	degreeGroup.GET("", controllers.DegreeSearch())
	degreeGroup.GET(":id", controllers.DegreeById())
}
