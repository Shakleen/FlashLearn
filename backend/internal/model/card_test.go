package model

import (
	"testing"
	"time"

	"github.com/attic-labs/testify/assert"
)

func TestNewCard(t *testing.T) {
	type testCase struct {
		deckID  int
		content string
		source  string
	}

	t.Run("Valid New Card", func(t *testing.T) {
		testCases := []testCase{
			{deckID: 1, content: "Content", source: "Source"},
			{deckID: 2, content: "Content 2", source: "Source 2"},
			{deckID: 2, content: "Content 2", source: ""},
		}

		for _, test := range testCases {
			card := NewCard(test.deckID, test.content, test.source)

			assert.Equal(t, test.deckID, card.DeckID)
			assert.Equal(t, test.content, card.Content)
			assert.Equal(t, test.source, card.Source)

			assert.True(t, time.Since(card.CreationTime) < 1*time.Millisecond)
			assert.True(t, time.Since(card.ModificationTime) < 1*time.Millisecond)
			assert.True(t, card.NextReviewTime.After(card.CreationTime))
			assert.Equal(t, 0, card.RetentionLevel)
			assert.Equal(t, 0, card.Flag)
		}
	})
}
