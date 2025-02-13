package queries

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anishsharma21/go-backend-starter-template/internal/types/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserByEmail(ctx context.Context, dbPool *pgxpool.Pool, email string) (models.User, error) {
	query := "SELECT * from users WHERE email = $1"

	rows, err := dbPool.Query(ctx, query, email)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve data for user with email %q: %v", email, err)
	}
	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return models.User{}, fmt.Errorf("failed to collect data from database for user with email %q: %v", email, err)
	}

	return user, nil
}

func SignUpNewUser(ctx context.Context, dbPool *pgxpool.Pool, user models.User) error {
	args := pgx.NamedArgs{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"password":   user.Password,
	}

	query := "INSERT INTO users (email, first_name, last_name, password) VALUES (@email, @first_name, @last_name, @password)"

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	ct, err := tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to execute query: %v", err)
	}

	if ct.RowsAffected() != 1 {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to insert user: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	slog.Info("User signed up successfully", "email", user.Email, "command_tag", ct.String())

	return nil
}
