package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func DiscountRoutes(router *gin.Engine) {
	// All routes related to discounts come here
	discountGroup := router.Group("/discountPrograms")

	discountGroup.GET("", controllers.DiscountSearch)
}
