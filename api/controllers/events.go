package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventsCollection *mongo.Collection = configs.GetCollection("events")

// @Id events
// @Router /events/{date} [get]
// @Description "Returns all sections with meetings on the specified date"
// @Produce json
// @Param date path string true "ISO date of the set of events to get"
// @Success 200 {array} schema.MultiBuildingEvents "All sections with meetings on the specified date"
func Events(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	date := c.Param("date")

	var events schema.MultiBuildingEvents

	defer cancel()

	// find and parse matching date
	err := eventsCollection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.MultiBuildingEventsResponse{Status: http.StatusOK, Message: "success", Data: events})
}

// @Id eventsByBuilding
// @Router /events/{date}/{building} [get]
// @Description "Returns all sections with meetings on the specified date in the specified building"
// @Produce json
// @Param date path string true "ISO date of the set of events to get"
// @Param building path string true "building abbreviation of event locations"
// @Success 200 {array} schema.SingleBuildingEvents "All sections with meetings on the specified date in the specified building"
func EventsByBuilding(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	date := c.Param("date")
	building := c.Param("building")

	var events schema.MultiBuildingEvents
	var eventsByBuilding schema.SingleBuildingEvents

	defer cancel()

	// find and parse matching date
	err := eventsCollection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// filter for the specified building
	for _, b := range events.Buildings {
		if b.Building == building {
			eventsByBuilding = b
			break
		}
	}

	// If no building is found, return an error
	if eventsByBuilding.Building == "" {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "No events found for the specified building",
		})
		return
	}

	c.JSON(http.StatusOK, responses.SingleBuildingEventsResponse{Status: http.StatusOK, Message: "success", Data: eventsByBuilding})
}
