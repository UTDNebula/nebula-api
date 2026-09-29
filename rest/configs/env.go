package configs

import (
	"fmt"
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

func GetEnvMongoURI() (string, error) {

	uri, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		return "", fmt.Errorf("environment variable %q is not set", "MONGODB_URI")
	}

	return uri, nil
}

func GetClubsDBUri() (string, error) {
	uri, exist := os.LookupEnv("CLUBS_DB_URI")
	if !exist {
		return "", fmt.Errorf("environment variable %q is not set", "CLUBS_DB_URI")
	}

	return uri, nil
}

func GetEnvLogin() (netID string, password string, err error) {

	netID, exist := os.LookupEnv("LOGIN_NETID")
	if !exist {
		return "", "", fmt.Errorf("environment variable %q is not set", "LOGIN_NETID")
	}
	password, exist = os.LookupEnv("LOGIN_PASSWORD")
	if !exist {
		return "", "", fmt.Errorf("environment variable %q is not set", "LOGIN_PASSWORD")
	}

	return netID, password, nil
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

func GetEnvSentryDSN() string {

	// Sentry is disabled when initialized without a DSN, so an unset
	// SENTRY_DSN keeps development errors out of Sentry. Only production
	// sets this variable.
	return os.Getenv("SENTRY_DSN")
}
