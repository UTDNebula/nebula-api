package configs

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

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

// A misconfigured variable is a single, process-wide mistake, but these readers
// run on every request that paginates or uploads. Each warning therefore fires
// at most once per process rather than once per request.
var (
	limitWarnOnce         sync.Once
	uploadSizeWarnOnce    sync.Once
	uploadSizeCapWarnOnce sync.Once
)

func GetPortString() string {

	portNumber, exist := os.LookupEnv("PORT")
	if !exist {
		portNumber = "8080"
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
		limitWarnOnce.Do(func() {
			log.Printf("Ignoring 'LIMIT' value %q: not a number, using the default of %d\n", limitString, defaultLimit)
		})
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
		uploadSizeWarnOnce.Do(func() {
			log.Printf("Ignoring 'MAX_UPLOAD_SIZE' value %q: not a number, using the default of %d\n", limitString, defaultLimit)
		})
		return defaultLimit
	}

	if limit > hardCapLimit {
		uploadSizeCapWarnOnce.Do(func() {
			log.Printf("Capping 'MAX_UPLOAD_SIZE' of %d to the maximum of %d\n", limit, hardCapLimit)
		})
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
