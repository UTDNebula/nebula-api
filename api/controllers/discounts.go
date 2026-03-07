package controllers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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

	var cursor *mongo.Cursor
	var err error

	_, hasQ := c.GetQuery("q")
	_, hasBusiness := c.GetQuery("business")
	_, hasAddress := c.GetQuery("address")
	_, hasDiscount := c.GetQuery("discount")
	_, hasCategory := c.GetQuery("category")
	if hasQ {
		if hasBusiness || hasAddress || hasDiscount || hasCategory {
			// q may only be used alone
			respond(c, http.StatusBadRequest, "Invalid query parameters", "Parameter q may not be used with other parameters")
			return
		}

		pipeline, err := buildFuzzySearchPipeline(c)
		if err != nil {
			respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
			return
		}
		cursor, err = discountCollection.Aggregate(ctx, pipeline)
		if err != nil {
			respondWithInternalError(c, err)
			return
		}

	} else {
		if !hasBusiness && !hasAddress && !hasDiscount && !hasCategory {
			respond(c, http.StatusBadRequest, "Invalid query parameters", "Unknown query")
			return
		}

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
		cursor, err = discountCollection.Find(ctx, query, optionLimit)
		if err != nil {
			respondWithInternalError(c, err)
			return
		}
	}

	defer cursor.Close(ctx)

	var discounts []schema.DiscountProgram
	if err = cursor.All(ctx, &discounts); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", discounts)

}

// buildDiscountSearchQuery constructs the Mongo query for FIELD-BASED SEARCH.
// Users only search for 4 main fields of discount
func buildDiscountSearchQuery(c *gin.Context) (bson.M, error) {
	business, hasBusiness := c.GetQuery("business")
	address, hasAddress := c.GetQuery("address")
	discount, hasDiscount := c.GetQuery("discount")
	category, hasCategory := c.GetQuery("category")

	query := bson.M{}

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
		categoryFound := false
		for _, discountCategory := range discountCategories {
			// Case insensitive equal
			if strings.EqualFold(discountCategory, category) {
				query["category"] = discountCategory
				categoryFound = true
				break
			}
		}
		if !categoryFound {
			return nil, fmt.Errorf("unknown category %s. Valid categories are %s", category, strings.Join(discountCategories, " | "))
		}
	}

	return query, nil
}

// buildFuzzySearchPipeline constructs the pipeline to perform fuzzy search on keyword q.
func buildFuzzySearchPipeline(c *gin.Context) (mongo.Pipeline, error) {
	q, _ := c.GetQuery("q")
	if strings.TrimSpace(q) == "" {
		return mongo.Pipeline{}, fmt.Errorf("empty q")
	}

	// Literally copy from getOptionsLimit()
	var offset int64
	var err error
	if c.Query("offset") == "" {
		offset = 0
	} else {
		offset, err = strconv.ParseInt(c.Query("offset"), 10, 64)
		if err != nil {
			return mongo.Pipeline{}, err
		}
	}

	var fuzzySearchArr bson.A
	fields := [4]string{"category", "discount", "business", "address"}
	maxEditsList := [4]int{2, 2, 2, 1}
	boostScores := [4]int{5, 3, 2, 1}
	for i, field := range fields {
		fuzzySearchArr = append(fuzzySearchArr, bson.D{
			{Key: "text", Value: bson.D{
				{Key: "query", Value: q},
				{Key: "path", Value: field},
				{Key: "fuzzy", Value: bson.D{
					{Key: "maxEdits", Value: maxEditsList[i]},
				}},
				{Key: "score", Value: bson.D{
					{Key: "boost", Value: bson.D{
						{Key: "value", Value: boostScores[i]},
					}},
				}},
			}},
		})
	}

	return mongo.Pipeline{
		// Fuzzy searches
		bson.D{
			{Key: "$search", Value: bson.D{
				{Key: "index", Value: "discount_searches"},
				{Key: "compound", Value: bson.D{
					{Key: "should", Value: fuzzySearchArr},
				}},
			}},
		},

		// Sort based on relevancy score for deterministism and paginate
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "score", Value: bson.D{
					{Key: "$meta", Value: "searchScore"},
				}},
			}},
		},
		bson.D{{Key: "$skip", Value: offset}},
		bson.D{{Key: "$limit", Value: configs.GetEnvLimit()}},
	}, nil
}
