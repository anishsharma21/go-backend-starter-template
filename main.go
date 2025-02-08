package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var env string
var dbConnStr string

func init() {
	env = os.Getenv("ENV")
	if env == "production" || env == "cicd" {
		dbConnStr = os.Getenv("DATABASE_URL")
	} else {
		dbConnStr = "postgresql://admin:secret@localhost:5432/mydb?sslmode=disable"
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func main() {
	conn, err := setupDB()
	if err != nil {
		slog.Error("Failed to initialise database connection", "error", err)
		return
	}
	defer conn.Close(context.Background())
	slog.Info("Connected to database successfully.")

	if os.Getenv("RUN_MIGRATION") == "true" {
		slog.Info("Attempting to run database migrations...")
		err := runMigrations()
		if err != nil {
			slog.Error("Failed to run database migrations", "error", err)
			return
		}
		slog.Info("Database migrations complete.")
	} else {
		slog.Info("Database migrations skipped.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: setupRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			url := "http://" + l.Addr().String()
			slog.Info(fmt.Sprintf("Server started on %s", url))
			return context.Background()
		},
	}

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server closed early", "error", err)
		}
		slog.Info("Stopped server new connections.")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Warn("Received signal", "signal", sig.String())

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error occurred", "error", err)
	}
	<-shutdownChan
	close(shutdownChan)

	slog.Info("Graceful server shutdown complete.")
}

func setupDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v\n", err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v\n", err)
	}

	return conn, nil
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			slog.Error("Failed to execute template", "error", err)
		}
	})

	return mux
}

func runMigrations() error {
	if gooseDriver := os.Getenv("GOOSE_DRIVER"); gooseDriver == "" {
		return fmt.Errorf("Goose driver not set: GOOSE_DRIVER=?")
	}

	if gooseDbString := os.Getenv("GOOSE_DBSTRING"); gooseDbString == "" {
		return fmt.Errorf("Goose db string not set: GOOSE_DBSTRING=?")
	}

	if gooseMigrationDir := os.Getenv("GOOSE_MIGRATION_DIR"); gooseMigrationDir == "" {
		return fmt.Errorf("Goose migration dir not set: GOOSE_MIGRATION_DIR=?")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return fmt.Errorf("Failed to open database connection for *sql.DB: %v\n", err)
	}

	if err = goose.Status(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to retrieve status of migrations: %v\n", err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to run `goose up` command: %v\n", err)
	}

	return nil
}
