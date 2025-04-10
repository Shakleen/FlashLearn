package database

import (
	"flash-learn/internal/model"
	"flash-learn/internal/utils"
)

type CardDBWrapperMock struct {
	db    map[int]map[int]model.Card
	index map[int]int
}

func NewCardDBWrapperMock() *CardDBWrapperMock {
	return &CardDBWrapperMock{}
}

func (wrapper *CardDBWrapperMock) CreateTable() error {
	if wrapper.db == nil {
		wrapper.db = make(map[int]map[int]model.Card)
	}
	wrapper.index = make(map[int]int)

	return nil
}

func (wrapper *CardDBWrapperMock) Insert(card model.Card) (int, error) {
	_, ok := wrapper.db[card.DeckID]

	if !ok {
		return -1, utils.ErrDeckNotExist
	}

	wrapper.db[card.DeckID][wrapper.index[card.DeckID]] = card
	wrapper.index[card.DeckID]++

	return wrapper.index[card.DeckID] - 1, nil
}

func (wrapper *CardDBWrapperMock) InsertDeck(deckID int) {
	wrapper.db[deckID] = make(map[int]model.Card)
}

func (wrapper *CardDBWrapperMock) GetTotalCards(deckID int) (int, error) {
	_, ok := wrapper.db[deckID]
	if !ok {
		return 0, nil
	}
	return len(wrapper.db[deckID]), nil
}
