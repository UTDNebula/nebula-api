package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: %v", err)
	}

	return os.Getenv("MONGODB_URI")
}

func EnvLimit() int64 {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: %v", err)
	}

	limit, err := strconv.ParseInt(os.Getenv("LIMIT"), 10, 64)
	if err != nil {
		limit = 20 // default value for limit
	}
	return limit
}

var Limit int64 = EnvLimit()
