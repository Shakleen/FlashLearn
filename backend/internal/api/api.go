package api

import (
	"encoding/json"
	"flash-learn/internal/database"
	"flash-learn/internal/model"
	"flash-learn/internal/utils"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const (
	InvalidDeckIDErrorMessage         string = "Invalid deck ID"
	InvalidBodyErrorMessage           string = "Invalid request body"
	GetSingleDeckNotFoundErrorMessage string = "Deck not found"
	InternalServerErrorMessage        string = "Internal server error"
	DuplicateKeyViolationErrorMessage string = "Duplicate key violation"
)

type APIServer struct {
	address string
	deck_db database.DBWrapper
	card_db database.CardDBWrapperInterface
	server  *http.Server
}

// NewAPIServer creates a new instance of APIServer.
// It initializes the server with the given address and database wrapper.
// The address is the server's listening address, and the db is the database wrapper
// used for database operations.
func NewAPIServer(address string, deck_db database.DBWrapper, card_db database.CardDBWrapperInterface) *APIServer {
	return &APIServer{
		address: address,
		deck_db: deck_db,
		card_db: card_db,
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

	slog.Debug("Creating router")
	router := http.NewServeMux()
	addRoutes(router, s)

	slog.Debug("Creating cors handler")
	corsHandler := corsMiddleware(router)

	s.server = &http.Server{
		Addr:    s.address,
		Handler: corsHandler,
	}

	slog.Debug("Starting server")
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
		slog.Debug(fmt.Sprintf("Invalid deck ID %s", idStr))
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Fetch from database
	deck, dbErr := s.deck_db.GetSingle(id)
	if dbErr != nil {
		if dbErr == utils.ErrRecordNotExist {
			slog.Debug("Deck not found", "error", dbErr)
			http.Error(w, GetSingleDeckNotFoundErrorMessage, http.StatusBadRequest)
		} else {
			slog.Debug("Error getting single deck", "error", dbErr)
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
	slog.Debug("Sent response", "deck", deck)
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
	deckArray, err := s.deck_db.GetAll()
	if err != nil {
		slog.Debug("Error getting all decks", "error", err)
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
	slog.Debug("Sent response", "deck", deckOutputArray)
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
	count, err := s.deck_db.GetCount()
	if err != nil {
		slog.Debug("Error getting deck count", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"count": count})
	if err != nil {
		slog.Debug("Error encoding deck count", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Sent response", "deck count", count)
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
		slog.Debug("Error decoding request body", "error", err)
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Process input data
	bodyInput.Name = strings.TrimSpace(bodyInput.Name)
	bodyInput.Description = strings.TrimSpace(bodyInput.Description)
	if bodyInput.Name == "" {
		slog.Debug("Invalid body input", "error", err)
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Insert into database
	deck := model.NewDeck(bodyInput.Name, bodyInput.Description)
	deckID, dbErr := s.deck_db.Insert(deck)
	if dbErr != nil {
		if dbErr == utils.ErrMaxLengthExceeded {
			slog.Debug("Max length exceeded", "error", dbErr)
			http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		} else if dbErr == utils.ErrDuplicateKeyViolation {
			slog.Debug("Duplicate key violation", "error", dbErr)
			http.Error(w, DuplicateKeyViolationErrorMessage, http.StatusConflict)
		} else {
			slog.Debug("Error inserting deck", "error", dbErr)
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"id": deckID})
	if err != nil {
		slog.Debug("Error encoding deck ID", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	slog.Debug("Sent response", "deck ID", deckID)
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
		slog.Debug(fmt.Sprintf("Invalid deck ID %s", idStr))
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
		slog.Debug("Error decoding request body", "error", err)
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Process input data
	bodyInput.Name = strings.TrimSpace(bodyInput.Name)
	bodyInput.Description = strings.TrimSpace(bodyInput.Description)
	if bodyInput.Name == "" {
		slog.Debug("Invalid body input", "error", err)
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	// Modify row in database
	deck := model.NewDeck(bodyInput.Name, bodyInput.Description)
	deck.ID = deckID
	dbErr := s.deck_db.Modify(deck)
	if dbErr != nil {
		if dbErr == utils.ErrMaxLengthExceeded {
			slog.Debug("Max length exceeded", "error", dbErr)
			http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		} else if dbErr == utils.ErrDuplicateKeyViolation {
			slog.Debug("Duplicate key violation", "error", dbErr)
			http.Error(w, DuplicateKeyViolationErrorMessage, http.StatusConflict)
		} else if dbErr == utils.ErrRecordNotExist {
			slog.Debug("Record not exist", "error", dbErr)
			http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		} else {
			slog.Debug("Error modifying deck", "error", dbErr)
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"id": deckID})
	if err != nil {
		slog.Debug("Error encoding deck ID", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Sent response", "deck ID", deckID)
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
		slog.Debug(fmt.Sprintf("Invalid deck ID %s", idStr))
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Fetch from database
	dbErr := s.deck_db.Delete(id)
	if dbErr != nil {
		if dbErr == utils.ErrRecordNotExist {
			slog.Debug("Record not exist", "error", dbErr)
			http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		} else {
			slog.Debug("Error deleting deck", "error", dbErr)
			http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	encodingErr := json.NewEncoder(w).Encode(map[string]int{"id": id})
	if encodingErr != nil {
		slog.Debug("Error encoding deck ID", "error", encodingErr)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Sent response", "deck ID", id)
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
		slog.Debug("Error encoding deck name max length", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Sent response", "deck name max length", database.DeckColumnNameMaxLength)
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
		slog.Debug("Error encoding deck description max length", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Sent response", "deck description max length", database.DeckColumnDescriptionMaxLength)
}

// HandleInsertCard handles the HTTP POST request for inserting a new card into a deck.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 400 Bad Request : If the deck ID is invalid or the request body is invalid.
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the card is found and the request is successful.
func (s *APIServer) HandleInsertCard(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idStr := strings.Split(r.URL.Path, "/")[2]
	deckID, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Debug(fmt.Sprintf("Invalid deck ID %s", idStr))
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Parse content body
	type Content struct {
		Fields []string `json:"fields"`
		Values []string `json:"values"`
	}
	type InsertInput struct {
		Content Content `json:"content"`
		Source  string  `json:"source"`
	}
	var bodyInput InsertInput
	err = json.NewDecoder(r.Body).Decode(&bodyInput)
	if err != nil {
		slog.Debug("Error decoding request body", "error", err)
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	} else if bodyInput.Content.Fields == nil || bodyInput.Content.Values == nil {
		slog.Debug("Missing mandatory field content")
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	} else if len(bodyInput.Content.Fields) != len(bodyInput.Content.Values) {
		slog.Debug("Fields and values length mismatch")
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	} else if len(bodyInput.Content.Fields) == 0 || len(bodyInput.Content.Values) == 0 {
		slog.Debug("Fields or values length is 0")
		http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
		return
	}

	for _, value := range bodyInput.Content.Values {
		if value == "" {
			slog.Debug("Value is empty")
			http.Error(w, InvalidBodyErrorMessage, http.StatusBadRequest)
			return
		}
	}

	contentBytes, err := json.Marshal(bodyInput.Content)
	if err != nil {
		slog.Debug("Error marshalling content", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	card := model.NewCard(deckID, string(contentBytes), bodyInput.Source)

	cardID, dbErr := s.card_db.Insert(card)
	if dbErr != nil {
		slog.Debug(fmt.Sprintf("Error inserting card %s", dbErr))
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"id": cardID})
	if err != nil {
		slog.Debug(fmt.Sprintf("Error encoding card ID %s", err))
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	slog.Debug(fmt.Sprintf("Sent response, card ID: %d", cardID))
	w.WriteHeader(http.StatusCreated)
}

// HandleGetTotalCards handles the HTTP GET request for retrieving the total number of cards in a deck.
//
// Parameters:
//   - w http.ResponseWriter : The response writer to send the response.
//   - r *http.Request : The HTTP request containing the deck ID in the URL path.
//
// Errors:
//   - 400 Bad Request : If the deck ID is invalid.
//   - 500 Internal Server Error : If there is an error while processing the request.
//   - 200 OK : If the total number of cards is found and the request is successful.
func (s *APIServer) HandleGetTotalCards(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idStr := strings.Split(r.URL.Path, "/")[2]
	deckID, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Debug(fmt.Sprintf("Invalid deck ID %s", idStr))
		http.Error(w, InvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	// Fetch from database
	count, dbErr := s.card_db.GetTotalCards(deckID)
	if dbErr != nil {
		slog.Debug("Error getting total cards", "error", dbErr)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"total": count})
	if err != nil {
		slog.Debug("Error encoding total cards", "error", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
