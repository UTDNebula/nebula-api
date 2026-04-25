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

var mazevoCollection *mongo.Collection = configs.GetCollection("mazevo")

// @Id				MazevoEvents
// @Router			/mazevo/{date} [get]
// @Tags			Events
// @Description	"Returns MazevoEvent based on the input date"
// @Produce		json
// @Param			date	path		string																true	"date (ISO format) to retrieve mazevo events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.MazevoEvent]]	"All MazevoEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]											"A string describing the error"
func MazevoEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	var mazevoEvents schema.MultiBuildingEvents[schema.MazevoEvent]

	// Find mazevo event for input date
	err := mazevoCollection.FindOne(ctx, bson.M{"date": date}).Decode(&mazevoEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			mazevoEvents.Date = date
			mazevoEvents.Buildings = []schema.SingleBuildingEvents[schema.MazevoEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	respond(c, http.StatusOK, "success", mazevoEvents)
}

// @Id				mazevoEventsByBuilding
// @Router			/mazevo/{date}/{building} [get]
// @Tags			Events
// @Description	"Returns all sections with MazevoEvent meetings on the specified date in the specified building"
// @Produce		json
// @Param			date		path		string																true	"date (ISO format) to retrieve mazevo events"
// @Param			building	path		string																true	"building abbreviation of the event location"
// @Success		200			{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.MazevoEvent]]	"All MazevoEvents sections with meetings on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]											"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]											"A string describing the error"
func MazevoEventsByBuilding(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")

	var mazevoEvents schema.MultiBuildingEvents[schema.SectionWithTime]
	var mazevoEventsByBuilding schema.SingleBuildingEvents[schema.SectionWithTime]

	// find and parse matching date
	err := mazevoCollection.FindOne(ctx, bson.M{"date": date}).Decode(&mazevoEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			mazevoEvents.Date = date
			mazevoEvents.Buildings = []schema.SingleBuildingEvents[schema.SectionWithTime]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	// filter by the specified building
	for _, b := range mazevoEvents.Buildings {
		if b.Building == building {
			mazevoEventsByBuilding = b
			break
		}
	}

	// if no building found return an error
	if mazevoEventsByBuilding.Building == "" {
		respond(c, http.StatusNotFound, "error", "No events found for the specified building")
		return
	}

	respond(c, http.StatusOK, "success", mazevoEventsByBuilding)
}

// @Id				mazevoEventsByRoom
// @Router			/mazevo/{date}/{building}/{room} [get]
// @Tags			Events
// @Description	"Returns all sections with MazevoEvent meetings on the specified date in the specified building and room"
// @Produce		json
// @Param			date		path		string																true	"date (ISO format) to retrieve mazevo events"
// @Param			building	path		string																true	"building abbreviation of the event location"
// @Param			room		path		string																true	"room number"
// @Success		200			{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.MazevoEvent]]	"All MazevoEvents sections with meetings on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]											"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]
func MazevoEventsByRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")
	room := c.Param("room")

	var mazevoEvents schema.MultiBuildingEvents[schema.SectionWithTime]
	var mazevoEventsByRoom schema.RoomEvents[schema.SectionWithTime]

	// find and parse matching date
	err := mazevoCollection.FindOne(ctx, bson.M{"date": date}).Decode(&mazevoEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			mazevoEvents.Date = date
			mazevoEvents.Buildings = []schema.SingleBuildingEvents[schema.SectionWithTime]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	// filter for the specified building and room
	for _, b := range mazevoEvents.Buildings {
		if b.Building == building {
			for _, r := range b.Rooms {
				if r.Room == room {
					mazevoEventsByRoom = r
					break
				}
			}
			break
		}
	}

	if mazevoEventsByRoom.Room == "" {
		respond(c, http.StatusNotFound, "error", "No events found for that specific building and room")
	}

	respond(c, http.StatusOK, "success", mazevoEventsByRoom)
}
