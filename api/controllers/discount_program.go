package controllers

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var discountProgramCollection *mongo.Collection = configs.GetCollection("discountPrograms")

var discountProgramCategories = map[string]string{
	"accommodations":             "Accommodations",
	"auto services":              "Auto Services",
	"clothes, flowers and gifts": "Clothes, Flowers and Gifts",
	"dining":                     "Dining",
	"entertainment":              "Entertainment",
}

var discountProgramCategoryList = []string{
	"Accommodations",
	"Auto Services",
	"Clothes, Flowers and Gifts",
	"Dining",
	"Entertainment",
}

// @Id				discountProgramsSearch
// @Router			/discountPrograms [get]
// @Tags			DiscountPrograms
// @Description	Returns paginated list of discount programs with optional filtering.
// @Description	Supported query params are category, business, address, and discount.
// @Description	The q query param performs keyword search across all fields and must be used alone.
// @Produce		json
// @Param			offset		query		number											false	"The starting position of the current page of discount programs"
// @Param			category	query		string											false	"One of Accommodations, Auto Services, Clothes, Flowers and Gifts, Dining, or Entertainment (case-insensitive)"
// @Param			business	query		string											false	"Keyword search in business"
// @Param			address		query		string											false	"Keyword search in address"
// @Param			discount	query		string											false	"Keyword search in discount"
// @Param			q			query		string											false	"Keyword search across all fields. Must be used alone"
// @Success		200			{object}	schema.APIResponse[[]schema.DiscountProgram]	"A list of discount programs"
// @Failure		400			{object}	schema.APIResponse[string]						"A string describing the error"
// @Failure		500			{object}	schema.APIResponse[string]						"A string describing the error"
func DiscountProgramSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{}

	category := strings.TrimSpace(c.Query("category"))
	business := strings.TrimSpace(c.Query("business"))
	address := strings.TrimSpace(c.Query("address"))
	discount := strings.TrimSpace(c.Query("discount"))
	keyword := strings.TrimSpace(c.Query("q"))

	if keyword != "" {
		if category != "" || business != "" || address != "" || discount != "" {
			respond(
				c,
				http.StatusBadRequest,
				"Invalid query parameters",
				"q must be used alone (except with offset)",
			)
			return
		}

		query["$or"] = []bson.M{
			{"category": regexFilter(keyword)},
			{"business": regexFilter(keyword)},
			{"address": regexFilter(keyword)},
			{"phone": regexFilter(keyword)},
			{"email": regexFilter(keyword)},
			{"website": regexFilter(keyword)},
			{"discount": regexFilter(keyword)},
		}
	} else {
		if category != "" {
			normalizedCategory := strings.ToLower(category)
			canonicalCategory, validCategory := discountProgramCategories[normalizedCategory]
			if !validCategory {
				respond(
					c,
					http.StatusBadRequest,
					"Invalid query parameters",
					"category must be one of: "+strings.Join(discountProgramCategoryList, ", "),
				)
				return
			}
			query["category"] = bson.M{
				"$regex":   "^" + regexp.QuoteMeta(canonicalCategory) + "$",
				"$options": "i",
			}
		}

		if business != "" {
			query["business"] = regexFilter(business)
		}
		if address != "" {
			query["address"] = regexFilter(address)
		}
		if discount != "" {
			query["discount"] = regexFilter(discount)
		}
	}

	options, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}
	options.SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := discountProgramCollection.Find(ctx, query, options)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	var discountPrograms []schema.DiscountProgram
	if err = cursor.All(ctx, &discountPrograms); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", discountPrograms)
}

func regexFilter(term string) bson.M {
	return bson.M{
		"$regex":   regexp.QuoteMeta(term),
		"$options": "i",
	}
}
