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
	"go.mongodb.org/mongo-driver/mongo/options"
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
// @Description	"Returns paginated list of discounts filtered using field-specific keyword searches or global fuzzy search. See offset for more details on pagination."
// @Produce		json
// @Param			offset		query		number											false	"The starting position of the current page of discounts (e.g. For starting at the 17th discount, offset=16)."
// @Param			category	query		string											false	"The discount's category (exact match with suggestions)."
// @Param			business	query		string											false	"The discount's business contains this keyword (case-insensitive)."
// @Param			address		query		string											false	"The discount's address contains this keyword (case-insensitive)."
// @Param			discount	query		string											false	"The discount's discount contains this keyword (case-insensitive)."
// @Param			q			query		string											false	"Fuzzy search, must be used alone."
// @Success		200			{object}	schema.APIResponse[[]schema.DiscountProgram]	"A list of discounts"
// @Failure		500			{object}	schema.APIResponse[string]						"A string describing the error"
// @Failure		400			{object}	schema.APIResponse[string]						"A string describing the error"
func DiscountSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cursor *mongo.Cursor
	var err error

	var params schema.DiscountQueryParams
	if err = c.ShouldBindQuery(&params); err != nil {
		respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}
	params.TrimSpace()

	if params.Q != "" {
		if params.Category != "" || params.Business != "" ||
			params.Address != "" || params.Discount != "" {
			respond(c, http.StatusBadRequest, "Invalid query parameters", "q must be used alone")
			return
		}

		fuzzyPipeline := buildFuzzySearchPipeline(params.Q, params.Offset)
		cursor, err = discountCollection.Aggregate(ctx, fuzzyPipeline)
		if err != nil {
			respondWithInternalError(c, err)
			return
		}
	} else {
		// Either fields are specified or not
		query, err := buildDiscountSearchQuery(params)
		if err != nil {
			respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
			return
		}

		optionLimit := options.Find().SetSkip(params.Offset).SetLimit(configs.GetEnvLimit())
		cursor, err = discountCollection.Find(ctx, query, optionLimit)
		if err != nil {
			respondWithInternalError(c, err)
			return
		}
	}
	defer cursor.Close(ctx)

	discounts := make([]schema.DiscountProgram, 0)
	if err = cursor.All(ctx, &discounts); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", discounts)

}

// buildDiscountSearchQuery constructs the Mongo query for FIELD-BASED SEARCH.
// Users only search for 4 main fields of discount
func buildDiscountSearchQuery(q schema.DiscountQueryParams) (bson.M, error) {
	query := bson.M{}

	// We use regexp.QuoteMeta and option i to essentially do string.toLower().contains(key) on fields
	if q.Business != "" {
		business := regexp.QuoteMeta(q.Business)
		query["business"] = bson.D{{Key: "$regex", Value: business}, {Key: "$options", Value: "i"}}
	}
	if q.Address != "" {
		address := regexp.QuoteMeta(q.Address)
		query["address"] = bson.D{{Key: "$regex", Value: address}, {Key: "$options", Value: "i"}}
	}
	if q.Discount != "" {
		discount := regexp.QuoteMeta(q.Discount)
		query["discount"] = bson.D{{Key: "$regex", Value: discount}, {Key: "$options", Value: "i"}}
	}

	if q.Category != "" {
		categoryFound := false
		for _, discountCategory := range discountCategories {
			if discountCategory == q.Category {
				query["category"] = discountCategory
				categoryFound = true
				break
			}
		}
		if !categoryFound {
			return nil, fmt.Errorf("unknown category, valid categories are %s", strings.Join(discountCategories, " | "))
		}
	}

	return query, nil
}

// buildFuzzySearchPipeline constructs the pipeline to perform fuzzy search on keyword q.
func buildFuzzySearchPipeline(q string, offset int64) mongo.Pipeline {
	type FuzzyConfig struct {
		Field      string
		maxEdits   int
		boostScore int
	}
	// Will need to tune the configuration to get better results
	fuzzyConfigs := []FuzzyConfig{
		{"category", 2, 5},
		{"discount", 2, 3},
		{"business", 2, 2},
		{"address", 1, 1},
	}

	var fuzzySearches bson.A
	for _, fuzzyConfig := range fuzzyConfigs {
		fuzzySearches = append(fuzzySearches, bson.D{
			{Key: "text", Value: bson.D{
				{Key: "query", Value: q},
				{Key: "path", Value: fuzzyConfig.Field},
				{Key: "fuzzy", Value: bson.D{
					{Key: "maxEdits", Value: fuzzyConfig.maxEdits},
					{Key: "prefixLength", Value: 2}, // Should match first 2 characters
				}},
				{Key: "score", Value: bson.D{
					{Key: "boost", Value: bson.D{
						{Key: "value", Value: fuzzyConfig.boostScore},
					}},
				}},
			}},
		})
	}

	return mongo.Pipeline{
		bson.D{
			{Key: "$search", Value: bson.D{
				{Key: "index", Value: "discount_searches"},
				{Key: "compound", Value: bson.D{
					{Key: "should", Value: fuzzySearches},
					{Key: "minimumShouldMatch", Value: 1}, // Prevent extremely unrelated docs
				}},
			}},
		},

		// Sort based on relevancy score for determinism and paginate
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "score", Value: bson.D{
					{Key: "$meta", Value: "searchScore"},
				}},
			}},
		},
		bson.D{{Key: "$skip", Value: offset}},
		bson.D{{Key: "$limit", Value: configs.GetEnvLimit()}},
	}
}
