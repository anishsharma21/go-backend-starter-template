package main

import (
	"context"
	"database/sql"
	"embed"
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

	"github.com/anishsharma21/go-backend-starter-template/internal/handlers"
	"github.com/anishsharma21/go-backend-starter-template/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	env       string
	dbConnStr string

	dbPool    *pgxpool.Pool
	templates *template.Template
)

func init() {
	env = os.Getenv("ENV")
	if env == "production" || env == "cicd" {
		dbConnStr = os.Getenv("DATABASE_URL")
	} else {
		dbConnStr = "postgresql://admin:secret@localhost:5432/mydb?sslmode=disable"
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	var err error
	templates, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		os.Exit(1)
	}
	slog.Info("Templates parsed successfully")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := setupDBPool(ctx)
	if err != nil {
		slog.Error("Failed to initialise database connection pool", "error", err)
		return
	}
	defer dbPool.Close()

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
		Handler: setupRoutes(dbPool),
		BaseContext: func(l net.Listener) context.Context {
			url := "http://" + l.Addr().String()
			slog.Info(fmt.Sprintf("Server started on %s", url))
			return ctx
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

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error occurred", "error", err)
	}
	<-shutdownChan
	close(shutdownChan)

	slog.Info("Graceful server shutdown complete.")
}

func setupDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse database connection string: %v", err)
	}

	// Sets the maximum time an idle connection can remain in the pool before being closed
	config.MaxConnIdleTime = 1 * time.Minute
	// To prevent database and backend from ever sleeping, uncomment the following line
	config.MinConns = 1

	var dbPool *pgxpool.Pool
	for i := 1; i <= 5; i++ {
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil && dbPool != nil {
			break
		}
		slog.Warn("Failed to initialise database connection pool", "error", err)
		slog.Info(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if dbPool == nil {
		return nil, fmt.Errorf("Failed to initialise database connection pool after 5 attempts")
	}

	for i := 1; i <= 5; i++ {
		err = dbPool.Ping(ctx)
		if err == nil && dbPool != nil {
			break
		}
		slog.Warn("Failed to ping database connection pool", "error", err)
		slog.Info(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database connection pool after 5 attempts")
	}

	return dbPool, nil
}

func setupRoutes(dbPool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	// Default subpath for endpoints return JSON
	// JSON subpath for endpoints returns JSON
	// JSON should be stable and not change much as it represents data
	// Consumers of these endpoints should be concerned with the JSON structure
	mux.Handle("GET /users", middleware.JWTAuthMiddleware(handlers.GetUsers(dbPool)))
	mux.Handle("POST /signup", handlers.HandleSignUpRequest(dbPool))
	mux.Handle("POST /login", handlers.HandleLoginRequest(dbPool))

	// HTML can be dynamic and change a lot as it represents server state
	// Consumers of these endpoints should not be concerned with the HTML structure
	// example: mux.Handle("GET /users/view", handlers.GetUsersView(dbPool, templates))

	mux.Handle("GET /", handlers.RenderBaseView(templates))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
	defer db.Close()

	if err = goose.Status(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to retrieve status of migrations: %v\n", err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to run `goose up` command: %v\n", err)
	}

	return nil
}
