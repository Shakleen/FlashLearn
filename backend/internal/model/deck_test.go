package model

import (
	"testing"
	"time"

	"github.com/attic-labs/testify/assert"
)

type MockTime struct {
	MockedTime time.Time
}

func (mt *MockTime) Now() time.Time {
	return mt.MockedTime
}

func TestNewDeck(t *testing.T) {
	type testCase struct {
		name        string
		description string
	}

	t.Run("Valid New Deck", func(t *testing.T) {
		testCases := []testCase{
			{name: "Name", description: "Test Description"},
			{name: "Name 2", description: "Test Description 2"},
		}

		for _, test := range testCases {
			timeNow := time.Now()
			deck := NewDeck(test.name, test.description)

			assert.Equal(t, test.name, deck.Name)
			assert.Equal(t, test.description, deck.Description)
			assert.Equal(t, 0, deck.TotalCards)

			assert.IsType(t, time.Time{}, deck.CreationDate)
			timeDiff := timeNow.Sub(deck.CreationDate)
			assert.True(t, timeDiff < 1*time.Millisecond, "CreationDate should be close to the current time")

			assert.IsType(t, time.Time{}, deck.ModificationDate)
			timeDiff = timeNow.Sub(deck.ModificationDate)
			assert.True(t, timeDiff < 1*time.Millisecond, "ModificationDate should be close to the current time")

			assert.IsType(t, time.Time{}, deck.LastStudyDate)
			timeDiff = timeNow.Sub(deck.LastStudyDate)
			assert.False(t, timeDiff < 1*time.Millisecond, "LastStudyDate shouldn't be close to the current time")
		}
	})
}
