package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var degreesCollection *mongo.Collection = configs.GetCollection("degrees")

// @Id				degrees
// @Router			/degrees [get]
// @Tags			Degrees
// @Description	"Returns paginated academic programs matching the query's string-typed key-value pairs. An areas_of_interest query matches any area in the program's areas_of_interest array."
// @Produce		json
// @Param			offset						query		number											false	"The starting position of the current page of degrees (e.g. For starting at the 17th degree, offset=16)."
// @Param			name						query		string											false	"The name of the academic program"
// @Param			school						query		string											false	"The school offering the academic program"
// @Param			degree_options.level		query		string											false	"The level of one degree option, such as BS or MS"
// @Param			degree_options.public_url	query		string											false	"The public URL of one degree option"
// @Param			degree_options.cip_code		query		string											false	"The Classification of Instructional Programs code of one degree option"
// @Param			areas_of_interest			query		string											false	"One area of interest for the academic program"
// @Success		200							{object}	schema.APIResponse[[]schema.AcademicProgram]	"A list of academic programs"
// @Failure		500							{object}	schema.APIResponse[string]						"A string describing the error"
// @Failure		400							{object}	schema.APIResponse[string]						"A string describing the error"
func Degrees(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	query, err := getQuery[schema.AcademicProgram]("Search", c)
	if err != nil {
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	cursor, err := degreesCollection.Find(ctx, query, optionLimit)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	degrees := make([]schema.AcademicProgram, 0)
	if err = cursor.All(ctx, &degrees); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", degrees)
}
