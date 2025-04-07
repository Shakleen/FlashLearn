package api

import (
	"flash-learn/internal/database"
	"flash-learn/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

const (
	GetSingleDeckInvalidDeckIDErrorMessage  string = "Invalid deck ID"
	GetSingleDeckNotFoundErrorMessage       string = "Deck not found"
	GetSingleDeckInternalServerErrorMessage string = "Internal server error"
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
	router.HandleFunc("/deck/{id}", s.HandleGetSingleDeck)
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
	idStr := strings.Split(r.URL.Path, "/")[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, GetSingleDeckInvalidDeckIDErrorMessage, http.StatusBadRequest)
		return
	}

	_, dbErr := s.db.GetSingle(id)
	if dbErr != nil {
		if dbErr == utils.ErrRecordNotExist {
			http.Error(w, GetSingleDeckNotFoundErrorMessage, http.StatusBadRequest)
		} else {
			http.Error(w, GetSingleDeckInternalServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
