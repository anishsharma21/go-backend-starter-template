package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/anishsharma21/go-backend-starter-template/internal/queries"
	"github.com/anishsharma21/go-backend-starter-template/internal/types/models"
	"github.com/anishsharma21/go-backend-starter-template/internal/types/selectors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type indexAddButtonUserModel struct {
	Users   []models.User
	OOBSwap bool
}

func AddUser(dbPool *pgxpool.Pool, templates *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := randString(5)
		userEmail := userName + "@gmail.com"
		userPassword := randString(10)

		args := pgx.NamedArgs{
			"name":     userName,
			"email":    userEmail,
			"password": userPassword,
		}

		query := `INSERT INTO users (name, email, password) VALUES (@name, @email, @password)`

		cmdTag, err := dbPool.Exec(r.Context(), query, args)
		if err != nil {
			slog.Error("Failed to insert user", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		slog.Info("User added successfully", "command tag", cmdTag.String(), "rows affected", cmdTag.RowsAffected())

		users, err := queries.GetAllUsers(r.Context(), dbPool)
		if err != nil {
			slog.Error("Failed to fetch users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var templateData indexAddButtonUserModel
		templateData.Users = users
		templateData.OOBSwap = true

		err = templates.ExecuteTemplate(w, selectors.IndexPage.AddUserButton, templateData)
		if err != nil {
			slog.Error("Failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

func GetUsers(dbPool *pgxpool.Pool, templates *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := queries.GetAllUsers(r.Context(), dbPool)
		if err != nil {
			slog.Error("Failed to fetch users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// This is where the go backend starter template allows you to choose the response type
		// The default response type is html, but you can choose to get the response in json format
		// For instance, if you are building a mobile app, you can choose json as the response type
		responseType := r.URL.Query().Get("response_type")

		switch responseType {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			if err = json.NewEncoder(w).Encode(users); err != nil {
				slog.Error("Failed to encode users", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		default:
			w.Header().Set("Content-Type", "text/html")
			err = templates.ExecuteTemplate(w, selectors.IndexPage.UsersList, users)
			if err != nil {
				slog.Error("Failed to execute template", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		slog.Info("Users fetched successfully")
	})
}

func DeleteAllUsers(dbPool *pgxpool.Pool, templates *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := `DELETE FROM users`

		cmdTag, err := dbPool.Exec(r.Context(), query)
		if err != nil {
			slog.Error("Failed to delete users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		slog.Info("All users deleted successfully", "command tag", cmdTag.String(), "rows affected", cmdTag.RowsAffected())

		users, err := queries.GetAllUsers(r.Context(), dbPool)
		if err != nil {
			slog.Error("Failed to fetch users", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var templateData indexAddButtonUserModel
		templateData.Users = users
		templateData.OOBSwap = true

		err = templates.ExecuteTemplate(w, selectors.IndexPage.AddUserButton, templateData)
		if err != nil {
			slog.Error("Failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
