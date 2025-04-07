package api

import (
	"context"
	"encoding/json"
	"flash-learn/internal/database"
	"flash-learn/internal/model"
	"fmt"
	"net/http"
	"net/http/httptest"
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
			name:           "Bad Request for valid deck ID (database doesn't exist)",
			deckID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   GetSingleDeckInternalServerErrorMessage + "\n",
		},
		{
			name:           "Bad Request for valid deck ID (deck table doesn't exist)",
			deckID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   GetSingleDeckInternalServerErrorMessage + "\n",
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

		assert.Equal(suite.T(), tc.expectedStatus, rr.Code, "Expected status code to be %d, got %d", tc.expectedStatus, rr.Code)
		assert.Equal(suite.T(), string(jsonData)+"\n", rr.Body.String(), "Expected response body to be '%s', got '%s'", tc.expectedBody, rr.Body.String())
	}
}
