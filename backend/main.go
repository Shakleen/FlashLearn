package main

import (
	"flash-learn/internal/api"
	"flash-learn/internal/database"
	"flash-learn/internal/utils"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db := utils.ConnectToPostgres()

	server := api.NewAPIServer("localhost:8080", database.NewDeckDBWrapper(db))
	err := server.Start()

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server started on localhost:8080")

	defer db.Close()
}
