package queries

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anishsharma21/go-backend-starter-template/internal/types/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllUsers(ctx context.Context, dbPool *pgxpool.Pool) ([]models.User, error) {
	query := `SELECT * FROM users`

	rows, err := dbPool.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch users: %v\n", err)
	}

	var users []models.User
	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("Failed to collect users: %v\n", err)
	}

	return users, nil
}

func DeleteAllUsers(ctx context.Context, dbPool *pgxpool.Pool) error {
	query := `DELETE FROM users`

	slog.Info("Delete users transaction beginning...")

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %v\n", err)
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("Failed to delete users: %v\n", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %v\n", err)
	}

	slog.Info("Delete users transaction completed successfully.", "count", ct.RowsAffected())

	return nil
}
