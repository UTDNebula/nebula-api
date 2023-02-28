package configs

import (
	"log"
	"os"
	"strconv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return uri
}

func EnvLimit() int64 {

	limit, err := strconv.ParseInt(os.Getenv("LIMIT"), 10, 64)
	if err != nil {
		limit = 20 // default value for limit
	}
	return limit
}

var Limit int64 = EnvLimit()
