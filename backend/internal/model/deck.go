package model

import "time"

type Deck struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
	LastStudyDate    time.Time `json:"last_study_date"`
	TotalCards       int       `json:"total_cards"`
}

func NewDeck(name string, description string) Deck {
	return Deck{
		Name:             name,
		Description:      description,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		LastStudyDate:    time.Time{},
		TotalCards:       0,
	}
}
