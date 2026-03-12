package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var discountCollection *mongo.Collection = configs.GetCollection("discounts")

var discountCategories []string
var discountCategoriesOnce sync.Once

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
	fetchDiscountCategories()
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
		if params.HasFields() {
			respond(c, http.StatusBadRequest, "Invalid query parameters", "q must be used alone")
			return
		}

		pipeline := buildFuzzySearchPipeline(params.Q, params.Offset)
		cursor, err = discountCollection.Aggregate(ctx, pipeline)
		if err != nil {
			respondWithInternalError(c, err)
			return
		}
	} else {
		// If no fields are specified, it just returns paginated collection
		discountQuery, err := buildDiscountSearchQuery(params)
		if err != nil {
			respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
			return
		}

		optionLimit := options.Find().SetSkip(params.Offset).SetLimit(configs.GetEnvLimit())
		cursor, err = discountCollection.Find(ctx, discountQuery, optionLimit)
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

// initDiscountCategories aggregates the list of discount categories from DB once
func fetchDiscountCategories() {
	discountCategoriesOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()

		results, err := discountCollection.Distinct(ctx, "category", bson.M{})
		if err != nil {
			panic(err)
		}
		for _, result := range results {
			category, ok := result.(string)
			if !ok {
				continue // Skip invalid category
			}
			discountCategories = append(discountCategories, category)
		}
		log.Printf("Available discount categories: %s.\n", discountCategories)
	})
}

// buildDiscountSearchQuery constructs the Mongo query for FIELD-BASED SEARCH.
// Users only search for 4 main fields of discount
func buildDiscountSearchQuery(p schema.DiscountQueryParams) (bson.M, error) {
	query := bson.M{}

	// Use regexp.QuoteMeta and option i for insensitive-matching
	if p.Business != "" {
		query["business"] = bson.D{
			{Key: "$regex", Value: regexp.QuoteMeta(p.Business)},
			{Key: "$options", Value: "i"},
		}
	}
	if p.Address != "" {
		query["address"] = bson.D{
			{Key: "$regex", Value: regexp.QuoteMeta(p.Address)},
			{Key: "$options", Value: "i"},
		}
	}
	if p.Discount != "" {
		query["discount"] = bson.D{
			{Key: "$regex", Value: regexp.QuoteMeta(p.Discount)},
			{Key: "$options", Value: "i"},
		}
	}
	if p.Category != "" {
		categoryFound := false
		for _, discountCategory := range discountCategories {
			if discountCategory == p.Category {
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
	var fuzzySearches bson.A
	// Tune the configuration to get better results
	fuzzyConfigs := []schema.FuzzySearchConfig{
		{Field: "category", MaxEdits: 2, BoostScore: 5},
		{Field: "discount", MaxEdits: 2, BoostScore: 3},
		{Field: "business", MaxEdits: 2, BoostScore: 2},
		{Field: "address", MaxEdits: 1, BoostScore: 1},
	}
	for _, fuzzyConfig := range fuzzyConfigs {
		fuzzySearches = append(fuzzySearches, bson.D{
			{Key: "text", Value: bson.D{
				{Key: "query", Value: q},
				{Key: "path", Value: fuzzyConfig.Field},
				{Key: "fuzzy", Value: bson.D{
					{Key: "maxEdits", Value: fuzzyConfig.MaxEdits},
					{Key: "prefixLength", Value: 2}, // Should match first 2 chars
				}},
				{Key: "score", Value: bson.D{
					{Key: "boost", Value: bson.D{
						{Key: "value", Value: fuzzyConfig.BoostScore},
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

					// Match at least 1 field to prevent super unrelated docs
					{Key: "minimumShouldMatch", Value: 1},
				}},
			}},
		},

		// Sort based on relevancy for determinism and paginate
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "score", Value: bson.D{
					{Key: "$meta", Value: "searchScore"},
				}},
			}},
		},

		// Paginate the results
		bson.D{{Key: "$skip", Value: offset}},
		bson.D{{Key: "$limit", Value: configs.GetEnvLimit()}},
	}
}
