package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DriverName string
	UserName   string
	Password   string
	Host       string
	Port       string
	DbName     string
	SSLMode    string
}

func main() {
	config := getPostgresConfig()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.UserName, config.Password, config.Host, config.Port, config.DbName, config.SSLMode)

	db, err := sql.Open(config.DriverName, connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

func getPostgresConfig() Config {
	const envFileName string = "postgres.env"
	err := godotenv.Load(envFileName)

	if err != nil {
		log.Fatalf("Error loading secret.env file: %v", err)
	}

	// Get environment variables
	config := Config{
		DriverName: os.Getenv("POSTGRES_DRIVER"),
		UserName:   os.Getenv("POSTGRES_USER"),
		Password:   os.Getenv("POSTGRES_PASSWORD"),
		Host:       os.Getenv("POSTGRES_HOST"),
		Port:       os.Getenv("POSTGRES_PORT"),
		DbName:     os.Getenv("POSTGRES_DB"),
		SSLMode:    os.Getenv("POSTGRES_SSL_MODE"),
	}
	return config
}
