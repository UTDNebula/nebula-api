package routes

import (
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/gin-gonic/gin"
)

func BudgetRoute(router *gin.Engine) {
	//All routes related to budgets come here
	budgetGroup := router.Group("/budget")

	budgetGroup.OPTIONS("", controllers.Preflight)
	budgetGroup.GET(":year", controllers.Budget)
}
