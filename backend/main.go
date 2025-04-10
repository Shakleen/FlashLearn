package main

import (
	"flash-learn/internal/api"
	"flash-learn/internal/database"
	"flash-learn/internal/utils"
	"log/slog"

	_ "github.com/lib/pq"
)

func main() {
	logger := utils.GetLogger()
	slog.SetDefault(logger)

	slog.Info("Starting server")
	db, err := utils.ConnectToPostgres()
	if err != nil {
		slog.Error("Error connecting to postgres", "error", err)
		return
	}
	db_wrapper := database.NewDeckDBWrapper(db)
	card_db_wrapper := database.NewCardDBWrapper(db)

	slog.Info("Creating table if not exists")
	db_wrapper.CreateTable()
	card_db_wrapper.CreateTable()

	slog.Info("Starting API server")
	server := api.NewAPIServer("localhost:8080", db_wrapper, card_db_wrapper)
	err = server.Start()

	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}

	slog.Info("Server started on localhost:8080")
	defer db.Close()
}
