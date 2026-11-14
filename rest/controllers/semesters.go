package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/UTDNebula/nebula-api/rest/configs"
	_ "github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
)

// @Id				semesters
// @Router			/semesters [get]
// @Tags			Other
// @Description	"Returns a list of all available terms in CourseBook, sorted by recency"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]string]	"A string list of all available terms, lowercased"
// @Failure		500	{object}	schema.APIResponse[string]		"A string describing the error"
func Semesters(c *gin.Context) {
	fileBytes, err := os.ReadFile(configs.SemestersJSONFilename)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", json.RawMessage(fileBytes))
}
