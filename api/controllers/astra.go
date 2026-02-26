package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/api/configs"

	"github.com/UTDNebula/nebula-api/api/schema"
)

var astraCollection *mongo.Collection = configs.GetCollection("astra")

// @Id				AstraEvents
// @Router			/astra/{date} [get]
// @Tags			Events
// @Description	"Returns AstraEvent based on the input date"
// @Produce		json
// @Param			date	path		string																true	"date (ISO format) to retrieve astra events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.AstraEvent]]	"All AstraEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]											"A string describing the error"
func AstraEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	var astra_events schema.MultiBuildingEvents[schema.AstraEvent]

	// Find astra event given date
	err := astraCollection.FindOne(ctx, bson.M{"date": date}).Decode(&astra_events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			astra_events.Date = date
			astra_events.Buildings = []schema.SingleBuildingEvents[schema.AstraEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	respond(c, http.StatusOK, "success", astra_events)
}

// @Id				AstraEventsByBuilding
// @Router			/astra/{date}/{building} [get]
// @Tags			Events
// @Description	"Returns AstraEvent based on the input date and building name"
// @Produce		json
// @Param			date		path		string																true	"date (ISO format) to retrieve astra events"
// @Param			building	path		string																true	"building abbreviation of event locations"
// @Success		200			{object}	schema.APIResponse[schema.SingleBuildingEvents[schema.AstraEvent]]	"All sections with meetings on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]											"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]											"A string describing the error"
func AstraEventsByBuilding(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")

	var astra_events schema.MultiBuildingEvents[schema.AstraEvent]
	var astra_eventsByBuilding schema.SingleBuildingEvents[schema.AstraEvent]

	// Find astra event given date
	err := astraCollection.FindOne(ctx, bson.M{"date": date}).Decode(&astra_events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			astra_events.Date = date
			astra_events.Buildings = []schema.SingleBuildingEvents[schema.AstraEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	//parse response for requested building
	for _, b := range astra_events.Buildings {
		if b.Building == building {
			astra_eventsByBuilding = b
			break
		}
	}

	if astra_eventsByBuilding.Building == "" {
		respond(c, http.StatusNotFound, "error", "No events found for the specified building")
		return
	}

	respond(c, http.StatusOK, "success", astra_eventsByBuilding)
}

// @Id				AstraEventsByBuildingandRoom
// @Router			/astra/{date}/{building}/{room} [get]
// @Tags			Events
// @Description	"Returns AstraEvent based on the input date building name and room number"
// @Produce		json
// @Param			date		path		string																true	"date (ISO format) to retrieve astra events"
// @Param			building	path		string																true	"building abbreviation of event locations"
// @Param			room		path		string																true	"room number for event"
// @Success		200			{object}	schema.APIResponse[schema.SingleBuildingEvents[schema.AstraEvent]]	"All sections with meetings on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]											"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]											"A string describing the error"
func AstraEventsByBuildingAndRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")
	room := c.Param("room")

	var astra_events schema.MultiBuildingEvents[schema.AstraEvent]
	var roomEvents schema.RoomEvents[schema.AstraEvent]

	// Find astra event given date
	err := astraCollection.FindOne(ctx, bson.M{"date": date}).Decode(&astra_events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			astra_events.Date = date
			astra_events.Buildings = []schema.SingleBuildingEvents[schema.AstraEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	//parse response for requested building and room
	for _, b := range astra_events.Buildings {
		if b.Building == building {
			for _, r := range b.Rooms {
				if r.Room == room {
					roomEvents = r
					break
				}
			}
			break
		}
	}

	if roomEvents.Room == "" {
		respond(c, http.StatusNotFound, "error", "No rooms found for the specified building or event")
		return
	}

	respond(c, http.StatusOK, "success", roomEvents)
}
