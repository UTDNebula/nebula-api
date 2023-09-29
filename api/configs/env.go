package configs

import (
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func EnvMongoURI() string {

	uri, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		log.Fatalf("Error loading 'MONGODB_URI' from the .env file")
	}

	return uri
}

func EnvLimit() int64 {

	const defaultLimit int64 = 20

	limitString, exist := os.LookupEnv("LIMIT")
	if !exist {
		return defaultLimit
	}

	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		return defaultLimit
	}

	return limit
}

var Limit int64 = EnvLimit()
