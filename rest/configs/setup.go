package configs

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBSingleton struct {
	client *mongo.Client
}

var dbInstance *DBSingleton
var once sync.Once

func ConnectDB() *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(GetEnvMongoURI()))
		if err != nil {
			log.Fatalf("Unable to create MongoDB client")
		}

		defer cancel()

		// ping the database
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Unable to ping database")
		}

		log.Printf("Connected to MongoDB")

		dbInstance = &DBSingleton{
			client: client,
		}
	})

	return dbInstance.client
}

// getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	client := ConnectDB()
	collection := client.Database("combinedDB").Collection(collectionName)
	return collection
}

// Returns *options.FindOptions with a limit and offset applied.
// Produces an error if user-provided offset isn't able to be parsed.
func GetOptionLimit(query *bson.M, c *gin.Context) (*options.FindOptions, error) {
	delete(*query, "offset") // removes offset (if present) in query --offset is not field in collections

	// parses offset if included in the query
	var offset int64
	var err error

	var limit int64 = GetEnvLimit()

	if c.Query("offset") == "" {
		offset = 0 // default value for offset
	} else {
		offset, err = strconv.ParseInt(c.Query("offset"), 10, 64)
		if err != nil {
			return options.Find().SetSkip(0).SetLimit(limit), err // default value for offset
		}
	}

	return options.Find().SetSkip(offset).SetLimit(limit), err
}

// Returns the offsets and limit for pagination stage for aggregate endpoints pipeline
func GetAggregateLimit(query *bson.M, c *gin.Context) (map[string]bson.D, error) {
	// Parses offsets if included in the query
	paginateMap := map[string]bson.D{
		"former_offset": {{Key: "$skip", Value: 0}}, // Init the default value of offset
		"latter_offset": {{Key: "$skip", Value: 0}},
		"limit":         {{Key: "$limit", Value: GetEnvLimit()}},
	}
	var err error

	// Loop through offset types (keys indicating offset values)
	for field := range paginateMap {
		// Only change values of the map if specified
		if field != "limit" && c.Query(field) != "" {
			// Remove offset field (if present) in the query
			delete(*query, field)

			// Build the stage from the parsed field
			offset, err := strconv.ParseInt(c.Query(field), 10, 64)
			if err != nil {
				// Return default value of offset
				return paginateMap, err
			}
			paginateMap[field] = bson.D{{Key: "$skip", Value: offset}}
		}
	}

	return paginateMap, err
}

var clubsDbInstance *sql.DB
var clubOnce sync.Once

func ConnectClubsDB() *sql.DB {
	clubOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		db, err := sql.Open("pgx", GetClubsDBUri())
		if err != nil {
			log.Panic("Unable to connect to clubs database.")
		}

		defer cancel()

		// ping the database
		err = db.PingContext(ctx)
		if err != nil {
			log.Panic("Unable to ping database")
		}

		log.Printf("Connected to Clubs DB")

		clubsDbInstance = db
	})

	return clubsDbInstance
}

// make sure this file is in gitignore if changed
const SemestersJSONFilename = "semesters.json"

func CreateSemestersJSON() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	terms, err := GetCollection("sections").Distinct(ctx, "academic_session.name", bson.D{})
	if err != nil {
		return err
	}

	type term struct {
		name   string
		year   int
		season string
	}

	var parsedTerms []term
	seen := make(map[string]struct{})

	for _, t := range terms {
		str, ok := t.(string)
		if !ok || len(str) < 3 {
			continue
		}

		str = strings.ToLower(strings.TrimSpace(str))

		// Solution for possible duplicate value case (if database query returns '24F' and '24f')
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

	jsonData, err := json.Marshal(sortedTerms)
	if err != nil {
		return err
	}

	err = os.WriteFile(SemestersJSONFilename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
