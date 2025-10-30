package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// Creates a context with the specified timeout and returns both context and cancel function.
// Common timeouts: 10s for standard queries, 30s for "all" operations.
func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// Generic function to handle Find operations with pagination.
// Reduces boilerplate for search endpoints by handling query building, finding, decoding, and responding.
func findAndRespond[T any](c *gin.Context, collection *mongo.Collection, timeout time.Duration) {
	ctx, cancel := createContext(timeout)
	defer cancel()

	var results []T

	// Build query key-value pairs
	query, err := getQuery[T]("Search", c)
	if err != nil {
		return // getQuery already responds with error
	}

	// Get pagination options
	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	// Execute find query
	cursor, err := collection.Find(ctx, query, optionLimit)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Decode all results
	if err = cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Return results
	respond(c, http.StatusOK, "success", results)
}

// Generic function to handle FindOne operations by ID.
// Reduces boilerplate for ById endpoints by handling query building, finding one, decoding, and responding.
func findOneByIdAndRespond[T any](c *gin.Context, collection *mongo.Collection, timeout time.Duration) {
	ctx, cancel := createContext(timeout)
	defer cancel()

	var result T

	// Parse object ID from parameter
	query, err := getQuery[T]("ById", c)
	if err != nil {
		return // getQuery already responds with error
	}

	// Find and decode matching document
	err = collection.FindOne(ctx, query).Decode(&result)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Return result
	respond(c, http.StatusOK, "success", result)
}

// Generic function to handle FindAll operations without filters.
// Reduces boilerplate for "all" endpoints by finding all documents and responding.
func findAllAndRespond[T any](c *gin.Context, collection *mongo.Collection, timeout time.Duration) {
	ctx, cancel := createContext(timeout)
	defer cancel()

	var results []T

	// Find all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Decode all results
	if err = cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Return results
	respond(c, http.StatusOK, "success", results)
}

// Generic function to handle Aggregate operations.
// Reduces boilerplate for aggregate endpoints by executing pipeline, decoding, and responding.
func aggregateAndRespond[T any](c *gin.Context, collection *mongo.Collection, pipeline mongo.Pipeline, timeout time.Duration) {
	ctx, cancel := createContext(timeout)
	defer cancel()

	var results []T

	// Execute aggregation pipeline
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Decode all results
	if err = cursor.All(ctx, &results); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Return results
	respond(c, http.StatusOK, "success", results)
}