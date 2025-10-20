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

var astraCollection *mongo.Collection = configs.GetCollection("astra")

// @Id				AstraEvents
// @Router			/astra/{date} [get]
// @Tags			Events
// @Description	"Returns AstraEvent based on the input date"
// @Produce		json
// @Param			date	path		string																true	"date (ISO format) to retrieve astra events"
// @Success		200		{object}	schema.APIResponse[schema.MultiBuildingEvents[schema.AstraEvent]]	"All AstraEvents with events on the inputted date"
// @Failure		500		{object}	schema.APIResponse[string]											"A string describing the error"
func AstraEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	date := c.Param("date")

	var astra_events schema.MultiBuildingEvents[schema.AstraEvent]

	// Find astra event given date
	err := astraCollection.FindOne(ctx, bson.M{"date": date}).Decode(&astra_events)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", astra_events)
}
