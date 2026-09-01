package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/rest/configs"

	"github.com/UTDNebula/nebula-api/rest/schema"
)

var astraTracer = otel.Tracer("astra-controller")

// @Id				AstraEvents
// @Router			/astra/{date} [get]
// @Tags			Events
// @Description	"Returns AstraEvent based on the input date"
// @Produce		json
// @Param			date	path		string																true	"date (ISO format) to retrieve astra events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.AstraEvent]]	"All AstraEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]											"A string describing the error"
func AstraEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	ctx, span := astraTracer.Start(ctx, "astra.events")
	defer span.End()

	date := c.Param("date")

	var astra_events schema.MultiBuildingEvents[schema.AstraEvent]

	// Find astra event given date
	err := configs.GetCollection("astra").FindOne(ctx, bson.M{
		"date": date,
	}).Decode(&astra_events)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	ctx, span := astraTracer.Start(ctx, "astra.events.by.building")
	defer span.End()

	date := c.Param("date")
	building := strings.TrimSpace(c.Param("building")) // trimming the input

	var astraEvents schema.MultiBuildingEvents[schema.AstraEvent]
	var astraEventsByBuilding schema.SingleBuildingEvents[schema.AstraEvent]

	// Find astra event given date
	err := configs.GetCollection("astra").
		FindOne(ctx, bson.M{"date": date}).
		Decode(&astraEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		}
		respondWithInternalError(c, err)
		return
	}

	// case insensitive matching
	for _, b := range astraEvents.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
			astraEventsByBuilding = b
			break
		}
	}

	if astraEventsByBuilding.Building == "" {
		// provide suggestion if not found
		var available []string
		for i := range len(astraEvents.Buildings) {
			available = append(available, strings.TrimSpace(astraEvents.Buildings[i].Building))
		}
		respond(c, http.StatusNotFound, "error", map[string][]string{
			"available": available,
		})
		return
	}

	respond(c, http.StatusOK, "success", astraEventsByBuilding)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	ctx, span := astraTracer.Start(ctx, "astra.events.by.building.and.room")
	defer span.End()

	date := c.Param("date")
	building := strings.TrimSpace(c.Param("building"))
	room := strings.TrimSpace(c.Param("room"))

	var astraEvents schema.MultiBuildingEvents[schema.AstraEvent]
	var roomEvents schema.RoomEvents[schema.AstraEvent]

	// Find astra event given date
	err := configs.GetCollection("astra").
		FindOne(ctx, bson.M{"date": date}).
		Decode(&astraEvents)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No events found for the specified date")
			return
		}
		respondWithInternalError(c, err)
		return
	}

	// matching buildings case-insensitively
	var matchedBuilding *schema.SingleBuildingEvents[schema.AstraEvent]
	for _, b := range astraEvents.Buildings {
		if strings.EqualFold(strings.TrimSpace(b.Building), building) {
			matchedBuilding = &b
			break
		}
	}

	if matchedBuilding == nil {
		respond(c, http.StatusNotFound, "error", "Building not found")
		return
	}

	// match room case-insesitively
	for _, r := range matchedBuilding.Rooms {
		if strings.EqualFold(strings.TrimSpace(r.Room), room) {
			roomEvents = r
			break
		}
	}

	if roomEvents.Room == "" {
		var available []string
		for i := range len(matchedBuilding.Rooms) {
			available = append(available, strings.TrimSpace(matchedBuilding.Rooms[i].Room))
		}
		respond(c, http.StatusNotFound, "error", map[string][]string{
			"available": available,
		})
		return
	}

	respond(c, http.StatusOK, "success", roomEvents)
}
