package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Initialize this file to load environment variables from .env file
func init() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func GetPortString() string {

	portNumber, exist := os.LookupEnv("PORT")
	if !exist {
		portNumber = "8080"
	}

	portString := fmt.Sprintf(":%s", portNumber)

	return portString
}

func GetEnvMongoURI() string {

	uri, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		log.Fatalf("Error loading 'MONGODB_URI' from the .env file")
	}

	return uri
}

func GetClubsDBUri() string {
	uri, exist := os.LookupEnv("CLUBS_DB_URI")
	if !exist {
		log.Panic("Error loading 'CLUBS_DB_URI' from the .env file")
	}

	return uri
}

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

func GetEnvMaxUploadSize() int64 {
	const (
		defaultLimit int64 = 30 * 1024 * 1024
		hardCapLimit int64 = 50 * 1024 * 1024
	)

	limitString, exist := os.LookupEnv("MAX_UPLOAD_SIZE")
	if !exist {
		return defaultLimit
	}

	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		return defaultLimit
	}

	if limit > hardCapLimit {
		return hardCapLimit
	}

	return limit
}

func GetSentryEnv() (string, string) {
	dsn, exist := os.LookupEnv("SENTRY_DSN")
	if !exist {
		log.Fatalf("Error loading 'SENTRY_DSN' from the .env file")
	}

	env, exist := os.LookupEnv("SENTRY_ENVIRONMENT")
	if !exist {
		log.Fatalf("Error loading 'SENTRY_ENVIRONMENT' from the .env file")
	}

	return dsn, env
}
