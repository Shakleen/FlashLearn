package main

import (
	"database/sql"
	"flash-learn/internal/model"
	"fmt"
	"log"
	"strings"
	"time"
)

const deckTableName string = "decks"
const deckColumnID string = "id"
const deckColumnName string = "name"
const deckColumnDescription string = "description"
const deckColumnCreationDate string = "creation_date"
const deckColumnModificationDate string = "modification_date"
const deckColumnLastStudyDate string = "last_study_date"
const deckColumnTotalCards string = "total_cards"
const deckColumnNameMaxLength int = 64
const deckColumnDescriptionMaxLength int = 255

func createDeckTable(db *sql.DB) {
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(deckTableName)
	sb.WriteString(" (")
	sb.WriteString(fmt.Sprintf("%s SERIAL PRIMARY KEY, ", deckColumnID))
	sb.WriteString(fmt.Sprintf("%s VARCHAR(%d) NOT NULL, ", deckColumnName, deckColumnNameMaxLength))
	sb.WriteString(fmt.Sprintf("%s VARCHAR(%d) NOT NULL, ", deckColumnDescription, deckColumnDescriptionMaxLength))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW(), ", deckColumnCreationDate))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW(), ", deckColumnModificationDate))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP, ", deckColumnLastStudyDate))
	sb.WriteString(fmt.Sprintf("%s INT DEFAULT 0", deckColumnTotalCards))
	sb.WriteString(")")

	query := sb.String()

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Error creating decks table: %v", err)
	}
}

func insertDeck(db *sql.DB, deck model.Deck) int {
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(deckTableName)
	sb.WriteString(" (")
	sb.WriteString(deckColumnName)
	sb.WriteString(", ")
	sb.WriteString(deckColumnDescription)
	sb.WriteString(") VALUES ($1, $2) RETURNING ")
	sb.WriteString(deckColumnID)

	query := sb.String()
	err := db.QueryRow(query, deck.Name, deck.Description).Scan(&deck.ID)

	if err != nil {
		log.Fatalf("Error inserting deck: %v", err)
	}

	return deck.ID
}

func getDeckDetails(db *sql.DB, deckID int) (model.Deck, error) {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString(deckColumnID)
	sb.WriteString(", ")
	sb.WriteString(deckColumnName)
	sb.WriteString(", ")
	sb.WriteString(deckColumnDescription)
	sb.WriteString(", ")
	sb.WriteString(deckColumnCreationDate)
	sb.WriteString(", ")
	sb.WriteString(deckColumnModificationDate)
	sb.WriteString(", ")
	sb.WriteString(deckColumnLastStudyDate)
	sb.WriteString(", ")
	sb.WriteString(deckColumnTotalCards)
	sb.WriteString(" FROM ")
	sb.WriteString(deckTableName)
	sb.WriteString(" WHERE ")
	sb.WriteString(deckColumnID)
	sb.WriteString(" = $1")

	query := sb.String()
	var deck model.Deck

	var lastStudyDate sql.NullTime
	err := db.QueryRow(query, deckID).Scan(
		&deck.ID,
		&deck.Name,
		&deck.Description,
		&deck.CreationDate,
		&deck.ModificationDate, &lastStudyDate,
		&deck.TotalCards)

	if err != nil {
		return model.Deck{}, err
	} else if err == sql.ErrNoRows {
		log.Printf("No deck found with ID: %d", deckID)
		return model.Deck{}, err
	}

	if lastStudyDate.Valid {
		deck.LastStudyDate = lastStudyDate.Time
	} else {
		deck.LastStudyDate = time.Time{}
	}

	return deck, nil
}

func getDeckArray(db *sql.DB) ([]model.Deck, error) {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString(deckColumnID)
	sb.WriteString(", ")
	sb.WriteString(deckColumnName)
	sb.WriteString(", ")
	sb.WriteString(deckColumnDescription)
	sb.WriteString(", ")
	sb.WriteString(deckColumnTotalCards)
	sb.WriteString(" FROM ")
	sb.WriteString(deckTableName)

	query := sb.String()
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		log.Printf("No decks found")
		return nil, err
	}

	defer rows.Close()

	var decks []model.Deck

	for rows.Next() {
		var deck model.Deck

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

func modifyDeckDetails(db *sql.DB, deck model.Deck) error {
	deck.ModificationDate = time.Now()
	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(deckTableName)
	sb.WriteString(" SET ")
	sb.WriteString(deckColumnName)
	sb.WriteString(" = $1, ")
	sb.WriteString(deckColumnDescription)
	sb.WriteString(" = $2, ")
	sb.WriteString(deckColumnModificationDate)
	sb.WriteString(" = $3 WHERE ")
	sb.WriteString(deckColumnID)
	sb.WriteString(" = $4")

	query := sb.String()
	_, err := db.Exec(query, deck.Name, deck.Description, deck.ModificationDate, deck.ID)

	if err != nil {
		return err
	}

	return nil
}

func deleteDeck(db *sql.DB, id int) error {
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	sb.WriteString(deckTableName)
	sb.WriteString(" WHERE ")
	sb.WriteString(deckColumnID)
	sb.WriteString(" = $1")

	query := sb.String()
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
