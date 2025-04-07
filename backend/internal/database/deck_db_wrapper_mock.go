package database

import (
	"flash-learn/internal/model"
	"flash-learn/internal/utils"
	"time"
)

type DeckDBWrapperMock struct {
	db    map[int]model.Deck
	index int
}

func NewDeckDBWrapperMock() *DeckDBWrapperMock {
	return &DeckDBWrapperMock{}
}

func (wrapper *DeckDBWrapperMock) CreateTable() error {
	if wrapper.db == nil {
		wrapper.db = make(map[int]model.Deck)
		wrapper.index = 0
	}

	return nil
}

func (wrapper *DeckDBWrapperMock) Insert(deck model.Deck) (int, error) {
	if wrapper.db == nil {
		return 0, utils.ErrDatabaseNotExist
	}

	wrapper.db[wrapper.index] = deck
	wrapper.index++
	return wrapper.index, nil
}

func (wrapper *DeckDBWrapperMock) GetSingle(deckID int) (model.Deck, error) {
	if wrapper.db == nil {
		return model.Deck{}, utils.ErrDatabaseNotExist
	}
	deck, exists := wrapper.db[deckID]

	if !exists {
		return model.Deck{}, utils.ErrRecordNotExist
	}

	return deck, nil
}

func (wrapper *DeckDBWrapperMock) GetAll() ([]model.Deck, error) {
	if wrapper.db == nil {
		return nil, utils.ErrDatabaseNotExist
	}

	decks := make([]model.Deck, 0, len(wrapper.db))

	for _, deck := range wrapper.db {
		decks = append(decks, deck)
	}

	return decks, nil
}

func (wrapper *DeckDBWrapperMock) Modify(deck model.Deck) error {
	if wrapper.db == nil {
		return utils.ErrDatabaseNotExist
	}

	oldDeck, exists := wrapper.db[deck.ID]

	if !exists {
		return utils.ErrRecordNotExist
	}

	oldDeck.Name = deck.Name
	oldDeck.Description = deck.Description
	oldDeck.ModificationDate = time.Now()

	wrapper.db[deck.ID] = oldDeck

	return nil
}

func (wrapper *DeckDBWrapperMock) Delete(id int) error {
	if wrapper.db == nil {
		return utils.ErrDatabaseNotExist
	}

	if _, exists := wrapper.db[id]; !exists {
		return utils.ErrRecordNotExist
	}

	delete(wrapper.db, id)

	return nil
}
