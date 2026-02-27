package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func EmailRoute(router *gin.Engine) {
	// All routes related to email come here
	emailGroup := router.Group("/email")

	emailGroup.OPTIONS("", controllers.Preflight)
	emailGroup.POST("", controllers.SendEmail)
}
