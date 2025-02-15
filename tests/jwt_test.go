package tests

import (
	"testing"
	"time"

	"github.com/anishsharma21/go-backend-starter-template/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAccessToken(t *testing.T) {
	email := "test@example.com"

	tokenString, err := middleware.CreateAccessToken(email)
	if err != nil {
		t.Fatalf("Expected no error, got %v\n", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &middleware.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return middleware.JWT_SECRET_KEY, nil
	})
	if err != nil {
		t.Fatalf("Expected no error while parsing token, got %v\n", err)
	}

	if !token.Valid {
		t.Fatalf("Expected token to be valid")
	}

	claims, ok := token.Claims.(*middleware.CustomClaims)
	if !ok {
		t.Fatalf("Expected claims to be of type CustomClaims")
	}

	if claims.Email != email {
		t.Errorf("Expected email %v, got %v\n", email, claims.Email)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Errorf("Expected token to be valid, but it is expired")
	}
}
