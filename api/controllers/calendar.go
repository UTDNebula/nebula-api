package controllers

import (
	"context"
	"errors"
	"log"
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
// @Param			date	path		string															true	"date (ISO format) to retrieve comet calendar events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.Event]]	"All CometCalendarEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]										"A string describing the error"
func CometCalendarEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	var cometCalendarEvents schema.MultiBuildingEvents[schema.Event]

	// Find comet calendar event given date
	err := cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			cometCalendarEvents.Date = date
			cometCalendarEvents.Buildings = []schema.SingleBuildingEvents[schema.Event]{}
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
// @Param			date		path		string															true	"date (ISO format) to retrieve comet calendar events"
// @Param			building	path		string															true	"building abbreviation of event locations"
// @Success		200			{object}	schema.APIResponse[schema.SingleBuildingEvents[schema.Event]]	"All events on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]										"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]										"A string describing the error"
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

	var cometCalendarEvents schema.MultiBuildingEvents[schema.Event]
	var cometCalendarEventsByBuilding schema.SingleBuildingEvents[schema.Event]

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
// @Param			date		path		string												true	"date (ISO format) to retrieve comet calendar events"
// @Param			building	path		string												true	"building abbreviation of event locations"
// @Param			room		path		string												true	"room number for event"
// @Success		200			{object}	schema.APIResponse[schema.RoomEvents[schema.Event]]	"All events on the specified date in the specified building and room"
// @Failure		500			{object}	schema.APIResponse[string]							"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]							"A string describing the error"
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

	var cometCalendarEvents schema.MultiBuildingEvents[schema.Event]
	var roomEvents schema.RoomEvents[schema.Event]

	// Find comet calendar event given date
	log.Printf("Querying cometCalendar collection for date: %s", date)
	err = cometCalendarCollection.FindOne(ctx, bson.M{"date": date}).Decode(&cometCalendarEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No documents found for date: %s", date)
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		} else {
			log.Printf("Database error for date %s: %v", date, err)
			respondWithInternalError(c, err)
			return
		}
	}
	log.Printf("Found data for date %s with %d buildings", date, len(cometCalendarEvents.Buildings))

	// Check if any buildings exist
	if len(cometCalendarEvents.Buildings) == 0 {
		respond(c, http.StatusNotFound, "error", "No buildings found for the specified date")
		return
	}

	// Parse response for requested building and room (case-insensitive matching)
	buildingFound := false
	var matchedBuilding *schema.SingleBuildingEvents[schema.Event]

	for _, b := range cometCalendarEvents.Buildings {
		buildingName := strings.TrimSpace(b.Building)

		if strings.EqualFold(buildingName, building) {
			buildingFound = true
			matchedBuilding = &b
			break
		}
	}

	if !buildingFound {
		// Collect available buildings only when needed (limit to first 10 for performance)
		maxBuildings := 10
		if len(cometCalendarEvents.Buildings) < maxBuildings {
			maxBuildings = len(cometCalendarEvents.Buildings)
		}
		availableBuildings := make([]string, 0, maxBuildings)
		for i := 0; i < maxBuildings; i++ {
			availableBuildings = append(availableBuildings, strings.TrimSpace(cometCalendarEvents.Buildings[i].Building))
		}
		buildingList := strings.Join(availableBuildings, ", ")
		if len(cometCalendarEvents.Buildings) > 10 {
			buildingList += " (and more...)"
		}
		respond(c, http.StatusNotFound, "error",
			"No events found for the specified building. Available buildings: "+buildingList)
		return
	}

	// Check if any rooms exist for this building
	if len(matchedBuilding.Rooms) == 0 {
		respond(c, http.StatusNotFound, "error", "No rooms found for the specified building")
		return
	}

	// Look for the room - try exact match first, then case-insensitive
	log.Printf("Searching for room '%s' in building '%s' with %d rooms", room, building, len(matchedBuilding.Rooms))
	for _, r := range matchedBuilding.Rooms {
		roomName := strings.TrimSpace(r.Room)

		// Try exact match first
		if roomName == room {
			log.Printf("Found exact match for room: %s", room)
			roomEvents = r
			break
		}
		// Try case-insensitive match
		if strings.EqualFold(roomName, room) {
			log.Printf("Found case-insensitive match for room: %s (matched: %s)", room, roomName)
			roomEvents = r
			break
		}
	}

	if roomEvents.Room == "" {
		log.Printf("Room '%s' not found in building '%s'", room, building)
		// Collect available rooms only when needed (limit to first 20 for performance)
		maxRooms := 20
		if len(matchedBuilding.Rooms) < maxRooms {
			maxRooms = len(matchedBuilding.Rooms)
		}
		availableRooms := make([]string, 0, maxRooms)
		for i := 0; i < maxRooms; i++ {
			availableRooms = append(availableRooms, strings.TrimSpace(matchedBuilding.Rooms[i].Room))
		}
		roomList := strings.Join(availableRooms, ", ")
		if len(matchedBuilding.Rooms) > 20 {
			roomList += " (and more...)"
		}
		respond(c, http.StatusNotFound, "error",
			"No events found for the specified room. Available rooms: "+roomList)
		return
	}

	log.Printf("Successfully found room events for %s/%s/%s", date, building, room)
	respond(c, http.StatusOK, "success", roomEvents)
}
