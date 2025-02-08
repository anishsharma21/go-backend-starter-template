package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anishsharma21/go-backend-starter-template/internal/types"
	"github.com/jackc/pgx/v5"
)

var (
	env types.Environment
)

func init() {
	env = types.StringToEnv(os.Getenv("ENV"))
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := setupServer(port)

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
	slog.Warn("Received signal", "signal", sig.String)

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error", "error", err)
	}
	<-shutdownChan
	close(shutdownChan)

	slog.Info("Graceful server shutdown complete.")
}

func setupDB() (*pgx.Conn, error) {
	var connStr string
	if env.IsProduction() || env.IsCICD() {
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

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}

func setupServer(port string) *http.Server {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: setupRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			slog.Info("Server started on port 8080...")
			return context.Background()
		},
	}

	return server
}
