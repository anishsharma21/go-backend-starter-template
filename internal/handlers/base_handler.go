package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-backend-starter-template/internal/middleware"
	"github.com/anishsharma21/go-backend-starter-template/internal/types/models"
	"github.com/anishsharma21/go-backend-starter-template/internal/types/selectors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func RenderBaseView(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, selectors.IndexPage.BaseHTML, nil)
		if err != nil {
			slog.Error("Failed to execute template", "error", err, "template", selectors.IndexPage.BaseHTML)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}

func HandleLoginRequest(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			slog.Error("Email or password is empty")
			http.Error(w, "Email or password is empty", http.StatusBadRequest)
			return
		}

		var user models.User

		args := pgx.NamedArgs{
			"email": email,
		}

		query := "SELECT * from users WHERE email = @email"

		rows, err := dbPool.Query(r.Context(), query, args)
		if err != nil {
			slog.Error("Failed to retrieve user from database for login", "error", err, "user_email", email)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}
		defer rows.Close()

		user, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
		if err != nil {
			slog.Error("Failed to collect user from database", "error", err, "user_email", email)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if user.Email != email {
			slog.Error("User email does not match", "user_email", user.Email, "email", email)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			slog.Error("Failed to compare password hashes", "error", err)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		jwtToken, err := middleware.CreateToken(user.Email)
		if err != nil {
			slog.Error("Failed to create JWT token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"token": jwtToken,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
