package main

import (
	"flash-learn/internal/api"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db := connectToPostgres()

	server := api.NewAPIServer("localhost:8080", db)
	err := server.Run()

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server started on localhost:8080")

	defer db.Close()
}
