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

var lettersCollection *mongo.Collection = configs.GetCollection("letters")

// @Id				letters
// @Router			/games/letters/{date} [get]
// @Description	"Returns letters for the specified date"
// @Produce		json
// @Param			date	path		string								true	"ISO date of the letters to get"
// @Success		200		{object}	schema.APIResponse[schema.Letters]	"Letters for the specified date"
// @Failure		500		{object}	schema.APIResponse[string]			"A string describing the error"
func Letters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	date := c.Param("date")

	var letters schema.Letters

	defer cancel()

	// find and parse matching date
	err := lettersCollection.FindOne(ctx, bson.M{"date": date}).Decode(&letters)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", letters)
}
