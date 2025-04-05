package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db := connectToPostgres()

	defer db.Close()

	createDeckTable(db)

	deck := Deck{
		Name:             "My First Deck",
		Description:      "This is a test deck",
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		LastStudyDate:    time.Time{},
		TotalCards:       0,
	}
	deckID := insertDeck(db, deck)
	log.Printf("Inserted deck with ID: %d", deckID)

	deck, err := getDeckDetails(db, 1)
	if err != nil {
		log.Fatalf("Error getting deck details: %v", err)
	}
	log.Printf("Deck ID: %+v", deck.ID)
	log.Printf("Deck Name: %+v", deck.Name)
	log.Printf("Deck Description: %+v", deck.Description)
	log.Printf("Deck Creation Date: %+v", deck.CreationDate)
	log.Printf("Deck Modification Date: %+v", deck.ModificationDate)
	log.Printf("Deck Last Study Date: %+v", deck.LastStudyDate)
	log.Printf("Deck Total Cards: %+v", deck.TotalCards)
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

type Deck struct {
	ID               int
	Name             string
	Description      string
	CreationDate     time.Time
	ModificationDate time.Time
	LastStudyDate    time.Time
	TotalCards       int
}

func insertDeck(db *sql.DB, deck Deck) int {
	query := `INSERT INTO decks (name, description) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, deck.Name, deck.Description).Scan(&deck.ID)

	if err != nil {
		log.Fatalf("Error inserting deck: %v", err)
	}

	return deck.ID
}

func getDeckDetails(db *sql.DB, deckID int) (Deck, error) {
	query := `SELECT id, name, description, creation_date, modification_date, last_study_date, total_cards FROM decks WHERE id = $1`
	var deck Deck

	var lastStudyDate sql.NullTime
	err := db.QueryRow(query, deckID).Scan(
		&deck.ID,
		&deck.Name,
		&deck.Description,
		&deck.CreationDate,
		&deck.ModificationDate, &lastStudyDate,
		&deck.TotalCards)

	if err != nil {
		return Deck{}, err
	}

	if lastStudyDate.Valid {
		deck.LastStudyDate = lastStudyDate.Time
	} else {
		deck.LastStudyDate = time.Time{}
	}

	return deck, nil
}
