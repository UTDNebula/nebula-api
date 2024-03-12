package configs

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {

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

	return client
}

var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
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
