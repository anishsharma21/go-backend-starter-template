package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET_KEY = make([]byte, 64)

func init() {
	secretKeyString := os.Getenv("JWT_SECRET_KEY")
	if secretKeyString == "" {
		slog.Error("JWT_SECRET_KEY environment variable not set")
		os.Exit(1)
	}

	JWT_SECRET_KEY = []byte(secretKeyString)
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header was empty")
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := VerifyToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			slog.Error("Error while verifying access token", "error", err, "access_token", tokenString)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CreateAccessToken(email string) (string, error) {
	return createToken(email, time.Now().Add(15*time.Minute))
}

func CreateRefreshToken(email string) (string, error) {
	return createToken(email, time.Now().Add(7*24*time.Hour))
}

type CustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func createToken(email string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	})

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v\n", t.Header["alg"])
		}

		return JWT_SECRET_KEY, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims: email not found")
	}

	return email, nil
}
