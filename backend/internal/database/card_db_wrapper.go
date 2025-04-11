package database

import (
	"database/sql"
	"flash-learn/internal/model"
	"flash-learn/internal/utils"
	"fmt"
	"log/slog"
	"strings"
)

const (
	cardTableName              = "cards"
	cardColumnID               = "id"
	cardColumnDeckID           = "deck_id"
	cardColumnContent          = "content"
	cardColumnCreationTime     = "creation_time"
	cardColumnModificationTime = "modification_time"
	cardColumnNextReviewTime   = "next_review_time"
	cardColumnRetentionLevel   = "retention_level"
	cardColumnFlag             = "flag"
	cardColumnSource           = "source"
	cardMinRetentionLevel      = 0
	cardMinFlag                = 0
	cardMaxFlag                = 9
)

// An interface that defines the methods for interacting with the card database.
// This interface abstracts the database operations for cards,
// allowing for easier testing and mocking.
type CardDBWrapperInterface interface {
	CreateTable() error
	Insert(card model.Card) (int, error)
	GetTotalCards(deckID int) (int, error)
}

// A struct that implements the CardDBWrapperInterface.
//
// This is the concrete implementation and should be used for actual
// database operations.
//
// Parameters:
//   - db *sql.DB : The database connection.
//
// Returns:
//   - *CardDBWrapper
type CardDBWrapper struct {
	db *sql.DB
}

// Creates and returns a new instance of CardDBWrapper.
//
// Parameters:
//   - db *sql.DB : The database connection.
//
// Returns:
//   - *CardDBWrapper
func NewCardDBWrapper(db *sql.DB) *CardDBWrapper {
	return &CardDBWrapper{db: db}
}

// Creates a new table in the database if it doesn't already exist.
//
// Returns:
//   - error : An error if the table creation fails, nil otherwise.

func (wrapper *CardDBWrapper) CreateTable() error {
	if wrapper.db == nil {
		slog.Error("Database connection is nil")
		return utils.ErrDatabaseNotExist
	}

	query := wrapper.buildCreateTableQueryString()
	slog.Debug("Creating cards table", "query", query)

	_, err := wrapper.db.Exec(query)

	if err != nil {
		slog.Error(fmt.Sprintf("Error creating cards table: %s", err))
	}

	return err
}

// A helper function that constructs the SQL query string
// to create the cards table.
//
// Returns:
//   - string : The SQL query string to create the cards table.
func (wrapper *CardDBWrapper) buildCreateTableQueryString() string {
	var sb strings.Builder

	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(cardTableName)
	sb.WriteString(" (")
	sb.WriteString(fmt.Sprintf("%s SERIAL PRIMARY KEY, ", cardColumnID))
	sb.WriteString(fmt.Sprintf("%s INT NOT NULL REFERENCES %s(%s), ", cardColumnDeckID, deckTableName, deckColumnID))
	sb.WriteString(fmt.Sprintf("%s TEXT NOT NULL, ", cardColumnContent))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW(), ", cardColumnCreationTime))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW() CHECK (%s >= %s), ", cardColumnModificationTime, cardColumnModificationTime, cardColumnCreationTime))
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP DEFAULT NOW() + INTERVAL '10 minutes' CHECK (%s >= NOW()), ", cardColumnNextReviewTime, cardColumnNextReviewTime))
	sb.WriteString(fmt.Sprintf("%s INT DEFAULT %d CHECK (%s >= %d), ", cardColumnRetentionLevel, cardMinRetentionLevel, cardColumnRetentionLevel, cardMinRetentionLevel))
	sb.WriteString(fmt.Sprintf("%s INT DEFAULT 0 CHECK (%s BETWEEN %d AND %d), ", cardColumnFlag, cardColumnFlag, cardMinFlag, cardMaxFlag))
	sb.WriteString(fmt.Sprintf("%s TEXT ", cardColumnSource))
	sb.WriteString(")")

	query := sb.String()
	return query
}

// Inserts a new card into the database and returns its unique ID.
//
// Parameters:
//   - card model.Card : Details of the card to be inserted as a model.Card object.
//
// Returns:
//   - int : The unique ID of the inserted card.
//   - error : An error if the insertion fails, nil otherwise.
func (wrapper *CardDBWrapper) Insert(card model.Card) (int, error) {
	query := wrapper.buildInsertQueryString(card)
	slog.Debug(fmt.Sprintf("Inserting card: %s", query))

	err := wrapper.db.QueryRow(query, card.DeckID, card.Content, card.Source).Scan(&card.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Error inserting card: %s", err))
		return 0, err
	}

	return card.ID, nil
}

// A helper function that constructs the SQL query string
// to insert a new card into the database.
//
// Returns:
//   - string : The SQL query string to insert a new card into the database.
func (wrapper *CardDBWrapper) buildInsertQueryString(card model.Card) string {
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString(cardTableName)
	sb.WriteString(" (")
	sb.WriteString(fmt.Sprintf("%s, %s, %s", cardColumnDeckID, cardColumnContent, cardColumnSource))
	sb.WriteString(") VALUES ($1, $2, $3) RETURNING ")
	sb.WriteString(cardColumnID)

	query := sb.String()
	return query
}

func (wrapper *CardDBWrapper) GetTotalCards(deckID int) (int, error) {
	query := wrapper.buildGetTotalCardsQueryString(deckID)
	slog.Debug(fmt.Sprintf("Getting total cards: %s", query))

	var count int
	err := wrapper.db.QueryRow(query, deckID).Scan(&count)
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting total cards: %s", err))
		return 0, err
	}
	return count, nil
}

func (wrapper *CardDBWrapper) buildGetTotalCardsQueryString(deckID int) string {
	var sb strings.Builder
	sb.WriteString("SELECT COUNT( ")
	sb.WriteString(cardColumnID)
	sb.WriteString(") FROM ")
	sb.WriteString(cardTableName)
	sb.WriteString(" WHERE ")
	sb.WriteString(cardColumnDeckID)
	sb.WriteString(" = $1")

	query := sb.String()
	return query
}
