package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/UTDNebula/nebula-api/api/common/log"

	_ "github.com/joho/godotenv/autoload"
)

func GetPortString() string {

	portNumber, exist := os.LookupEnv("Port")
	if !exist {
		portNumber = "8080"
	}

	portString := fmt.Sprintf(":%s", portNumber)

	return portString
}

func GetEnvMongoURI() string {

	uri, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		log.WriteErrorMsg("Error loading 'MONGODB_URI' from the .env file")
		os.Exit(1)
	}

	return uri
}

func GetEnvLogin() (netID string, password string) {

	netID, exist := os.LookupEnv("LOGIN_NETID")
	if !exist {
		log.WriteErrorMsg("Error loading 'LOGIN_NETID' from the .env file")
		os.Exit(1)
	}
	password, exist = os.LookupEnv("LOGIN_PASSWORD")
	if !exist {
		log.WriteErrorMsg("Error loading 'LOGIN_PASSWORD' from the .env file")
		os.Exit(1)
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
