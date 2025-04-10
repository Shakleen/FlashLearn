package model

import "time"

type Card struct {
	ID               int       `json:"id"`
	DeckID           int       `json:"deck_id"`
	Content          string    `json:"content"`
	CreationTime     time.Time `json:"creation_time"`
	ModificationTime time.Time `json:"modification_time"`
	NextReviewTime   time.Time `json:"next_review_time"`
	RetentionLevel   int       `json:"retention_level"`
	Flag             int       `json:"flag"`
	Source           string    `json:"source"`
}

func NewCard(deckID int, content string, source string) Card {
	return Card{
		DeckID:           deckID,
		Content:          content,
		Source:           source,
		CreationTime:     time.Now(),
		ModificationTime: time.Now(),
		NextReviewTime:   time.Now().Add(time.Minute * 10),
		RetentionLevel:   0,
		Flag:             0,
	}
}
