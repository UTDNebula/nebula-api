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

// @Id				Budget
// @Router			/budget/{year} [get]
// @Tags			Other
// @Description	"Returns discal year Budget based on the input year"
// @Produce		json
// @Param			year	path		string								true	"year to retrieve budget for"
// @Success		200		{object}	schema.APIResponse[schema.Budget]	"Single budget from the given fiscal year"
// @Failure		500		{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		404		{object}	schema.APIResponse[string]			"A string describing the error"
func Budget(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	year := c.Param("year")

	var budget schema.Budget

	// Find budget given date
	err := configs.GetCollection("budgets").FindOne(ctx, bson.M{"_id": year}).Decode(&budget)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respond(c, http.StatusNotFound, "error", "No budgets found for the specified fiscal year")
			return
		} else {
			respondWithInternalError(c, err)
			return
		}
	}

	respond(c, http.StatusOK, "success", budget)
}
