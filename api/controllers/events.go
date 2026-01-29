package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"
	"strings" // adding missing import

	"github.com/UTDNebula/nebula-api/api/configs"

	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventsCollection *mongo.Collection = configs.GetCollection("events")

// @Id				events
// @Router			/events/{date} [get]
// @Tags			Events
// @Description	"Returns all sections with meetings on the specified date"
// @Produce		json
// @Param			date	path		string																	true	"ISO date of the set of events to get"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.SectionWithTime]]	"All sections with meetings on the specified date"
// @Failure		500		{object}	schema.APIResponse[string]												"A string describing the error"
func Events(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	date := c.Param("date")

	var events schema.MultiBuildingEvents[schema.SectionWithTime]

	defer cancel()

	// find and parse matching date
	err := eventsCollection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			events.Date = date
			events.Buildings = []schema.SingleBuildingEvents[schema.SectionWithTime]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}
	respond(c, http.StatusOK, "success", events)
}

// @Id				eventsByBuilding
// @Router			/events/{date}/{building} [get]
// @Tags			Events
// @Description	"Returns all sections with meetings on the specified date in the specified building"
// @Produce		json
// @Param			date		path		string																	true	"ISO date of the set of events to get"
// @Param			building	path		string																	true	"building abbreviation of event locations"
// @Success		200			{object}	schema.APIResponse[schema.SingleBuildingEvents[schema.SectionWithTime]]	"All sections with meetings on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]												"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]												"A string describing the error"
func EventsByBuilding(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")

	var events schema.MultiBuildingEvents[schema.SectionWithTime]
	var eventsByBuilding schema.SingleBuildingEvents[schema.SectionWithTime]


	// find and parse matching date
	err := eventsCollection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			events.Date = date
			events.Buildings = []schema.SingleBuildingEvents[schema.SectionWithTime]{}
		} else {
			respondWithInternalError(c, err)
			return
		}
	}
	
	// case insensitive filter after data is retrieved
	for _, b := range events.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
			eventsByBuilding = b
			break
		}
	}
	
	// if no building is found, return an err with suggestion
	if eventsByBuilding.Building == "" {
		maxBuildings := min(len(events.Buildings), 10)
		available := make([]string, 0, maxBuildings)
		for i := 0; i < maxBuildings; i++ {
			available = append(available, strings.TrimSpace(events.Buildings[i].Building))
		}
		respond(c, http.StatusNotFound, "error", "Building not found. Available: "+strings.Join(available, ", "))
		return
	}
	respond(c, http.StatusOK, "success", eventsByBuilding)
}

// @Id				eventsByRoom
// @Router			/events/{date}/{building}/{room} [get]
// @Tags			Events
// @Description	"Returns all sections with meetings on the specified date in the specified building and room"
// @Produce		json
// @Param			date		path		string															true	"ISO date of the set of events to get"
// @Param			building	path		string															true	"building abbreviation of the event location"
// @Param			room		path		string															true	"room number"
// @Success		200			{object}	schema.APIResponse[schema.RoomEvents[schema.SectionWithTime]]	"All sections with meetings on the specified date in the specified building and room"
// @Failure		500			{object}	schema.APIResponse[string]										"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]										"A string describing the error"
func EventsByRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := strings.TrimSpace(c.Param("building"))
	room := strings.TrimSpace(c.Param("room"))

	var events schema.MultiBuildingEvents[schema.SectionWithTime]
	var eventsByRoom schema.RoomEvents[schema.SectionWithTime]

	err := eventsCollection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		}
		respondWithInternalError(c, err)
		return
	}

	// 3. Updated to use Case-Insensitive matching for Building
	var matchedBuilding *schema.SingleBuildingEvents[schema.SectionWithTime]
	for _, b := range events.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
			matchedBuilding = &b
			break
		}
	}

	if matchedBuilding == nil {
		respond(c, http.StatusNotFound, "error", "Building not found")
		return
	}

	// 4. Updated to use Case-Insensitive matching for Room
	for _, r := range matchedBuilding.Rooms {
		if strings.EqualFold(strings.TrimSpace(r.Room), room) {
			eventsByRoom = r
			break
		}
	}

	if eventsByRoom.Room == "" {
		maxRooms := min(len(matchedBuilding.Rooms), 20)
		available := make([]string, 0, maxRooms)
		for i := 0; i < maxRooms; i++ {
			available = append(available, strings.TrimSpace(matchedBuilding.Rooms[i].Room))
		}
		respond(c, http.StatusNotFound, "error", "Room not found. Available in this building: "+strings.Join(available, ", "))
		return
	}

	respond(c, http.StatusOK, "success", eventsByRoom)
}

// @Id				sectionsByRoomDetailed
// @Router			/events/{date}/{building}/{room}/sections [get]
// @Tags			Events
// @Description	"Returns full section objects with meetings on the specified date in the specified building and room"
// @Produce		json
// @Param			date		path		string													true	"ISO date of the set of events to get"
// @Param			building	path		string													true	"building abbreviation of the event location"
// @Param			room		path		string													true	"room number"
// @Success		200			{object}	schema.APIResponse[schema.RoomEvents[schema.Section]]	"Full section objects with meetings on the specified date in the specified building and room"
// @Failure		500			{object}	schema.APIResponse[string]								"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]								"A string describing the error"
func SectionsByRoomDetailed(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := strings.TrimSpace(c.Param("building"))
	room := strings.TrimSpace(c.Param("room"))

	var events schema.MultiBuildingEvents[schema.SectionWithTime]
	var sectionsByRoom schema.RoomEvents[schema.Section]

	// Find events for the specified date
	err := eventsCollection.FindOne(ctx, bson.M{"date": date}).Decode(&events)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		}
		respondWithInternalError(c, err)
		return
	}

	// Extract section IDs for the specified building and room using case-insensitive matching
	var sectionIDs []primitive.ObjectID
	buildingFound := false
	for _, b := range events.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
			buildingFound = true
			for _, r := range b.Rooms {
				if strings.EqualFold(strings.TrimSpace(r.Room), room) {
					sectionsByRoom.Room = r.Room
					for _, event := range r.Events {
						sectionIDs = append(sectionIDs, event.Section)
					}
					break
				}
			}
			break
		}
	}

	if !buildingFound {
		respond(c, http.StatusNotFound, "error", "Building not found")
		return
	}

	if len(sectionIDs) == 0 {
		respond(c, http.StatusNotFound, "error", "No sections found for the specified room")
		return
	}

	// Fetch full section objects from the sections collection using the extracted IDs
	cursor, err := sectionCollection.Find(ctx, bson.M{"_id": bson.M{"$in": sectionIDs}})
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &sectionsByRoom.Events); err != nil {
		respondWithInternalError(c, err)
		return
	}

	if len(sectionsByRoom.Events) == 0 {
		respond(c, http.StatusNotFound, "error", "No section details found")
		return
	}

	respond(c, http.StatusOK, "success", sectionsByRoom)
}
