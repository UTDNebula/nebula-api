package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/api/configs"

	"github.com/UTDNebula/nebula-api/api/schema"
)

var mazevoCollection *mongo.Collection = configs.GetCollection("mazevo")

// @Id				MazevoEvents
// @Router			events/mazevo/{date} [get]
// @Description	"Returns MazevoEvent based on the input date"
// @Produce		json
// @Param			date	path		string																true	"date (ISO format) to retrieve mazevo events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.MazevoEvent]]	"All MazevoEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]											"A string describing the error"
func MazevoEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	var mazevo_events schema.MultiBuildingEvents[schema.MazevoEvent]

	// Find mazevo event for input date
	err := mazevoCollection.FindOne(ctx, bson.M{"date": date}).Decode(&mazevo_events)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", mazevo_events)
}
