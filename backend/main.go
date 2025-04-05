package main

import (
	_ "github.com/lib/pq"
)

func main() {
	db := connectToPostgres()

	defer db.Close()
}
