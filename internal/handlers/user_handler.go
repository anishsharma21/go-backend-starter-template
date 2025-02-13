package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-backend-starter-template/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUsers(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := queries.GetAllUsers(r.Context(), dbPool)
		if err != nil {
			slog.Error("Failed to fetch users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(users); err != nil {
			slog.Error("Failed to encode users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		slog.Info("Users JSON data fetched successfully")
	})
}

func DeleteUsers(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := queries.DeleteAllUsers(r.Context(), dbPool)
		if err != nil {
			slog.Error("Failed to delete users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	})
}
