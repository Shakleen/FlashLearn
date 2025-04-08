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
	InvalidDeckIDErrorMessage         string = "Invalid deck ID"
	InvalidBodyErrorMessage           string = "Invalid request body"
	GetSingleDeckNotFoundErrorMessage string = "Deck not found"
	InternalServerErrorMessage        string = "Internal server error"
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
	router.HandleFunc("GET /deck/nameMaxLength", s.HandleGetDeckNameMaxLength)
	router.HandleFunc("GET /deck/descriptionMaxLength", s.HandleGetDeckDescriptionMaxLength)
	router.HandleFunc("POST /deck", s.HandleInsertDeck)
	router.HandleFunc("POST /deck/{id}", s.HandleModifyDeck)
	router.HandleFunc("DELETE /deck/{id}", s.HandleDeleteDeck)

	corsHandler := corsMiddleware(router)

	s.server = &http.Server{
		Addr:    s.address,
		Handler: corsHandler,
	}
	return s.server.ListenAndServe()
}

// corsMiddleware adds CORS headers to the HTTP responses.
//
// Parameters:
//   - handler http.Handler
//
// Returns:
//   - http.Handler
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass control to the next handler
		handler.ServeHTTP(w, r)
	})
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
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
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

	type DeckOutput struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	deckOutputArray := make([]DeckOutput, len(deckArray))
	for i, deck := range deckArray {
		deckOutputArray[i] = DeckOutput{
			ID:          strconv.Itoa(deck.ID),
			Name:        deck.Name,
			Description: deck.Description,
		}
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deckOutputArray)
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleGetDeckCount handles the HTTP GET request for retrieving the count of decks.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the decks are found and the request is successful.
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

// HandleInsertDeck handles the HTTP POST request for inserting a new deck.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 400 Bad Request : If the request body is invalid or violates max length constraints.
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the decks are found and the request is successful.
func (s *APIServer) HandleInsertDeck(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from request body
	type InsertInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var bodyInput InsertInput
	err := json.NewDecoder(r.Body).Decode(&bodyInput)
	if err != nil {
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Process input data
	bodyInput.Name = strings.TrimSpace(bodyInput.Name)
	bodyInput.Description = strings.TrimSpace(bodyInput.Description)
	if bodyInput.Name == "" {
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Insert into database
	deck := model.NewDeck(bodyInput.Name, bodyInput.Description)
	deckID, dbErr := s.db.Insert(deck)
	if dbErr != nil {
		if dbErr == utils.ErrMaxLengthExceeded {
			http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
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

// HandleModifyDeck handles the HTTP POST request for modifying an existing deck.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 400 Bad Request : If id is invalid or the request body is invalid or violates max length constraints.
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the decks are found and the request is successful.
func (s *APIServer) HandleModifyDeck(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idStr := strings.Split(r.URL.Path, "/")[2]
	deckID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Parse JSON data from request body
	type InsertInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var bodyInput InsertInput
	err = json.NewDecoder(r.Body).Decode(&bodyInput)
	if err != nil {
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Process input data
	bodyInput.Name = strings.TrimSpace(bodyInput.Name)
	bodyInput.Description = strings.TrimSpace(bodyInput.Description)
	if bodyInput.Name == "" {
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Modify row in database
	deck := model.NewDeck(bodyInput.Name, bodyInput.Description)
	deck.ID = deckID
	dbErr := s.db.Modify(deck)
	if dbErr != nil {
		if dbErr == utils.ErrMaxLengthExceeded {
			http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		} else if dbErr == utils.ErrRecordNotExist {
			http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
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
	w.WriteHeader(http.StatusOK)
}

// HandleDeleteDeck handles the HTTP GET request for deleting a single deck.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 400 Bad Request : If the deck ID is invalid or deck is not found.
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the deck is found and the request is successful.
func (s *APIServer) HandleDeleteDeck(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idStr := strings.Split(r.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Fetch from database
	dbErr := s.db.Delete(id)
	if dbErr != nil {
		if dbErr == utils.ErrRecordNotExist {
			http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		} else {
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	encodingErr := json.NewEncoder(w).Encode(map[string]int{"id": id})
	if encodingErr != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleGetDeckNameMaxLength handles the HTTP GET request for retrieving the max length of the deck name.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
func (s *APIServer) HandleGetDeckNameMaxLength(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]int{"maxLength": database.DeckColumnNameMaxLength})
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleGetDeckDescriptionMaxLength handles the HTTP GET request for retrieving the max length of the deck description.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
func (s *APIServer) HandleGetDeckDescriptionMaxLength(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]int{"maxLength": database.DeckColumnDescriptionMaxLength})
	if err != nil {
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
