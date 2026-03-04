package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func DiscountProgramRoute(router *gin.Engine) {
	// All routes related to discount programs come here
	discountProgramGroup := router.Group("/discountPrograms")

	discountProgramGroup.OPTIONS("", controllers.Preflight)
	discountProgramGroup.GET("", controllers.DiscountProgramSearch)
}
