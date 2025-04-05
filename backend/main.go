package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	const driverName string = "postgres"
	const connStr string = "postgres://postgres:flashlearn@localhost:5432/fl-db?sslmode=disable"

	db, err := sql.Open(driverName, connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
