package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func AstraRoute(router *gin.Engine) {
	//All routes related to astra events come here
	astraGroup := router.Group("/astra")

	astraGroup.OPTIONS("", controllers.Preflight)
	astraGroup.GET(":date", controllers.AstraEvents)
}
