package configs

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
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
			log.WriteErrorMsg("Unable to create MongoDB client")
			os.Exit(1)
		}

		defer cancel()

		//ping the database
		err = client.Ping(ctx, nil)
		if err != nil {
			log.WriteErrorMsg("Unable to ping database")
			os.Exit(1)
		}

		log.WriteDebug("Connected to MongoDB")

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

// Returns *options.FindOptions with a limit and offset applied. Returns error if any
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
// (former offset, latter offset, limit, err)
func GetAggregateLimit(query *bson.M, c *gin.Context) (int64, int64, int64, error) {
	// remove formerOffset and latterOffset field (if present) in the query
	delete(*query, "former_offset")
	delete(*query, "latter_offset")

	// parses offset if included in the query
	var formerOffset, latterOffset int64
	var err error

	var limit int64 = GetEnvLimit()

	// get the offset on the "former" part of the endpoint
	if c.Query("former_offset") == "" {
		formerOffset = 0
	} else {
		formerOffset, err = strconv.ParseInt(c.Query("former_offset"), 10, 64)
		if err != nil {
			return 0, 0, limit, err
		}
	}

	// get offset on the "latter" part of the endpoint
	if c.Query("latter_offset") == "" {
		latterOffset = 0
	} else {
		latterOffset, err = strconv.ParseInt(c.Query("latter_offset"), 10, 64)
		if err != nil {
			return 0, 0, limit, err
		}
	}

	return formerOffset, latterOffset, limit, err
}
