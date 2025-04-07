package api

import (
	"encoding/json"
	"flash-learn/internal/database"
	"flash-learn/internal/model"
	"flash-learn/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

const (
	GetSingleDeckInvalidDeckIDErrorMessage string = "Invalid deck ID"
	InsertDeckInvalidBodyErrorMessage      string = "Invalid request body"
	GetSingleDeckNotFoundErrorMessage      string = "Deck not found"
	InternalServerErrorMessage             string = "Internal server error"
)

type APIServer struct {
	address string
	db      database.DBWrapper
	server  *http.Server
}

// NewAPIServer creates a new instance of APIServer.
// It initializes the server with the given address and database wrapper.
// The address is the server's listening address, and the db is the database wrapper
// used for database operations.
func NewAPIServer(address string, db database.DBWrapper) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

// Start initializes the server and starts listening for incoming requests.
//
// Returns:
//   - error : An error if the server fails to start, nil otherwise.
func (s *APIServer) Start() error {
	if s.server != nil {
		return nil
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /deck/{id}", s.HandleGetSingleDeck)
	router.HandleFunc("GET /deck", s.HandleGetAllDecks)
	router.HandleFunc("GET /deck/count", s.HandleGetDeckCount)
	router.HandleFunc("POST /deck", s.HandleInsertDeck)
	s.server = &http.Server{
		Addr:    s.address,
		Handler: router,
	}
	return s.server.ListenAndServe()
}

// Stop stops the server if it is running.
//
// Returns:
//   - error : An error if the server fails to stop, nil otherwise.
func (s *APIServer) Stop() error {
	if s.server == nil {
		return nil
	}

	return s.server.Close()
}

// HandleGetSingleDeck handles the HTTP GET request for retrieving a single deck.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 400 Bad Request : If the deck ID is invalid or deck is not found.
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the deck is found and the request is successful.
func (s *APIServer) HandleGetSingleDeck(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idStr := strings.Split(r.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, GetSingleDeckInvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Fetch from database
	deck, dbErr := s.db.GetSingle(id)
	if dbErr != nil {
		if dbErr == utils.ErrRecordNotExist {
			http.Error(w, GetSingleDeckNotFoundErrorMessage, http.StatusBadRequest)
		} else {
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	encodingErr := json.NewEncoder(w).Encode(deck)
	if encodingErr != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleGetAllDecks handles the HTTP GET request for retrieving all decks.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the decks are found and the request is successful.
func (s *APIServer) HandleGetAllDecks(w http.ResponseWriter, r *http.Request) {
	// Fetch from database
	deckArray, err := s.db.GetAll()
	if err != nil {
		if err == utils.ErrDatabaseNotExist {
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deckArray)
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *APIServer) HandleGetDeckCount(w http.ResponseWriter, r *http.Request) {
	// Fetch from database
	count, err := s.db.GetCount()
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"count": count})
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *APIServer) HandleInsertDeck(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from request body
	type InsertInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var bodyInput InsertInput
	err := json.NewDecoder(r.Body).Decode(&bodyInput)
	if err != nil {
		http.Error(w, InsertDeckInvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Process input data
	bodyInput.Name = strings.TrimSpace(bodyInput.Name)
	bodyInput.Description = strings.TrimSpace(bodyInput.Description)
	if bodyInput.Name == "" {
		http.Error(w, InsertDeckInvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Insert into database
	deck := model.NewDeck(bodyInput.Name, bodyInput.Description)
	deckID, dbErr := s.db.Insert(deck)
	if dbErr != nil {
		if dbErr == utils.ErrMaxLengthExceeded {
			http.Error(w, InsertDeckInvalidBodyErrorMessage, http.StatusBadRequest)
		} else {
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"id": deckID})
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
