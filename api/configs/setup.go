package configs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetEnvMongoURI()))
	if err != nil {
		log.Fatalf("Unable to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
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
