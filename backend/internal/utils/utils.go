package utils

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
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

func getPostgresConfig() (Config, error) {
	slog.Debug("Getting postgres config")
	const envFileName string = "postgres.env"
	err := godotenv.Load(envFileName)

	if err != nil {
		slog.Error("Error loading secret.env file", "error", err)
		return Config{}, err
	}

	config := Config{
		DriverName: os.Getenv("POSTGRES_DRIVER"),
		UserName:   os.Getenv("POSTGRES_USER"),
		Password:   os.Getenv("POSTGRES_PASSWORD"),
		Host:       os.Getenv("POSTGRES_HOST"),
		Port:       os.Getenv("POSTGRES_PORT"),
		DbName:     os.Getenv("POSTGRES_DB"),
		SSLMode:    os.Getenv("POSTGRES_SSL_MODE"),
	}

	return config, nil
}

func ConnectToPostgres() (*sql.DB, error) {
	config, err := getPostgresConfig()

	if err != nil {
		slog.Error("Error getting postgres config", "error", err)
		return nil, err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.UserName, config.Password, config.Host, config.Port, config.DbName, config.SSLMode)

	db, err := sql.Open(config.DriverName, connStr)

	if err != nil {
		slog.Error("Error connecting to postgres", "error", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		slog.Error("Error pinging postgres", "error", err)
		return nil, err
	}

	slog.Info("Connected to postgres")
	return db, nil
}
