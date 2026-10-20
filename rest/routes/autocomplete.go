package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/rest/controllers"
)

func AutocompleteRoute(router *gin.Engine) {
	// All routes related to autocomplete come here
	autocompleteGroup := router.Group("/autocomplete")

	autocompleteGroup.GET("/dag", controllers.AutocompleteDAG)
}
