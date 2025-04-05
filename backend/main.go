package main

import (
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
		Name:             "My First Deck Edit3",
		Description:      "This is a test deck Edit3",
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
