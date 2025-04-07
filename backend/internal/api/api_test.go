package api

import (
	"context"
	"encoding/json"
	"flash-learn/internal/database"
	"flash-learn/internal/model"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIServerTestSuite struct {
	suite.Suite
	address string
	db      *database.DeckDBWrapperMock
	server  *APIServer
}

func (suite *APIServerTestSuite) SetupTest() {
	suite.address = "localhost:8080"
	suite.db = database.NewDeckDBWrapperMock()
	suite.server = NewAPIServer(suite.address, suite.db)
}

func (suite *APIServerTestSuite) TearDownTest() {
	suite.server = nil
}

func (suite *APIServerTestSuite) TestNewAPIServer() {
	assert.Equal(suite.T(), suite.address, suite.server.address, "Expected address to be set correctly")
	assert.NotNil(suite.T(), suite.server.db, "Expected db to be initialized")
}

func TestAPIServerTestSuite(t *testing.T) {
	suite.Run(t, new(APIServerTestSuite))
}

func (suite *APIServerTestSuite) TestGetSingleDeckHandlerWithError() {
	testCases := []struct {
		name           string
		deckID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Bad Request for no deck ID",
			deckID:         "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   GetSingleDeckInvalidDeckIDErrorMessage + "\n",
		},
		{
			name:           "Bad Request for invalid deck ID (non-numeric)",
			deckID:         "a",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   GetSingleDeckInvalidDeckIDErrorMessage + "\n",
		},
		{
			name:           "Bad Request for invalid deck ID (digit + non-numeric)",
			deckID:         "1a",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   GetSingleDeckInvalidDeckIDErrorMessage + "\n",
		},
		{
			name:           "Internal server error for valid deck ID (database doesn't exist)",
			deckID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   InternalServerErrorMessage + "\n",
		},
		{
			name:           "Internal server error for valid deck ID (deck table doesn't exist)",
			deckID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   InternalServerErrorMessage + "\n",
		},
	}

	for _, tc := range testCases {
		// Create a new request
		req := httptest.NewRequest(http.MethodGet, "/deck/"+tc.deckID, nil)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		ctx := req.Context()
		req = req.WithContext(context.WithValue(ctx, "id", tc.deckID))

		suite.server.HandleGetSingleDeck(rr, req)

		assert.Equal(suite.T(), tc.expectedStatus, rr.Code, "Expected status code to be %d, got %d", tc.expectedStatus, rr.Code)
		assert.Equal(suite.T(), tc.expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", tc.expectedBody, rr.Body.String())
	}
}

func (suite *APIServerTestSuite) TestGetSingleDeckHandlerWithValid() {
	suite.db.CreateTable()

	decks := []model.Deck{}
	decks = append(decks, model.NewDeck("Deck #1", "This is a first deck"))
	decks = append(decks, model.NewDeck("Deck #2", "This is a second deck"))

	suite.db.Insert(decks[0])
	suite.db.Insert(decks[1])

	testCases := []struct {
		name           string
		deckID         string
		expectedStatus int
		expectedBody   model.Deck
	}{
		{
			name:           "Valid deck id",
			deckID:         "0",
			expectedStatus: http.StatusOK,
			expectedBody:   decks[0],
		},
		{
			name:           "Valid deck id",
			deckID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   decks[1],
		},
	}

	for _, tc := range testCases {
		// Create a new request
		req := httptest.NewRequest(http.MethodGet, "/deck/"+tc.deckID, nil)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		ctx := req.Context()
		req = req.WithContext(context.WithValue(ctx, "id", tc.deckID))

		suite.server.HandleGetSingleDeck(rr, req)

		jsonData, err := json.Marshal(tc.expectedBody)
		if err != nil {
			fmt.Println("Error marshaling to JSON:", err)
			return
		}
		expectedBody := string(jsonData) + "\n"

		assert.Equal(suite.T(), tc.expectedStatus, rr.Code, "Expected status code to be %d, got %d", tc.expectedStatus, rr.Code)
		assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
	}
}

func (suite *APIServerTestSuite) TestGetAllDecksHandlerWithError() {
	expectedStatus := http.StatusInternalServerError
	expectedBody := InternalServerErrorMessage + "\n"

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleGetAllDecks(rr, req)

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}

func (suite *APIServerTestSuite) TestGetAllDecksHandlerEmpty() {
	suite.db.CreateTable()

	expectedStatus := http.StatusOK

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleGetAllDecks(rr, req)
	expectedBody := "[]\n"

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}

func (suite *APIServerTestSuite) TestGetAllDecksHandlerNonEmpty() {
	suite.db.CreateTable()

	decks := []model.Deck{}
	decks = append(decks, model.NewDeck("Deck #1", "This is a first deck"))
	decks = append(decks, model.NewDeck("Deck #2", "This is a second deck"))

	suite.db.Insert(decks[0])
	suite.db.Insert(decks[1])

	expectedStatus := http.StatusOK

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleGetAllDecks(rr, req)

	jsonData, err := json.Marshal(decks)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	expectedBody := string(jsonData) + "\n"

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}

func (suite *APIServerTestSuite) TestGetDeckCountHandlerWithError() {
	expectedStatus := http.StatusInternalServerError
	expectedBody := InternalServerErrorMessage + "\n"

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleGetDeckCount(rr, req)

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}

func (suite *APIServerTestSuite) TestGetDeckCountHandlerEmpty() {
	suite.db.CreateTable()

	expectedStatus := http.StatusOK

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleGetDeckCount(rr, req)

	jsonData, err := json.Marshal(map[string]int{"count": 0})
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	expectedBody := string(jsonData) + "\n"

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}

func (suite *APIServerTestSuite) TestGetDeckCountHandlerNonEmpty() {
	suite.db.CreateTable()

	decks := []model.Deck{}
	decks = append(decks, model.NewDeck("Deck #1", "This is a first deck"))
	decks = append(decks, model.NewDeck("Deck #2", "This is a second deck"))

	suite.db.Insert(decks[0])
	suite.db.Insert(decks[1])

	expectedStatus := http.StatusOK

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleGetDeckCount(rr, req)

	jsonData, err := json.Marshal(map[string]int{"count": len(decks)})
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	expectedBody := string(jsonData) + "\n"

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}

func (suite *APIServerTestSuite) TestInsertDeckHandlerWithError() {
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Bad Request (Empty request body)",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InsertDeckInvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (corrupted request body)",
			requestBody:    `{"name": "Interview Preparation", "description": "A deck containing cards for interview Preparation"`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InsertDeckInvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (No name in request body)",
			requestBody:    `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InsertDeckInvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Name is empty space)",
			requestBody:    `{"name": " "}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InsertDeckInvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Name is longer than max length)",
			requestBody:    `{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InsertDeckInvalidBodyErrorMessage + "\n",
		},
		{
			name: "Bad Request (Description is longer than max length)",
			requestBody: `{"name": "a",
			"description": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InsertDeckInvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Database doesn't exist)",
			requestBody:    `{"name": "Test Name", "description": "Test Description"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   InternalServerErrorMessage + "\n",
		},
	}

	for _, tc := range testCases {
		// Create a new request
		req := httptest.NewRequest(http.MethodGet, "/deck", nil)
		req.Body = io.NopCloser(strings.NewReader(tc.requestBody))

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		suite.server.HandleInsertDeck(rr, req)

		assert.Equal(suite.T(), tc.expectedStatus, rr.Code, "Expected status code to be %d, got %d", tc.expectedStatus, rr.Code)
		assert.Equal(suite.T(), tc.expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", tc.expectedBody, rr.Body.String())
	}
}

func (suite *APIServerTestSuite) TestInsertDeckHandlerWithValid() {
	suite.db.CreateTable()
	expectedStatus := http.StatusOK
	expectedBody := `{"id":0}` + "\n"

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/deck", nil)
	req.Body = io.NopCloser(strings.NewReader(`{"name": "Test Name", "description": "Test Description"}`))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	suite.server.HandleInsertDeck(rr, req)

	assert.Equal(suite.T(), expectedStatus, rr.Code, "Expected status code to be %d, got %d", expectedStatus, rr.Code)
	assert.Equal(suite.T(), expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", expectedBody, rr.Body.String())
}
