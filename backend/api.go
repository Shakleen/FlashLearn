package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /deck/{id}", s.handleDeck)
	router.HandleFunc("GET /deck", s.handleDeckArray)
	router.HandleFunc("POST /deck", s.handleInsertDeck)
	router.HandleFunc("POST /deck/update", s.handleUpdateDeck)

	server := &http.Server{
		Addr:    s.address,
		Handler: router,
	}

	return server.ListenAndServe()
}

func (s *APIServer) handleDeck(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := getDeckDetails(s.db, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deck)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *APIServer) handleDeckArray(w http.ResponseWriter, r *http.Request) {
	deckArray, err := getDeckArray(s.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deckArray)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

type DeckInput struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *APIServer) handleInsertDeck(w http.ResponseWriter, r *http.Request) {
	var deckInput DeckInput
	err := json.NewDecoder(r.Body).Decode(&deckInput)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Deck name:", deckInput.Name)
	fmt.Println("Deck description:", deckInput.Description)

	deck := NewDeck(deckInput.Name, deckInput.Description)
	deckID := insertDeck(s.db, deck)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]int{"deck_id": deckID}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	fmt.Println("Deck inserted with ID:", deckID)
}

func (s *APIServer) handleUpdateDeck(w http.ResponseWriter, r *http.Request) {
	var deckInput DeckInput
	err2 := json.NewDecoder(r.Body).Decode(&deckInput)
	if err2 != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Deck ID:", deckInput.ID)
	fmt.Println("Deck name:", deckInput.Name)
	fmt.Println("Deck description:", deckInput.Description)

	deck := NewDeck(deckInput.Name, deckInput.Description)
	deck.ID = deckInput.ID

	err3 := modifyDeckDetails(s.db, deck)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Deck updated successfully"}
	err4 := json.NewEncoder(w).Encode(response)
	if err4 != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println("Deck updated with ID:", deckInput.ID)
}
