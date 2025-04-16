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

var astraCollection *mongo.Collection = configs.GetCollection("astra")

// @Id				AstraEvents
// @Router			/astra/{date} [get]
// @Description	"Returns AstraEvent based on the input date"
// @Produce		json
// @Param			date	path	string											true	"date (ISO format) to retrieve astra events"
// @Success		200		{array}	schema.MultiBuildingEvents[schema.AstraEvent]	"All AstraEvents with events on the inputted date"
func AstraEvents(c *gin.Context) {
	//gin context has info about request and allows converting to json format

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //make resources available if function returns before timeout

	date := c.Param("date") //input date

	var astra_events schema.MultiBuildingEvents[schema.AstraEvent] //stores astra events for input date

	//find astra event given date
	err := astraCollection.FindOne(ctx, bson.M{"date": date}).Decode(&astra_events)

	//if there is an error
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}
	//no error
	c.JSON(http.StatusOK, responses.MultiBuildingEventsResponse[schema.AstraEvent]{Status: http.StatusOK, Message: "success", Data: astra_events})
}
