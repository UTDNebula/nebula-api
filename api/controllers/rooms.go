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

var buildingCollection *mongo.Collection = configs.GetCollection("rooms")

// @Id				rooms
// @Router			/rooms [get]
// @Description	"Returns all schedulable rooms being used in the current and futures semesters from CourseBook, Astra, and Mazevo"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]schema.BuildingRooms]	"All schedulable rooms being used in the current and futures semesters from CourseBook, Astra, and Mazevo"
// @Failure		500	{object}	schema.APIResponse[string]					"A string describing the error"
func Rooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var buildingRooms []schema.BuildingRooms // buildings and rooms to be returned

	//cursor is the pointer for the returned documents
	cursor, err := buildingCollection.Find(ctx, bson.M{})
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// use cursor to fill rooms slice to return
	err = cursor.All(ctx, &buildingRooms)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	// serialize RoomsResponse struct data into JSON format
	respond(c, http.StatusOK, "success", buildingRooms)
}
