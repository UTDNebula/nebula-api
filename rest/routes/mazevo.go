package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func MazevoRoute(router *gin.Engine, repository shared.EventRepository[schema.MazevoEvent]) {
	handler := controllers.MazevoController{Repository: repository}
	//All routes related to mazevo events come here
	mazevoGroup := router.Group("/mazevo")

	mazevoGroup.OPTIONS("", controllers.Preflight)
	mazevoGroup.GET(":date", handler.MazevoEvents)
	mazevoGroup.GET(":date/:building", handler.MazevoEventsByBuilding)
	mazevoGroup.GET(":date/:building/:room", handler.MazevoEventsByRoom)
}
