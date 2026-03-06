package controllers

import (
	"context"
	"fmt"
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

var discountCollection *mongo.Collection = configs.GetCollection("discounts")

// discountCategories
// potentially we may want to add an init function to create this list dynamically
var discountCategories = []string{
	"Accommodations",
	"Auto Services",
	"Child Care",
	"Clothes, Flowers and Gifts",
	"Dining",
	"Entertainment",
	"Health and Beauty",
	"Home and Garden",
	"Housing",
	"Miscellaneous",
	"Pet Care",
	"Professional Services",
	"Technology",
}

// @Id				discountPrograms
// @Router			/discountPrograms [get]
// @Tags			Discounts
// @Description	"Returns paginated list of discounts matching the query's string-typed key-value pairs. See offset for more details on pagination."
// @Produce		json
// @Param			offset		query		number											false	"The starting position of the current page of discounts (e.g. For starting at the 17th discount, offset=16)."
// @Param			category	query		string											false	"The discount's category."
// @Param			business	query		string											false	"The discount's business contains this keyword (case-insensitive)."
// @Param			address		query		string											false	"The discount's address contains this keyword (case-insensitive)."
// @Param			discount	query		string											false	"The discount's discount contains this keyword (case-insensitive)."
// @Param			q			query		string											false	"Full text search of all discount's fields."
// @Success		200			{object}	schema.APIResponse[[]schema.DiscountProgram]	"A list of discounts"
// @Failure		500			{object}	schema.APIResponse[string]						"A string describing the error"
// @Failure		400			{object}	schema.APIResponse[string]						"A string describing the error"
func DiscountSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query, err := buildDiscountSearchQuery(c)
	if err != nil {
		respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
		return
	}

	cursor, err := discountCollection.Find(ctx, query, optionLimit)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	defer cursor.Close(ctx)

	discounts := make([]schema.DiscountProgram, 0)
	if err = cursor.All(ctx, &discounts); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", discounts)

}

func buildDiscountSearchQuery(c *gin.Context) (bson.M, error) {
	business, hasBusiness := c.GetQuery("business")
	address, hasAddress := c.GetQuery("address")
	discount, hasDiscount := c.GetQuery("discount")
	category, hasCategory := c.GetQuery("category")
	q, hasQ := c.GetQuery("q")

	query := bson.M{}

	if hasQ {
		// q may only be used alone
		if hasBusiness || hasAddress || hasDiscount || hasCategory {
			return nil, fmt.Errorf("parameter q may not be used with other parameters")
		}
		query["$text"] = bson.D{{Key: "$search", Value: q}}
		return query, nil
	}

	// We use regexp.QuoteMeta and option i to essentially do string.toLower().contains(key) on fields
	if hasBusiness {
		cleanedBusiness := strings.TrimSpace(regexp.QuoteMeta(business))
		query["business"] = bson.D{{Key: "$regex", Value: cleanedBusiness}, {Key: "$options", Value: "i"}}
	}
	if hasAddress {
		cleanedAddress := strings.TrimSpace(regexp.QuoteMeta(address))
		query["address"] = bson.D{{Key: "$regex", Value: cleanedAddress}, {Key: "$options", Value: "i"}}
	}
	if hasDiscount {
		cleanedDiscount := strings.TrimSpace(regexp.QuoteMeta(discount))
		query["discount"] = bson.D{{Key: "$regex", Value: cleanedDiscount}, {Key: "$options", Value: "i"}}
	}

	if hasCategory {
		found := false
		for _, discountCategory := range discountCategories {
			//case insensitive equal
			if strings.EqualFold(discountCategory, category) {
				query["category"] = discountCategory
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("Unknown category %s. Valid categories are %s", category, strings.Join(discountCategories, ", "))
		}
	}

	return query, nil
}
