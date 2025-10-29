package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

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

	date := c.Param("date")
	building := c.Param("building")

	var events schema.MultiBuildingEvents[schema.SectionWithTime]
	var eventsByBuilding schema.SingleBuildingEvents[schema.SectionWithTime]

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

	// filter for the specified building
	for _, b := range events.Buildings {
		if b.Building == building {
			eventsByBuilding = b
			break
		}
	}

	// If no building is found, return an error
	if eventsByBuilding.Building == "" {
		c.JSON(http.StatusNotFound, schema.APIResponse[string]{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "No events found for the specified building",
		})
		return
	}

	respond(c, http.StatusOK, "success", eventsByBuilding)
}

// @Id				eventsByRoom
// @Router			/events/{date}/{building}/{room} [get]
// @Description	"Returns all sections with meetings on the specified date in the specified building and room"
// @Produce		json
// @Param			date		path		string											true	"ISO date of the set of events to get"
// @Param			building	path		string											true	"building abbreviation of the event location"
// @Param			room		path		string											true	"room number"
// @Success		200			{object}	schema.APIResponse[[]schema.SectionWithTime]	"All sections with meetings on the specified date in the specified building and room"
// @Failure		500			{object}	schema.APIResponse[string]						"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]						"A string describing the error"
func EventsByRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")
	room := c.Param("room")

	var events schema.MultiBuildingEvents[schema.SectionWithTime]
	var foundSections []schema.SectionWithTime

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

	// filter for the specified building and room
	for _, b := range events.Buildings {
		if b.Building == building {
			for _, r := range b.Rooms {
				if r.Room == room {
					foundSections = r.Events
					break
				}
			}
			break
		}
	}

	if len(foundSections) == 0 {
		c.JSON(http.StatusNotFound, schema.APIResponse[string]{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "No events found for the specified building and room",
		})
		return
	}

	respond(c, http.StatusOK, "success", foundSections)
}

// @Id				sectionsByRoomDetailed
// @Router			/events/{date}/{building}/{room}/sections [get]
// @Description	"Returns full section objects with meetings on the specified date in the specified building and room"
// @Produce		json
// @Param			date		path		string									true	"ISO date of the set of events to get"
// @Param			building	path		string									true	"building abbreviation of the event location"
// @Param			room		path		string									true	"room number"
// @Success		200			{object}	schema.APIResponse[[]schema.Section]	"Full section objects with meetings on the specified date in the specified building and room"
// @Failure		500			{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]				"A string describing the error"
func SectionsByRoomDetailed(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")
	room := c.Param("room")

	var events schema.MultiBuildingEvents[schema.SectionWithTime]

	// Step 1: Find events for the specified date
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

	// Step 2: Extract section IDs for the specified building and room
	var sectionIDs []primitive.ObjectID
	for _, b := range events.Buildings {
		if b.Building == building {
			for _, r := range b.Rooms {
				if r.Room == room {
					for _, event := range r.Events {
						sectionIDs = append(sectionIDs, event.Section)
					}
					break
				}
			}
			break
		}
	}

	if len(sectionIDs) == 0 {
		c.JSON(http.StatusNotFound, schema.APIResponse[string]{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "No sections found for the specified building and room",
		})
		return
	}

	// Step 3: Fetch full section objects from the sections collection
	sectionsCollection := configs.GetCollection("sections")
	cursor, err := sectionsCollection.Find(ctx, bson.M{"_id": bson.M{"$in": sectionIDs}})
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	var sections []schema.Section
	if err = cursor.All(ctx, &sections); err != nil {
		respondWithInternalError(c, err)
		return
	}

	if len(sections) == 0 {
		c.JSON(http.StatusNotFound, schema.APIResponse[string]{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "No section details found",
		})
		return
	}

	respond(c, http.StatusOK, "success", sections)
}
