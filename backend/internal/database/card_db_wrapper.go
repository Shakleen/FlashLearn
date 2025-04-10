package database

import "database/sql"

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
)

type CardDBWrapperInterface interface {
	CreateTable() error
}

type CardDBWrapper struct {
	db *sql.DB
}

func NewCardDBWrapper(db *sql.DB) *CardDBWrapper {
	return &CardDBWrapper{db: db}
}

func (c *CardDBWrapper) CreateTable() error {
	return nil
}
