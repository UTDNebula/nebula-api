package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"

	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventsCollection *mongo.Collection = configs.GetCollection("events")

// @Id				events
// @Router			/events/{date} [get]
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
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", events)
}

// @Id				eventsByBuilding
// @Router			/events/{date}/{building} [get]
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
		respondWithInternalError(c, err)
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
		c.JSON(http.StatusNotFound, schema.APIResponse[string]{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "No events found for the specified building",
		})
		return
	}

	respond(c, http.StatusOK, "success", eventsByBuilding)
}

// @Id				eventsByRoomSection
// @Router			/events/{date}/{building}/{room}/{sections} [get]
// @Description	"Returns all sections (now, as a section object) with meetings on the specified date in the specified building"
// @Produce		json
// @Param			date		path		string																	true	"ISO date of the set of events to get"
// @Param			building	path		string	                                                                true    "building abbreviation of the event location"
// @Param           room        path        string                                                                true    "section ID (ObjectID)"
// @Success		200			{object}	schema.APIResponse[schema.SingleRoomEvents[schema.SectionWithTime]]	       	"All sections with meetings on the specified date in the specified building"
// @Failure		500			{object}	schema.APIResponse[string]												"A string describing the error"
// @Failure		404			{object}	schema.APIResponse[string]												"A string describing the error"
func EventsByRoomSection(c *gin.Context) {
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
		respondWithInternalError(c, err)
		return
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

	// return all sections in the room
	respond(c, http.StatusOK, "success", foundSections)
}
