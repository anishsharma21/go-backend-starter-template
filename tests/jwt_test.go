package tests

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/anishsharma21/go-backend-starter-template/internal/handlers"
	"github.com/anishsharma21/go-backend-starter-template/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// dbPool is the database connection pool used for the tests that require database interaction
var dbPool *pgxpool.Pool
var ctx, cancel = context.WithCancel(context.Background())

func TestMain(m *testing.M) {
	// Set up the database connection pool
	dbConnStr := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to parse database connection string.\n")
	}

	for i := 1; i <= 5; i++ {
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil && dbPool != nil {
			break
		}
		log.Printf("Failed to initialise database connection pool")
		log.Printf(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if dbPool == nil {
		log.Fatalf("Failed to initialise database connection pool after 5 attempts")
	}
	defer dbPool.Close()

	// Run the tests
	code := m.Run()

	cancel()
	// Exit after running the tests
	os.Exit(code)
}

// Integration test for the user sign up flow
func TestUserSignUpFlow(t *testing.T) {
	// Prepare
	ts := httptest.NewServer(handlers.SignUp(dbPool))
	defer ts.Close()

	email := "person@gmail.com"
	firstName := "per"
	lastName := "son"
	password := "password"

	// Execute
	resp, err := ts.Client().PostForm(ts.URL, map[string][]string{
		"email":      {email},
		"first_name": {firstName},
		"last_name":  {lastName},
		"password":   {password},
	})
	if err != nil {
		t.Fatalf("Expected no error when sending POST request to /signup, got %v\n", err)
	}
	defer resp.Body.Close()

	// Verify
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v\n", resp.StatusCode)
	}

	// Tear down
	_, err = dbPool.Exec(ctx, "DELETE FROM users WHERE email = $1", email)
	if err != nil {
		t.Fatalf("Failed to delete user from database, %v\n", err)
	}
}

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
