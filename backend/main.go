package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	db := connectToPostgres()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /deck/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid deck ID", http.StatusBadRequest)
			return
		}

		deck, err := getDeckDetails(db, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(deck)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe("localhost:8080", mux)

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server started on localhost:8080")

	defer db.Close()
}
