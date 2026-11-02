package responses

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// Sets the API's response to a request, producing valid JSON given a status code and data.
func Respond[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(
		status,
		schema.APIResponse[T]{
			Status:  status,
			Message: message,
			Data:    data,
		},
	)
}

// Helper function for logging and responding to a generic internal server error.
func RespondWithInternalError(c *gin.Context, err error) {
	// Note that we use log.Output here to be able to set the stack depth to the frame above this one (2),
	// which allows us to log the location this function was called from
	log.Output(2, fmt.Sprintf("INTERNAL SERVER ERROR: %s", err.Error()))
	// Capture error with Sentry
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			hub.CaptureException(err)
		})
	}
	Respond(c, http.StatusInternalServerError, "error", err.Error())
}
