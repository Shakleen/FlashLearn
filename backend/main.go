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
		ID:               1,
		Name:             "My First Deck Edit2",
		Description:      "This is a test deck Edit2",
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		LastStudyDate:    time.Time{},
		TotalCards:       0,
	}

	err := modifyDeckDetails(db, deck)
	if err != nil {
		log.Fatalf("Error modifying deck: %v", err)
	}

	deck2, err := getDeckDetails(db, deck.ID)
	if err != nil {
		log.Fatalf("Error getting deck details: %v", err)
	}
	log.Printf("Deck Name: %s", deck2.Name)
	log.Printf("Deck Description: %s", deck2.Description)
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
	} else if err == sql.ErrNoRows {
		log.Printf("No deck found with ID: %d", deckID)
		return Deck{}, err
	}

	if lastStudyDate.Valid {
		deck.LastStudyDate = lastStudyDate.Time
	} else {
		deck.LastStudyDate = time.Time{}
	}

	return deck, nil
}

func getDeckArray(db *sql.DB) ([]Deck, error) {
	query := `SELECT id, name, description, total_cards FROM decks`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		log.Printf("No decks found")
		return nil, err
	}

	defer rows.Close()

	var decks []Deck

	for rows.Next() {
		var deck Deck

		err := rows.Scan(
			&deck.ID,
			&deck.Name,
			&deck.Description,
			&deck.TotalCards)

		if err != nil {
			return nil, err
		}

		decks = append(decks, deck)
	}

	return decks, nil
}

func modifyDeckDetails(db *sql.DB, deck Deck) error {
	deck.ModificationDate = time.Now()
	query := `UPDATE decks SET name = $1, description = $2, modification_date = $3 WHERE id = $4`
	_, err := db.Exec(query, deck.Name, deck.Description, deck.ModificationDate, deck.ID)

	if err != nil {
		return err
	}

	return nil
}
