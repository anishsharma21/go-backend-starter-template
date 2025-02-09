package handlers

import (
	"context"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/jackc/pgx/v5"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func AddUser(dbConn *pgx.Conn, templates *template.Template) http.Handler {
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

		cmdTag, err := dbConn.Exec(context.Background(), query, args)
		if err != nil {
			slog.Error("Failed to insert user", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		slog.Info("User added successfully", "command tag", cmdTag.String(), "rows affected", cmdTag.RowsAffected())

		err = templates.ExecuteTemplate(w, "button", nil)
		if err != nil {
			slog.Error("Failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
