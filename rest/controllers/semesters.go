package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type cacheFile struct {
	LastUpdated time.Time `json:"lastUpdated"`
	Terms       []string  `json:"terms"`
}

const (
	filename = "semesters.json"
	fileTTL  = 30 * 24 * time.Hour // 30 days
)

var (
	semesters   []string
	lastUpdated time.Time
	isLoaded    bool
	mutex       sync.RWMutex
)

// @Id				semesters
// @Router			/semesters [get]
// @Tags			Other
// @Description	"Returns a list of all available terms in CourseBook, sorted by recency"
// @Produce		json
// @Success		200	{object}	schema.APIResponse[[]string]	"A string list of all available terms, lowercased"
// @Failure		500	{object}	schema.APIResponse[string]		"A string describing the error"
func Semesters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	mutex.RLock()
	needsRefresh := !isLoaded || cacheExpired(lastUpdated)
	mutex.RUnlock()

	if needsRefresh {
		err := loadOrRefresh(ctx)
		if err != nil {
			respondWithInternalError(c, err)
			return
		}
	}

	mutex.RLock()
	data := semesters
	mutex.RUnlock()

	respond(c, http.StatusOK, "success", data)
}

func loadOrRefresh(ctx context.Context) error {
	// Prevent other calls from r/w json file
	mutex.Lock()
	defer mutex.Unlock()

	if isLoaded && !cacheExpired(lastUpdated) {
		return nil
	}

	if !isLoaded {
		// Check if json file exists
		fileBytes, err := os.ReadFile(filename)
		if err == nil {
			var fileData cacheFile
			if err := json.Unmarshal(fileBytes, &fileData); err == nil {
				// Check if file data is outdated, load and return if not
				if time.Since(fileData.LastUpdated) <= fileTTL {
					semesters = fileData.Terms
					lastUpdated = fileData.LastUpdated
					isLoaded = true
					return nil
				}
			}
		}
	}

	sortedTerms, err := getSortedTerms(ctx)
	if err != nil {
		return err
	}
	currTime := time.Now()

	newData := cacheFile{
		LastUpdated: currTime,
		Terms:       sortedTerms,
	}

	fileBytes, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	// Create new json file with updated data
	err = os.WriteFile(filename, fileBytes, 0644)
	if err != nil {
		return err
	}

	// Update global variables
	semesters = sortedTerms
	lastUpdated = currTime
	isLoaded = true

	return nil
}

func getSortedTerms(ctx context.Context) ([]string, error) {

	type term struct {
		name   string
		year   int
		season string
	}

	terms, err := sectionCollection.Distinct(ctx, "academic_session.name", bson.D{})
	if err != nil {
		return nil, err
	}

	var parsedTerms []term
	seen := make(map[string]struct{})

	for _, t := range terms {
		str, ok := t.(string)
		if !ok || len(str) < 3 {
			continue
		}

		str = strings.ToLower(strings.TrimSpace(str))

		// Solution to possible duplicate value case (if database query returns '24F' and '24f')
		if _, exists := seen[str]; exists {
			continue
		}

		year, err := strconv.Atoi(str[:2])
		if err != nil {
			continue
		}
		parsedTerms = append(parsedTerms, term{
			name:   str,
			year:   year,
			season: str[2:],
		})
	}

	// Sort terms by recency
	sort.Slice(parsedTerms, func(i, j int) bool {
		if parsedTerms[i].year != parsedTerms[j].year {
			return parsedTerms[i].year > parsedTerms[j].year
		}

		// Sort by season (comparing ascii values)
		return parsedTerms[i].season > parsedTerms[j].season
	})

	sortedTerms := make([]string, len(parsedTerms))
	for i, t := range parsedTerms {
		sortedTerms[i] = t.name
	}

	return sortedTerms, nil
}

func cacheExpired(updatedAt time.Time) bool {
	return updatedAt.IsZero() || time.Since(updatedAt) > fileTTL
}
