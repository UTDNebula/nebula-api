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

var DAGCollection *mongo.Collection = configs.GetCollection("DAG")

// @Id				autocompleteDAG
// @Router			/autocomplete/dag [get]
// @Description	"Returns an aggregation of courses for use in generating autocomplete DAGs"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]schema.Autocomplete]	"An aggregation of courses for use in generating autocomplete DAGs"
// @Failure		500	{object}	schema.APIResponse[string]					"A string describing the error"
func AutocompleteDAG(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var autocompleteDAG []schema.Autocomplete

	cursor, err := DAGCollection.Find(ctx, bson.M{})

	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Retrieve and parse all valid documents from view
	if err = cursor.All(ctx, &autocompleteDAG); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", autocompleteDAG)
}
