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

func AcademicCalendarsById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var academicCalendar schema.AcademicCalendar

	//parse string id from id parameter
	id := c.Param("id")
	query := bson.M{"_id": id}

	//find and parse matching academic calendar
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

func AcadenemicCalendarsCurrent(c *gin.Context) {
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
