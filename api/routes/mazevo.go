package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func MazevoRoute(router *gin.Engine) {
	//All routes related to mazevo events come here
	mazevoGroup := router.Group("/mazevo")

	mazevoGroup.OPTIONS("", controllers.Preflight)
	mazevoGroup.GET(":date", controllers.MazevoEvents)
}
