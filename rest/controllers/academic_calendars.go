package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/rest/configs"

	"github.com/UTDNebula/nebula-api/rest/schema"
)

var academicCalendarsCollection *mongo.Collection = configs.GetCollection("academicCalendars")

// AcademicCalendarsById returns an academic calendar by its string ID.
//
//	@Id					academicCalendarsById
//	@Router				/academicCalendars/{id} [get]
//	@Tags				Academic Calendars
//	@Description		Returns the academic calendar for the given term ID.
//	@Produce			json
//	@Param			id	path		string										true							"Academic term ID, for example 26F"
//	@Success		200	{object}	schema.APIResponse[schema.AcademicCalendar]	"An academic calendar"
//	@Failure		404	{object}	schema.APIResponse[string]					"No matching academic calendar"
//	@Failure		500	{object}	schema.APIResponse[string]					"An internal server error"
func AcademicCalendarsById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var academicCalendar schema.AcademicCalendar

	// Academic calendar IDs are strings such as "26F", not MongoDB ObjectIDs.
	id := c.Param("id")
	query := bson.M{"_id": id}

	// Find and parse matching academic calendar
	err := academicCalendarsCollection.FindOne(ctx, query).Decode(&academicCalendar)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No academic calendars with the given ID")
		} else {
			respondWithInternalError(c, err)
		}
		return
	}

	//return result
	respond(c, http.StatusOK, "success", academicCalendar)
}

// AcademicCalendarsCurrent returns an academic calendar marked as current.
//
//	@Id					academicCalendarsCurrent
//	@Router				/academicCalendars/current [get]
//	@Tags				Academic Calendars
//	@Description		Returns one academic calendar whose timeline is current.
//	@Produce			json
//	@Success		200	{object}	schema.APIResponse[schema.AcademicCalendar]	"The current academic calendar"
//	@Failure		404	{object}	schema.APIResponse[string]					"No current academic calendar"
//	@Failure		500	{object}	schema.APIResponse[string]					"An internal server error"
func AcademicCalendarsCurrent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var academicCalendar schema.AcademicCalendar

	query := bson.M{"timeline": "current"}

	//find and parse matching academic calendar
	err := academicCalendarsCollection.FindOne(ctx, query).Decode(&academicCalendar)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No current academic calendars")
		} else {
			respondWithInternalError(c, err)
		}
		return
	}

	respond(c, http.StatusOK, "success", academicCalendar)

}
