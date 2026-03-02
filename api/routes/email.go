package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/UTDNebula/nebula-api/api/schema"
)

func EmailRoute(router *gin.Engine) {
	// Rescrict with password
	authMiddleware := func(c *gin.Context) {
		secret := c.GetHeader("x-email-key")
		expected, exist := os.LookupEnv("EMAIL_ROUTE_KEY")
		if !exist || secret != expected {
			c.AbortWithStatusJSON(http.StatusForbidden, schema.APIResponse[string]{Status: http.StatusForbidden, Message: "error", Data: "Forbidden"})
			return
		}
		c.Next()
	}

	// All routes related to email come here
	emailGroup := router.Group("/email")

	// Use auth
	emailGroup.Use(authMiddleware)

	emailGroup.OPTIONS("", controllers.Preflight)
	emailGroup.POST("/send", controllers.SendEmail)
	emailGroup.POST("/queue", controllers.QueueEmail)
}
