package database

import (
	"database/sql"
	"testing"

	"github.com/attic-labs/testify/mock"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Open(driverName string, dataSourceName string) (*sql.DB, error) {
	return &sql.DB{}, nil
}

func TestNewDeckDBWrapper(t *testing.T) {
	t.Run("Valid New DeckDBWrapper", func(t *testing.T) {
		mockDB := &MockDB{}
		db, err := mockDB.Open("sqlite3", ":memory:")

		if err != nil {
			t.Fatalf("Failed to open database: %v", err)
		}

		deckDBWrapper := NewDeckDBWrapper(db)

		assert.NotNil(t, deckDBWrapper)
		assert.IsType(t, &DeckDBWrapper{}, deckDBWrapper)
	})
}
