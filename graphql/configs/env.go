package configs

import (
	"os"
	"strconv"

	"log"

	_ "github.com/joho/godotenv/autoload"
)

// Get the connection port (if any) from environment
func GetPortString() string {
	portNumber, exist := os.LookupEnv("PORT")
	if !exist {
		portNumber = "8000" // 8080 for REST, 8000 for GRAPHQL
	}

	return portNumber
}

// Get the connection string of the MongoDB from environment
func GetEnvMongoURI() string {
	uri, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		log.Fatalf("Error loading 'MONGODB_URI' from the .env file")
	}

	return uri
}

// Get the netID and password from env
func GetEnvLogin() (netID string, password string) {
	netID, exist := os.LookupEnv("LOGIN_NETID")
	if !exist {
		log.Fatalf("Error loading 'LOGIN_NETID' from the .env file")
	}
	password, exist = os.LookupEnv("LOGIN_PASSWORD")
	if !exist {
		log.Fatalf("Error loading 'LOGIN_PASSWORD' from the .env file")
	}

	return netID, password
}

func GetEnvLimit() int64 {
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
