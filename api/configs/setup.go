package configs

import (
	"context"
	"strconv"
	"sync"
	"time"

	"log"

	"github.com/gin-gonic/gin"

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

		//ping the database
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

// Returns *options.FindOptions with a limit and offset applied. Produces an error if user-provided offset isn't able to be parsed.
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

// Returns the offsets and limit for pagination stage for aggregate endpoints pipeline (map, err)
func GetAggregateLimit(query *bson.M, c *gin.Context) (map[string]int64, error) {
	// remove formerOffset and latterOffset field (if present) in the query
	delete(*query, "former_offset")
	delete(*query, "latter_offset")

	// parses offsets if included in the query
	paginateMap := map[string]int64{
		"former_offset": 0, // initialize the default value of offset & limit right in the map
		"latter_offset": 0,
		"limit":         GetEnvLimit(),
	}
	var err error

	// loop through offset types (keys indicating offset values)
	for key := range paginateMap {
		// only change values of the map if specified
		if key != "limit" && c.Query(key) != "" {
			offset, parseErr := strconv.ParseInt(c.Query(key), 10, 64)
			if parseErr != nil {
				return paginateMap, parseErr // return default value of offset
			}
			paginateMap[key] = offset
		}
	}

	return paginateMap, err
}
