package configs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

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
var dbErr error
var once sync.Once

func ConnectDB(uri string) error {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			dbErr = fmt.Errorf("create MongoDB client: %w", err)
			return
		}

		// ping the database
		err = client.Ping(ctx, nil)
		if err != nil {
			dbErr = fmt.Errorf("ping MongoDB: %w", err)
			return
		}

		log.Printf("Connected to MongoDB")

		dbInstance = &DBSingleton{
			client: client,
		}
	})

	return dbErr
}

// getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	collection := dbInstance.client.Database("combinedDB").Collection(collectionName)
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
var clubsDBErr error
var clubOnce sync.Once

func ConnectClubsDB(uri string) error {
	clubOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db, err := sql.Open("pgx", uri)
		if err != nil {
			clubsDBErr = fmt.Errorf("open clubs database: %w", err)
			return
		}

		// ping the database
		err = db.PingContext(ctx)
		if err != nil {
			clubsDBErr = fmt.Errorf("ping clubs database: %w", err)
			return
		}

		log.Printf("Connected to Clubs DB")

		clubsDbInstance = db
	})

	return clubsDBErr
}

func GetClubsDB() *sql.DB {
	return clubsDbInstance
}
