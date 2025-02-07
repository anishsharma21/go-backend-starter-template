package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func main() {
	_, err := setupDB()
	if err != nil {
		slog.Error("Failed to initialise database connection", "error", err)
		return
	}

	slog.Info("Connected to database successfully.")
}

func setupDB() (*pgx.Conn, error) {
	env := os.Getenv("ENVIRONMENT")
	var connStr string
	if env == "production" || env == "cicd" {
		connStr = os.Getenv("DATABASE_URL")
	} else {
		connStr = "host=localhost port=5432 user=admin password=secret dbname=mydb sslmode=disable"
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v\n", err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v\n", err)
	}

	return conn, nil
}
