package events

import (
	"net/http"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func MazevoRoute(router *gin.Engine, repository shared.EventRepository[schema.MazevoEvent]) {
	handler := mazevoController{repository: repository}
	//All routes related to mazevo events come here
	mazevoGroup := router.Group("/mazevo")

	mazevoGroup.OPTIONS("", func(c *gin.Context) { c.JSON(http.StatusOK, struct{}{}) })
	mazevoGroup.GET(":date", handler.MazevoEvents)
	mazevoGroup.GET(":date/:building", handler.MazevoEventsByBuilding)
	mazevoGroup.GET(":date/:building/:room", handler.MazevoEventsByRoom)
}
