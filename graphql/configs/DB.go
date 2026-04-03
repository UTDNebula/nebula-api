package configs

import (
	"context"
	"sync"
	"time"

	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBSingleton struct {
	client *mongo.Client
}

var dbInstance *DBSingleton
var once sync.Once

// Connect to the Mongo database
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

// Get database collections
func GetCollection(collectionName string) *mongo.Collection {
	client := ConnectDB()
	collection := client.Database("combinedDB").Collection(collectionName)
	return collection
}
