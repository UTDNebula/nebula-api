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

var buildingCollection *mongo.Collection = configs.GetCollection("rooms")

// @Id rooms
// @Router /rooms [get]
// @Description "Returns all schedulable rooms being used in the current and futures semesters from CourseBook, Astra, and Mazevo"
// @Produce json
// @Success 200 {array} schema.BuildingRooms "All schedulable rooms being used in the current and futures semesters from CourseBook, Astra, and Mazevo"
func Rooms(c *gin.Context) {
	//gin context has info about request and allows converting to json format

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //make resources available if Rooms returns before timeout

	var buildingRooms []schema.BuildingRooms //buildings and rooms to be returned

	//cursor is the pointer for the returned documents
	cursor, err := buildingCollection.Find(ctx, bson.M{})

	//if there is an error
	if err != nil {
		//serialize ErrorResponse struct data into JSON format
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	//use cursor to fill rooms slice to return
	err = cursor.All(ctx, &buildingRooms)

	//if there is an error filling rooms slice
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	//serialize RoomsResponse struct data into JSON format
	c.JSON(http.StatusOK, responses.RoomsResponse{Status: http.StatusOK, Message: "success", Data: buildingRooms})
}
