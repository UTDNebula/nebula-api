package controllers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
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

	// URL decode the building parameter in case it contains special characters
	// Use PathUnescape for path parameters (not QueryUnescape)
	buildingParam, err := url.PathUnescape(c.Param("building"))
	if err != nil {
		buildingParam = c.Param("building")
	}
	building := strings.TrimSpace(buildingParam)

	var cometCalendarEvents schema.MultiBuildingEvents[schema.CometCalendarEvent]
	var cometCalendarEventsByBuilding schema.SingleBuildingEvents[schema.CometCalendarEvent]

	// Find comet calendar event given date
	err = cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	// Check if any buildings exist
	if len(cometCalendarEvents.Buildings) == 0 {
		respond(c, http.StatusNotFound, "error", "No buildings found for the specified date")
		return
	}

	// Parse response for requested building (case-insensitive matching)
	for _, b := range cometCalendarEvents.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
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

	// URL decode the building and room parameters in case they contain special characters
	// Use PathUnescape for path parameters (not QueryUnescape)
	buildingParam, err := url.PathUnescape(c.Param("building"))
	if err != nil {
		buildingParam = c.Param("building")
	}
	building := strings.TrimSpace(buildingParam)

	roomParam, err := url.PathUnescape(c.Param("room"))
	if err != nil {
		roomParam = c.Param("room")
	}
	room := strings.TrimSpace(roomParam)

	var cometCalendarEvents schema.MultiBuildingEvents[schema.CometCalendarEvent]
	var roomEvents schema.RoomEvents[schema.CometCalendarEvent]

	// Find comet calendar event given date
	err = cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	// Check if any buildings exist
	if len(cometCalendarEvents.Buildings) == 0 {
		respond(c, http.StatusNotFound, "error", "No buildings found for the specified date")
		return
	}

	// Parse response for requested building and room (case-insensitive matching)
	buildingFound := false
	var availableBuildings []string
	var availableRooms []string

	for _, b := range cometCalendarEvents.Buildings {
		buildingName := strings.TrimSpace(b.Building)
		availableBuildings = append(availableBuildings, buildingName)

		if strings.EqualFold(buildingName, building) {
			buildingFound = true
			// Check if any rooms exist for this building
			if len(b.Rooms) == 0 {
				respond(c, http.StatusNotFound, "error", "No rooms found for the specified building")
				return
			}
			// Look for the room - try exact match first, then case-insensitive
			for _, r := range b.Rooms {
				roomName := strings.TrimSpace(r.Room)
				availableRooms = append(availableRooms, roomName)

				// Try exact match first
				if roomName == room {
					roomEvents = r
					break
				}
				// Try case-insensitive match
				if strings.EqualFold(roomName, room) {
					roomEvents = r
					break
				}
			}
			break
		}
	}

	if !buildingFound {
		// Return helpful error with available buildings
		respond(c, http.StatusNotFound, "error",
			"No events found for the specified building. Available buildings: "+strings.Join(availableBuildings, ", "))
		return
	}

	if roomEvents.Room == "" {
		// Return helpful error with available rooms
		respond(c, http.StatusNotFound, "error",
			"No events found for the specified room. Available rooms: "+strings.Join(availableRooms, ", "))
		return
	}

	respond(c, http.StatusOK, "success", roomEvents)
}
