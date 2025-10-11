package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sets the API's response to a request, producing valid JSON given a status code and data.
func respond[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(
		status,
		schema.APIResponse[T]{
			Status:  status,
			Message: message,
			Data:    data,
		},
	)
}

// Builds a MongoDB filter for type T based on the given flag search or byid
func getQuery[T any](flag string, c *gin.Context) (bson.M, error) {
	switch flag {
	case "Search":
		q, err := schema.FilterQuery[T](c)
		if err != nil {
			respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
			return nil, err
		}
		return q, nil

	case "ById":
		objId, err := objectIDFromParam(c, "id")
		if err != nil {
			// objectIDFromParam already responds with 400 if conversion fails
			return nil, err
		}
		return bson.M{"_id": objId}, nil

	default:
		err := fmt.Errorf("invalid flag for getQuery: %s", flag)
		respondWithInternalError(c, err)
		return nil, err
	}
}

// Helper function for logging and responding to a generic internal server error.
func respondWithInternalError(c *gin.Context, err error) {
	// Note that we use log.Output here to be able to set the stack depth to the frame above this one (2),
	// which allows us to log the location this function was called from
	log.Output(2, fmt.Sprintf("INTERNAL SERVER ERROR: %s", err.Error()))
	// Capture error with Sentry
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			hub.CaptureException(err)
		})
	}
	respond(c, http.StatusInternalServerError, "error", err.Error())
}

// Attempts to convert the given parameter to an ObjectID for use with MongoDB.
// Automatically responds with http.StatusBadRequest if conversion fails.
func objectIDFromParam(c *gin.Context, paramName string) (*primitive.ObjectID, error) {
	idHex := c.Param(paramName)
	objectId, convertIdErr := primitive.ObjectIDFromHex(idHex)
	if convertIdErr != nil {
		// Respond with an error if we can't covert successfully
		log.Println(convertIdErr)
		respond(c,
			http.StatusBadRequest,
			fmt.Sprintf("Parameter \"%s\" is not a valid ObjectID.", paramName),
			convertIdErr.Error(),
		)
		return nil, convertIdErr
	}
	return &objectId, nil
}
