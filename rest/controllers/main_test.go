package controllers

import (
	"fmt"
	"os"
	"testing"

	"github.com/UTDNebula/nebula-api/rest/configs"
)

func TestMain(m *testing.M) {
	uri, err := configs.GetEnvMongoURI()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := configs.ConnectDB(uri); err != nil {
		fmt.Fprintf(os.Stderr, "connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	InitializeCollections()
	os.Exit(m.Run())
}
