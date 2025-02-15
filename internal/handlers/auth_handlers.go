package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/anishsharma21/go-backend-starter-template/internal/middleware"
	"github.com/anishsharma21/go-backend-starter-template/internal/queries"
	"github.com/anishsharma21/go-backend-starter-template/internal/types/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := template.HTMLEscapeString(r.FormValue("email"))
		firstName := template.HTMLEscapeString(r.FormValue("first_name"))
		lastName := template.HTMLEscapeString(r.FormValue("last_name"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if email == "" || firstName == "" || lastName == "" || password == "" {
			slog.Error("Email, first name, last name or password is empty")
			http.Error(w, "Email, first name, last name or password is empty", http.StatusBadRequest)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		passwordHashString := string(passwordHash)

		err = queries.SignUpNewUser(r.Context(), dbPool, models.User{
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
			Password:  &passwordHashString,
		})
		if err != nil {
			slog.Error("Failed to sign up new user", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		accessToken, err := middleware.CreateAccessToken(email)
		if err != nil {
			slog.Error("Failed to create JWT token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		refreshToken, err := middleware.CreateRefreshToken(email)
		if err != nil {
			slog.Error("Failed to create refresh token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(24 * 7 * time.Hour),
		})

		response := map[string]string{
			"token": accessToken,
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

func Login(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := template.HTMLEscapeString(r.FormValue("email"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if email == "" || password == "" {
			slog.Error("Email or password is empty")
			http.Error(w, "Email or password is empty", http.StatusBadRequest)
			return
		}

		user, err := queries.GetUserByEmail(r.Context(), dbPool, email)
		if err != nil {
			slog.Error("Failed to find user when logging in", "error", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Email != email {
			slog.Error("User email does not match", "user_email", user.Email, "email", email)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
		if err != nil {
			slog.Error("Failed to compare password hashes", "error", err)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		accessToken, err := middleware.CreateAccessToken(user.Email)
		if err != nil {
			slog.Error("Failed to create JWT token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		refreshToken, err := middleware.CreateRefreshToken(email)
		if err != nil {
			slog.Error("Failed to create refresh token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(24 * 7 * time.Hour),
		})

		response := map[string]string{
			"token": accessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		slog.Info(fmt.Sprintf("User logged in: %s", email))
	})
}

func RefreshToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			slog.Error("Failed to get refresh token cookie from request", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value
		email, err := middleware.VerifyToken(refreshToken)
		if err != nil {
			slog.Error("Error validating refresh token", "error", err)
			http.Error(w, "Session ended. Login again.", http.StatusUnauthorized)
			return
		}

		accessToken, err := middleware.CreateAccessToken(email)
		if err != nil {
			slog.Error("Failed to create new access token", "error", err, "email", email)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"token": accessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Failed to encode refresh token response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
