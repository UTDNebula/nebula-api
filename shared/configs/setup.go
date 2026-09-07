package configs

import (
	"context"
	"database/sql"
	"strconv"
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
		defer cancel()

		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(GetEnvMongoURI()))
		if err != nil {
			log.Fatalf("Unable to create MongoDB client")
		}

		if err = client.Ping(ctx, nil); err != nil {
			log.Fatalf("Unable to ping database")
		}

		log.Printf("Connected to MongoDB")
		dbInstance = &DBSingleton{client: client}
	})

	return dbInstance.client
}

func GetCollection(collectionName string) *mongo.Collection {
	return ConnectDB().Database("combinedDB").Collection(collectionName)
}

func GetOptionLimit(query *bson.M, c *gin.Context) (*options.FindOptions, error) {
	delete(*query, "offset")

	var offset int64
	var err error
	limit := GetEnvLimit()

	if c.Query("offset") == "" {
		offset = 0
	} else {
		offset, err = strconv.ParseInt(c.Query("offset"), 10, 64)
		if err != nil {
			return options.Find().SetSkip(0).SetLimit(limit), err
		}
	}

	return options.Find().SetSkip(offset).SetLimit(limit), err
}

func GetAggregateLimit(query *bson.M, c *gin.Context) (map[string]bson.D, error) {
	paginateMap := map[string]bson.D{
		"former_offset": {{Key: "$skip", Value: 0}},
		"latter_offset": {{Key: "$skip", Value: 0}},
		"limit":         {{Key: "$limit", Value: GetEnvLimit()}},
	}

	for field := range paginateMap {
		if field != "limit" && c.Query(field) != "" {
			delete(*query, field)
			offset, err := strconv.ParseInt(c.Query(field), 10, 64)
			if err != nil {
				return paginateMap, err
			}
			paginateMap[field] = bson.D{{Key: "$skip", Value: offset}}
		}
	}

	return paginateMap, nil
}

var clubsDbInstance *sql.DB
var clubOnce sync.Once

func ConnectClubsDB() *sql.DB {
	clubOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db, err := sql.Open("pgx", GetClubsDBUri())
		if err != nil {
			log.Panic("Unable to connect to clubs database.")
		}

		if err = db.PingContext(ctx); err != nil {
			log.Panic("Unable to ping clubs database")
		}

		log.Printf("Connected to Clubs DB")
		clubsDbInstance = db
	})

	return clubsDbInstance
}
