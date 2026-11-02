package events

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/rest/responses"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
)

type mazevoController struct {
	repository shared.EventRepository[schema.MazevoEvent]
}

// @Id				MazevoEvents
// @Router			/mazevo/{date} [get]
// @Tags			Events
// @Description	"Returns MazevoEvent based on the input date"
// @Produce		json
// @Param			date	path		string																true	"date (ISO format) to retrieve mazevo events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.MazevoEvent]]	"All MazevoEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]											"A string describing the error"
func (h mazevoController) MazevoEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	// Find mazevo event for input date
	mazevoEvents, err := h.repository.GetEventsByDate(ctx, date)
	if err != nil {
		if errors.Is(err, shared.ErrEventsNotFound) {
			mazevoEvents.Date = date
			mazevoEvents.Buildings = []schema.SingleBuildingEvents[schema.MazevoEvent]{}
		} else {
			responses.RespondWithInternalError(c, err)
			return
		}
	}

	responses.Respond(c, http.StatusOK, "success", mazevoEvents)
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
func (h mazevoController) MazevoEventsByBuilding(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")

	var mazevoEventsByBuilding schema.SingleBuildingEvents[schema.MazevoEvent]

	// find and parse matching date
	mazevoEvents, err := h.repository.GetEventsByDate(ctx, date)
	if err != nil {
		if errors.Is(err, shared.ErrEventsNotFound) {
			mazevoEvents.Date = date
			mazevoEvents.Buildings = []schema.SingleBuildingEvents[schema.MazevoEvent]{}
		} else {
			responses.RespondWithInternalError(c, err)
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
		responses.Respond(c, http.StatusNotFound, "error", "No events found for the specified building")
		return
	}

	responses.Respond(c, http.StatusOK, "success", mazevoEventsByBuilding)
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
func (h mazevoController) MazevoEventsByRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	date := c.Param("date")
	building := c.Param("building")
	room := c.Param("room")

	var mazevoEventsByRoom schema.RoomEvents[schema.MazevoEvent]

	// find and parse matching date
	mazevoEvents, err := h.repository.GetEventsByDate(ctx, date)
	if err != nil {
		if errors.Is(err, shared.ErrEventsNotFound) {
			mazevoEvents.Date = date
			mazevoEvents.Buildings = []schema.SingleBuildingEvents[schema.MazevoEvent]{}
		} else {
			responses.RespondWithInternalError(c, err)
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
		responses.Respond(c, http.StatusNotFound, "error", "No events found for that specific building and room")
		return
	}

	responses.Respond(c, http.StatusOK, "success", mazevoEventsByRoom)
}
