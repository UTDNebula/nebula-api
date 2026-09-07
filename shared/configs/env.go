package configs

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

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

func GetPortString(defaultPort string) string {
	portNumber, exists := os.LookupEnv("PORT")
	if !exists {
		portNumber = defaultPort
	}

	return portNumber
}

func GetEnvMongoURI() string {
	uri, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		log.Fatalf("Error loading 'MONGODB_URI' from the .env file")
	}

	return uri
}

func GetEnvLogin() (netID string, password string) {
	netID, exists := os.LookupEnv("LOGIN_NETID")
	if !exists {
		log.Fatalf("Error loading 'LOGIN_NETID' from the .env file")
	}
	password, exists = os.LookupEnv("LOGIN_PASSWORD")
	if !exists {
		log.Fatalf("Error loading 'LOGIN_PASSWORD' from the .env file")
	}

	return netID, password
}

func GetEnvLimit() int64 {
	const defaultLimit int64 = 20

	limitString, exists := os.LookupEnv("LIMIT")
	if !exists {
		return defaultLimit
	}

	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		return defaultLimit
	}

	return limit
}

func GetClubsDBUri() string {
	uri, exists := os.LookupEnv("CLUBS_DB_URI")
	if !exists {
		log.Panic("Error loading 'CLUBS_DB_URI' from the .env file")
	}

	return uri
}

func GetEnvMaxUploadSize() int64 {
	const (
		defaultLimit int64 = 30 * 1024 * 1024
		hardCapLimit int64 = 50 * 1024 * 1024
	)

	limitString, exists := os.LookupEnv("MAX_UPLOAD_SIZE")
	if !exists {
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
