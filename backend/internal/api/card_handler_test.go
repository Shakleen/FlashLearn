package api

import (
	"flash-learn/internal/database"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APICardServerTestSuite struct {
	suite.Suite
	address string
	db      *database.DeckDBWrapperMock
	server  *APIServer
}

func (suite *APICardServerTestSuite) SetupTest() {
	suite.address = "localhost:8080"
	suite.db = database.NewDeckDBWrapperMock()
	suite.server = NewAPIServer(suite.address, suite.db, nil)
}

func (suite *APICardServerTestSuite) TearDownTest() {
	suite.server = nil
}

func (suite *APICardServerTestSuite) TestNewAPIServer() {
	assert.Equal(suite.T(), suite.address, suite.server.address, "Expected address to be set correctly")
	assert.NotNil(suite.T(), suite.server.deck_db, "Expected db to be initialized")
}

func TestAPICardServerTestSuite(t *testing.T) {
	suite.Run(t, new(APICardServerTestSuite))
}

func (suite *APICardServerTestSuite) TestInsertCardHandlerWithBadRequests() {
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
			name:           "Bad Request (Front and back are empty)",
			deckID:         "1",
			requestBody:    `{"content": {"front": "", "back": ""}}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   InvalidBodyErrorMessage + "\n",
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
