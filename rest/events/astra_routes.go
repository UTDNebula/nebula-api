package events

import (
	"net/http"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func AstraRoute(router *gin.Engine, repository shared.EventRepository[schema.AstraEvent]) {
	handler := astraController{repository: repository}
	//All routes related to astra events come here
	astraGroup := router.Group("/astra")

	astraGroup.OPTIONS("", func(c *gin.Context) { c.JSON(http.StatusOK, struct{}{}) })
	astraGroup.GET(":date", handler.AstraEvents)
	astraGroup.GET(":date/:building", handler.AstraEventsByBuilding)
	astraGroup.GET(":date/:building/:room", handler.AstraEventsByBuildingAndRoom)
}
