package api

import "net/http"

func addRoutes(router *http.ServeMux, s *APIServer) {
	addDeckRoutes(router, s)
	addCardRoutes(router, s)
}

// addDeckRoutes adds the routes for the deck API.
//
// Parameters:
//   - router *http.ServeMux
//   - s *APIServer
func addDeckRoutes(router *http.ServeMux, s *APIServer) {
	router.HandleFunc("GET /deck/{id}", s.HandleGetSingleDeck)
	router.HandleFunc("GET /deck", s.HandleGetAllDecks)
	router.HandleFunc("GET /deck/count", s.HandleGetDeckCount)
	router.HandleFunc("GET /deck/nameMaxLength", s.HandleGetDeckNameMaxLength)
	router.HandleFunc("GET /deck/descriptionMaxLength", s.HandleGetDeckDescriptionMaxLength)
	router.HandleFunc("POST /deck", s.HandleInsertDeck)
	router.HandleFunc("POST /deck/{id}", s.HandleModifyDeck)
	router.HandleFunc("DELETE /deck/{id}", s.HandleDeleteDeck)
}

// addCardRoutes adds the routes for the card API.
//
// Parameters:
//   - router *http.ServeMux
//   - s *APIServer
func addCardRoutes(router *http.ServeMux, s *APIServer) {
	router.HandleFunc("POST /deck/{id}/card", s.HandleInsertCard)
}
