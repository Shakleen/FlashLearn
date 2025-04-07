package database

import (
	"database/sql"
	"flash-learn/internal/model"
	"flash-learn/internal/utils"
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

// An interface that defines the methods for interacting with the database.
// This interface abstracts the database operations for decks,
// allowing for easier testing and mocking.
type DBWrapper interface {
	CreateTable() error
	Insert(deck model.Deck) (int, error)
	GetSingle(deckID int) (model.Deck, error)
	GetCount() (int, error)
	GetAll() ([]model.Deck, error)
	Modify(deck model.Deck) error
	Delete(id int) error
}

type DeckDBWrapper struct {
	db *sql.DB
}

// Creates and returns a new instance of DeckDBWrapper.
//
// Parameters:
//   - db *sql.DB : The database connection.
//
// Returns:
//   - *DeckDBWrapper
func NewDeckDBWrapper(db *sql.DB) *DeckDBWrapper {
	return &DeckDBWrapper{db: db}
}

// Creates a new table in the database if it doesn't already exist.
//
// Returns:
//   - error : An error if the table creation fails, nil otherwise.
func (wrapper *DeckDBWrapper) CreateTable() error {
	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return utils.ErrDatabaseNotExist
	}

	query := buildCreateTableQueryString()

	_, err := wrapper.db.Exec(query)

	if err != nil {
		log.Fatalf("Error creating decks table: %v", err)
	}

	return err
}

// A helper function that constructs the SQL query string
// to create the decks table.
//
// Returns:
//   - string : The SQL query string to create the decks table.
func buildCreateTableQueryString() string {
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(deckTableName)
	sb.WriteString(" (")
	sb.WriteString(fmt.Sprintf("%s SERIAL PRIMARY KEY, ", deckColumnID))
	sb.WriteString(fmt.Sprintf("%s VARCHAR(%d) NOT NULL UNIQUE, ", deckColumnName, deckColumnNameMaxLength))
	sb.WriteString(fmt.Sprintf("%s VARCHAR(%d) NOT NULL, ", deckColumnDescription, deckColumnDescriptionMaxLength))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW(), ", deckColumnCreationDate))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW(), ", deckColumnModificationDate))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP, ", deckColumnLastStudyDate))
	sb.WriteString(fmt.Sprintf("%s INT DEFAULT 0", deckColumnTotalCards))
	sb.WriteString(")")

	query := sb.String()
	return query
}

// Inserts a new deck into the database and returns its unique ID.
//
// Parameters:
//   - deck model.Deck : Details of the deck to be inserted as a model.Deck object.
//
// Returns:
//   - int : The unique ID of the inserted deck.
//   - error : An error if the insertion fails, nil otherwise.
func (wrapper *DeckDBWrapper) Insert(deck model.Deck) (int, error) {
	if len(deck.Name) > deckColumnNameMaxLength || len(deck.Description) > deckColumnDescriptionMaxLength {
		log.Fatal("Deck name or description exceeds maximum length")
		return -1, utils.ErrMaxLengthExceeded
	}

	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return -1, utils.ErrDatabaseNotExist
	}

	query := buildInsertQueryString()
	err := wrapper.db.QueryRow(query, deck.Name, deck.Description).Scan(&deck.ID)

	if err != nil {
		log.Fatalf("Error inserting deck: %v", err)
		return -1, err
	}

	return deck.ID, nil
}

// A helper function that constructs the SQL query string to insert a new deck.
//
// Returns:
//   - string : The SQL query string to insert a new deck.
func buildInsertQueryString() string {
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
	return query
}

// Retrieves a single deck from the database based on its unique ID.
//
// Parameters:
//   - deckID int : The unique ID of the deck to be retrieved.
//
// Returns:
//   - model.Deck : The details of the retrieved deck as a model.Deck object.
//   - error : An error if the retrieval fails, nil otherwise.
func (wrapper *DeckDBWrapper) GetSingle(deckID int) (model.Deck, error) {
	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return model.Deck{}, utils.ErrDatabaseNotExist
	}

	var deck model.Deck
	var lastStudyDate sql.NullTime

	query := buildGetSingleQueryString()

	err := wrapper.db.QueryRow(query, deckID).Scan(
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

// Helper function that constructs the SQL query string to retrieve a single deck.
//
// Returns:
//   - string : The SQL query string to retrieve a single deck.
func buildGetSingleQueryString() string {
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
	return query
}

// Retrieves all decks from the database.
//
// Returns:
//   - []model.Deck : A slice of model.Deck objects representing all decks.
//   - error : An error if the retrieval fails, nil otherwise.
func (wrapper *DeckDBWrapper) GetAll() ([]model.Deck, error) {
	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return nil, utils.ErrDatabaseNotExist
	}

	query := buildGetAllQueryString()
	rows, err := wrapper.db.Query(query)
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

// Counts the total number of decks in the database.
//
// Returns:
//   - int : The total number of decks.
//   - error : An error if the count fails, nil otherwise.
func (wrapper *DeckDBWrapper) GetCount() (int, error) {
	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return 0, utils.ErrDatabaseNotExist
	}

	query := buildGetCountQueryString()
	var count int

	err := wrapper.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatalf("Error getting deck count: %v", err)
		return -1, err
	}

	return count, nil
}

// Helper function that constructs the SQL query string to count all decks.
//
// Returns:
//   - string : The SQL query string to count all decks.
func buildGetCountQueryString() string {
	var sb strings.Builder
	sb.WriteString("SELECT COUNT(")
	sb.WriteString(deckColumnID)
	sb.WriteString(") FROM ")
	sb.WriteString(deckTableName)

	query := sb.String()
	return query
}

// Helper function that constructs the SQL query string to retrieve all decks.
//
// Returns:
//   - string : The SQL query string to retrieve all decks.
func buildGetAllQueryString() string {
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
	return query
}

// Modifies name and/or description of an existing deck in the database.
//
// Parameters:
//   - deck model.Deck : The deck object containing the new name and/or description.
//
// Returns:
//   - error : An error if the modification fails, nil otherwise.
func (wrapper *DeckDBWrapper) Modify(deck model.Deck) error {
	if len(deck.Name) > deckColumnNameMaxLength || len(deck.Description) > deckColumnDescriptionMaxLength {
		log.Fatal("Deck name or description exceeds maximum length")
		return utils.ErrMaxLengthExceeded
	}

	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return utils.ErrDatabaseNotExist
	}

	// Modification date is set to the current time
	deck.ModificationDate = time.Now()

	query := buildModifyQueryString()
	_, err := wrapper.db.Exec(query, deck.Name, deck.Description, deck.ModificationDate, deck.ID)

	if err != nil {
		log.Fatalf("Error modifying deck: %v", err)
		return err
	}

	return nil
}

// Helper function that constructs the SQL query string to modify a deck.
//
// Returns:
//   - string : The SQL query string to modify a deck.
func buildModifyQueryString() string {
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
	return query
}

// Deletes a deck from the database based on its unique ID.
//
// Parameters:
//   - id int : The unique ID of the deck to be deleted.
//
// Returns:
//   - error : An error if the deletion fails, nil otherwise.
func (wrapper *DeckDBWrapper) Delete(id int) error {
	if wrapper.db == nil {
		log.Fatal("Database connection is nil")
		return utils.ErrDatabaseNotExist
	}

	query := buildDeleteQueryString()
	_, err := wrapper.db.Exec(query, id)
	if err != nil {
		log.Fatalf("Error deleting deck: %v", err)
		return err
	}

	return nil
}

// Helper function that constructs the SQL query string to delete a deck.
//
// Returns:
//   - string : The SQL query string to delete a deck.
func buildDeleteQueryString() string {
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	sb.WriteString(deckTableName)
	sb.WriteString(" WHERE ")
	sb.WriteString(deckColumnID)
	sb.WriteString(" = $1")

	query := sb.String()
	return query
}
