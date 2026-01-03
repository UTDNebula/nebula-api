package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/api/configs"

	"github.com/UTDNebula/nebula-api/api/schema"
)

var cometCalendarCollection *mongo.Collection = configs.GetCollection("cometCalendar")

// @Id				CometCalendarEvents
// @Router			/calendar/{date} [get]
// @Tags			Events
// @Description	"Returns CometCalendarEvent based on the input date"
// @Produce		json
// @Param			date	path		string																		true	"date (ISO format) to retrieve comet calendar events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.CometCalendarEvent]]	"All CometCalendarEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]													"A string describing the error"
func CometCalendarEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	var cometCalendarEvents schema.MultiBuildingEvents[schema.CometCalendarEvent]

	// Find comet calendar event given date
	err := cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			cometCalendarEvents.Date = date
			cometCalendarEvents.Buildings = []schema.SingleBuildingEvents[schema.CometCalendarEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	respond(c, http.StatusOK, "success", cometCalendarEvents)
}

// @Id				CometCalendarEventsByBuilding
// @Router			/calendar/{date}/{building} [get]
// @Tags			Events
// @Description	"Returns CometCalendarEvent based on the input date and building name"
// @Produce		json
// @Param			date		path		string																			true	"date (ISO format) to retrieve comet calendar events"
// @Param			building	path		string																			true	"building abbreviation of event locations"
// @Success		200			{object}	schema.APIResponse[schema.SingleBuildingEvents[schema.CometCalendarEvent]]	"All events on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]													"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]													"A string describing the error"
func CometCalendarEventsByBuilding(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")

	var cometCalendarEvents schema.MultiBuildingEvents[schema.CometCalendarEvent]
	var cometCalendarEventsByBuilding schema.SingleBuildingEvents[schema.CometCalendarEvent]

	// Find comet calendar event given date
	err := cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			cometCalendarEvents.Date = date
			cometCalendarEvents.Buildings = []schema.SingleBuildingEvents[schema.CometCalendarEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	//parse response for requested building
	for _, b := range cometCalendarEvents.Buildings {
		if b.Building == building {
			cometCalendarEventsByBuilding = b
			break
		}
	}

	if cometCalendarEventsByBuilding.Building == "" {
		respond(c, http.StatusNotFound, "error", "No events found for the specified building")
		return
	}

	respond(c, http.StatusOK, "success", cometCalendarEventsByBuilding)
}

// @Id				CometCalendarEventsByBuildingAndRoom
// @Router			/calendar/{date}/{building}/{room} [get]
// @Tags			Events
// @Description	"Returns CometCalendarEvent based on the input date building name and room number"
// @Produce		json
// @Param			date		path		string																		true	"date (ISO format) to retrieve comet calendar events"
// @Param			building	path		string																		true	"building abbreviation of event locations"
// @Param			room		path		string																		true	"room number for event"
// @Success		200			{object}	schema.APIResponse[schema.RoomEvents[schema.CometCalendarEvent]]	"All events on the specified date in the specified building and room"
// @Failure		500			{object}	schema.APIResponse[string]												"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]												"A string describing the error"
func CometCalendarEventsByBuildingAndRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := strings.TrimSpace(c.Param("building"))
	room := strings.TrimSpace(c.Param("room"))

	var cometCalendarEvents schema.MultiBuildingEvents[schema.CometCalendarEvent]
	var roomEvents schema.RoomEvents[schema.CometCalendarEvent]

	// Find comet calendar event given date
	err := cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			cometCalendarEvents.Date = date
			cometCalendarEvents.Buildings = []schema.SingleBuildingEvents[schema.CometCalendarEvent]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	//parse response for requested building and room (case-insensitive matching)
	for _, b := range cometCalendarEvents.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
			for _, r := range b.Rooms {
				if strings.EqualFold(strings.TrimSpace(r.Room), room) {
					roomEvents = r
					break
				}
			}
			break
		}
	}

	if roomEvents.Room == "" {
		respond(c, http.StatusNotFound, "error", "No events found for the specified building or room")
		return
	}

	respond(c, http.StatusOK, "success", roomEvents)
}
