package api

import (
	"encoding/json"
	"flash-learn/internal/database"
	"flash-learn/internal/model"
	"fmt"
	"net/http"
	"strconv"
)

type APIServer struct {
	address     string
	ddb_wrapper database.DBWrapper
}

func NewAPIServer(address string, ddb_wrapper database.DBWrapper) *APIServer {
	return &APIServer{
		address:     address,
		ddb_wrapper: ddb_wrapper,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /deck/{id}", s.handleDeck)
	router.HandleFunc("GET /deck", s.handleDeckArray)
	router.HandleFunc("POST /deck", s.handleInsertDeck)
	router.HandleFunc("POST /deck/{id}", s.handleUpdateDeck)
	router.HandleFunc("DELETE /deck/{id}", s.handleDeleteDeck)

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

	deck, err := s.ddb_wrapper.GetSingle(id)

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
	deckArray, err := s.ddb_wrapper.GetAll()

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

	deck := model.NewDeck(deckInput.Name, deckInput.Description)
	deckID, err := s.ddb_wrapper.Insert(deck)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	var deckInput DeckInput
	err2 := json.NewDecoder(r.Body).Decode(&deckInput)
	if err2 != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Deck ID:", id)
	fmt.Println("Deck name:", deckInput.Name)
	fmt.Println("Deck description:", deckInput.Description)

	deck := model.NewDeck(deckInput.Name, deckInput.Description)
	deck.ID = id

	err3 := s.ddb_wrapper.Modify(deck)
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
	fmt.Println("Deck updated with ID:", id)
}

func (s *APIServer) handleDeleteDeck(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	err = s.ddb_wrapper.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Deck deleted successfully"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println("Deck deleted with ID:", id)
}
