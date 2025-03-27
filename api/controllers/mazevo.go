package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"
)

var mazevoCollection *mongo.Collection = configs.GetCollection("mazevo")

// @Id MazevoEvents
// @Router /mazevo/{date} [get]
// @Description "Returns MazevoEvent based on the input date"
// @Produce json
// @Param date path string true "date (ISO format) to retrieve mazevo events"
// @Success 200 {array} schema.MultiBuildingEvents[schema.MazevoEvent] "All MazevoEvents with events on the inputted date"
func MazevoEvents(c *gin.Context) {
	//gin context has info about request and allows converting to json format

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //make resources available if function returns before timeout

	date := c.Param("date") //input date

	var mazevo_events schema.MultiBuildingEvents[schema.MazevoEvent] //stores mazevo events occuring on input date

	//find mazevo event for input date
	err := mazevoCollection.FindOne(ctx, bson.M{"date": date}).Decode(&mazevo_events)

	//error
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}
	//no error
	c.JSON(http.StatusOK, responses.MultiBuildingEventsResponse[schema.MazevoEvent]{Status: http.StatusOK, Message: "success", Data: mazevo_events})
}
