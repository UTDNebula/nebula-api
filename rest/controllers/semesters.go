package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const testfilename = "semesters-data.json"

func Semesters(c *gin.Context) {
	terms := []string{
		"26u",
		"26s",
		"26f",
		"25u",
		"25s",
		"25f",
		"24u",
		"24s",
		"24f",
		"23u",
		"23s",
		"23f",
		"22u",
		"22s",
		"22f",
		"21u",
		"21s",
		"21f",
		"20u",
		"20s",
		"20f",
		"19u",
		"19s",
		"19f",
		"18u",
		"18s",
	}

	jsonData, err := json.Marshal(terms)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	err = os.WriteFile(testfilename, jsonData, 0644)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	fileBytes, err := os.ReadFile(testfilename)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", json.RawMessage(fileBytes))
}
