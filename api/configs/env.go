package configs

import (
	"log"
	"os"
	"strconv"
)

func EnvMongoURI() string {

	uri, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		log.Fatal("Error loading 'MONGODB_URI' from the .env file")
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
