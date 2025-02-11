package queries

import (
	"context"
	"fmt"

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
