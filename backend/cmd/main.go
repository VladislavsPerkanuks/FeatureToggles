package main

import (
	"log"
	"vladislavsperkanuks/feature-toggles/pkg/db"
	"vladislavsperkanuks/feature-toggles/pkg/server"
)

func main() {
	db, closer, err := db.New("../db.sqlite3")
	if err != nil {
		log.Fatalf("new repo: %s", err)
	}

	defer func() {
		if err := closer(); err != nil {
			log.Fatalf("close db: %s", err)
		}
	}()

	if err := server.New(db).Run("0.0.0.0:8081"); err != nil {
		log.Fatalf("run server: %s", err)
	}
}
