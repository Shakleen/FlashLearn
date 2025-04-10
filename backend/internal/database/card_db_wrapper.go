package database

import (
	"database/sql"
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
	sb.WriteString(fmt.Sprintf("%s TIMESTAMP CHECK (%s >= NOW()), ", cardColumnNextReviewTime, cardColumnNextReviewTime))
	sb.WriteString(fmt.Sprintf("%s INT DEFAULT %d CHECK (%s >= %d), ", cardColumnRetentionLevel, cardMinRetentionLevel, cardColumnRetentionLevel, cardMinRetentionLevel))
	sb.WriteString(fmt.Sprintf("%s INT DEFAULT 0 CHECK (%s BETWEEN %d AND %d), ", cardColumnFlag, cardColumnFlag, cardMinFlag, cardMaxFlag))
	sb.WriteString(fmt.Sprintf("%s TEXT ", cardColumnSource))
	sb.WriteString(")")

	query := sb.String()
	return query
}
