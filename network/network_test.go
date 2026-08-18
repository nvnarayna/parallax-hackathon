package network

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"wifi-diagnostic/database"
)

type testEngine struct {
	running bool
}

func TestNetwork(t *testing.T) {
	const databasePath = "network-test.db"

	_ = os.Remove(databasePath)
	defer os.Remove(databasePath)

	db, err := database.Open(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	logger := log.New(
		os.Stdout,
		"[network-test] ",
		log.LstdFlags,
	)

	_ = context.Background()

	_ = logger
	_ = http.MethodGet

	_ = httptest.NewRecorder

	_ = db
}