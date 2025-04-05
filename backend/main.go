package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db := connectToPostgres()

	defer db.Close()

	createDeckTable(db)
}

func createDeckTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS decks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(64) NOT NULL,
		description VARCHAR(255) NOT NULL,
		creation_date TIMESTAMP DEFAULT NOW(),
		modification_date TIMESTAMP DEFAULT NOW(),
		last_study_date TIMESTAMP,
		total_cards INT DEFAULT 0
	)`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Error creating decks table: %v", err)
	}
}
