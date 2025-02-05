package controllers

import (
	"github.com/gin-gonic/gin"
)

// @Id rooms
// @Router /rooms [get]
// @Description "Returns all classrooms being used in the current and futures semesters"
// @Produce json
// @Success 200 {array} schema.BuildingRooms "All classrooms being used in the current and futures semesters"
func Rooms(c *gin.Context) {

}
