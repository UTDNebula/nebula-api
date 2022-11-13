package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Preflight(c *gin.Context) {
	c.Header("Access-Contol-Allow-Origin", "*")
	c.Header("Access-Contol-Allow-Headers", "x-api-key, Accept")
	c.JSON(http.StatusOK, struct{}{})
}
