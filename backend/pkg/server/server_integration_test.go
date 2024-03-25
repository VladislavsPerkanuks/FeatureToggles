package server

import (
	"log"
	"net/http"
	"os"
	"testing"
	"vladislavsperkanuks/feature-toggles/pkg/db"
)

var testServer *Server

const (
	port    = ":8082"
	baseURL = "http://localhost" + port
	apiURL  = baseURL + "/api/v1"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	db, closer, err := db.New("test.db")
	if err != nil {
		log.Fatalf("new repo: %s", err)
	}

	defer func() {
		if err := closer(); err != nil {
			log.Fatalf("close db: %s", err)
		}

		if err := os.Remove("test.db"); err != nil {
			log.Fatalf("remove test db: %s", err)
		}
	}()

	testServer = New(db)

	go testServer.Run(port)

	for {
		_, err := http.DefaultClient.Get(baseURL)
		if err == nil {
			break
		}
	}

	return m.Run()
}
