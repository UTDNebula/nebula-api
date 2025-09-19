package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sets the API's response to a request, producing valid JSON given a status code and data.
func respond[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(status, schema.APIResponse[T]{Status: status, Message: message, Data: data})
}
// Builds a MongoDB filter from request query parameters for the given schema type T.
// Automatically responds with HTTP 400 if the parameters are invalid.
func getQuery[T any](c *gin.Context) (bson.M, error) {
	q, err := schema.FilterQuery[T](c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return nil, err
	}
	return q, nil
}


// Helper function for logging and responding to a generic internal server error.
func respondWithInternalError(c *gin.Context, err error) {
	// Note that we use log.Output here to be able to set the stack depth to the frame above this one (2), which allows us to log the location this function was called from
	log.Output(2, fmt.Sprintf("INTERNAL SERVER ERROR: %s", err.Error()))
	respond(c, http.StatusInternalServerError, "error", err.Error())
}

// Attempts to convert the given parameter to an ObjectID for use with MongoDB. Automatically responds with http.StatusBadRequest if conversion fails.
func objectIDFromParam(c *gin.Context, paramName string) (*primitive.ObjectID, error) {
	idHex := c.Param(paramName)
	objectId, convertIdErr := primitive.ObjectIDFromHex(idHex)
	if convertIdErr != nil {
		// Respond with an error if we can't covert successfully
		log.Println(convertIdErr)
		respond(c, http.StatusBadRequest, fmt.Sprintf("Parameter \"%s\" is not a valid ObjectID.", paramName), convertIdErr.Error())
		return nil, convertIdErr
	}
	return &objectId, nil
}
