package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Preflight(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}
