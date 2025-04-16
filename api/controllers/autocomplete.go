package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"
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
// @Success		200	{array}	schema.Autocomplete	"An aggregation of courses for use in generating autocomplete DAGs"
func AutocompleteDAG(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var autocompleteDAG []schema.Autocomplete

	cursor, err := DAGCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &autocompleteDAG); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, responses.AutocompleteResponse{Status: http.StatusOK, Message: "success", Data: autocompleteDAG})
}
