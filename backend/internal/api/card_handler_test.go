package api

import (
	"flash-learn/internal/database"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APICardServerTestSuite struct {
	suite.Suite
	address string
	deck_db *database.DeckDBWrapperMock
	card_db *database.CardDBWrapperMock
	server  *APIServer
}

func (suite *APICardServerTestSuite) SetupTest() {
	suite.address = "localhost:8080"
	suite.deck_db = database.NewDeckDBWrapperMock()
	suite.card_db = database.NewCardDBWrapperMock()
	suite.server = NewAPIServer(suite.address, suite.deck_db, suite.card_db)
}

func (suite *APICardServerTestSuite) TearDownTest() {
	suite.server = nil
}

func (suite *APICardServerTestSuite) TestNewAPIServer() {
	assert.Equal(suite.T(), suite.address, suite.server.address, "Expected address to be set correctly")
	assert.NotNil(suite.T(), suite.server.deck_db, "Expected db to be initialized")
	assert.NotNil(suite.T(), suite.server.card_db, "Expected card db to be initialized")
}

func TestAPICardServerTestSuite(t *testing.T) {
	suite.Run(t, new(APICardServerTestSuite))
}

func (suite *APICardServerTestSuite) TestInsertCardHandlerWithBadRequests() {
	suite.card_db.CreateTable()

	testCases := []struct {
		name           string
		deckID         string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Bad Request (Deck ID not number)",
			deckID:         "a",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidDeckIDErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Deck ID not number)",
			deckID:         "1a",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidDeckIDErrorMessage + "\n",
		},
		{
			name:           "Bad Request (empty request body)",
			deckID:         "1",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Mandatory field missing)",
			deckID:         "1",
			requestBody:    "{}",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Content is null)",
			deckID:         "1",
			requestBody:    `{"content": null}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Content missing fields)",
			deckID:         "1",
			requestBody:    `{"content": "{}"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Fields and values length mismatch)",
			deckID:         "0",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": [""]}}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Fields and values length mismatch)",
			deckID:         "0",
			requestBody:    `{"content": {"fields": [], "values": []}}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Fields have empty strings)",
			deckID:         "0",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": ["", ""]}}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
		},
		{
			name:           "Bad Request (Deck doesn't exist)",
			deckID:         "0",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": ["Test front", "Test back"]}}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidDeckIDErrorMessage + "\n",
		},
	}

	for _, tc := range testCases {
		// Create a new request
		req := httptest.NewRequest(http.MethodGet, "/deck/"+tc.deckID+"/card", nil)
		req.Body = io.NopCloser(strings.NewReader(tc.requestBody))

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		suite.server.HandleInsertCard(rr, req)

		assert.Equal(suite.T(), tc.expectedStatus, rr.Code, "Expected status code to be %d, got %d", tc.expectedStatus, rr.Code)
		assert.Equal(suite.T(), tc.expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", tc.expectedBody, rr.Body.String())
	}
}

func (suite *APICardServerTestSuite) TestInsertCardHandlerWithValid() {
	suite.card_db.CreateTable()
	suite.card_db.InsertDeck(0)
	suite.card_db.InsertDeck(1)
	suite.card_db.InsertDeck(2)

	testCases := []struct {
		name           string
		deckID         string
		cardID         string
		requestBody    string
		expectedStatus int
		expectedBody   string
		expectedCount  int
	}{
		{
			name:           "Bad Request (Deck doesn't exist)",
			deckID:         "0",
			cardID:         "0",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": ["Test front", "Test back"]}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":0}\n",
			expectedCount:  1,
		},
		{
			name:           "Bad Request (Deck doesn't exist)",
			deckID:         "0",
			cardID:         "1",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": ["Test front 2", "Test back 2"]}, "source": "Test source"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1}\n",
			expectedCount:  2,
		},
		{
			name:           "Bad Request (Deck doesn't exist)",
			deckID:         "1",
			cardID:         "0",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": ["Test front 3", "Test back 3"]}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":0}\n",
			expectedCount:  1,
		},
		{
			name:           "Bad Request (Deck doesn't exist)",
			deckID:         "1",
			cardID:         "1",
			requestBody:    `{"content": {"fields": ["front", "back"], "values": ["Test front 4", "Test back 4"]}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1}\n",
			expectedCount:  2,
		},
	}

	for _, tc := range testCases {
		// Create a new request
		req := httptest.NewRequest(http.MethodGet, "/deck/"+tc.deckID+"/card", nil)
		req.Body = io.NopCloser(strings.NewReader(tc.requestBody))

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		suite.server.HandleInsertCard(rr, req)

		assert.Equal(suite.T(), tc.expectedStatus, rr.Code, "Expected status code to be %d, got %d", tc.expectedStatus, rr.Code)
		assert.Equal(suite.T(), tc.expectedBody, rr.Body.String(), "Expected response body to be '%s', got '%s'", tc.expectedBody, rr.Body.String())

		deckID, _ := strconv.Atoi(tc.deckID)
		count, _ := suite.card_db.GetCardCount(deckID)
		assert.Equal(suite.T(), tc.expectedCount, count, "Expected card count to be %d, got %d", tc.expectedCount, count)
	}

	count, _ := suite.card_db.GetCardCount(2)
	assert.Equal(suite.T(), 0, count, "Expected card count to be %d, got %d", 0, count)
}
