package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

func New(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgresql://admin:admin@localhost:5432/contest?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	return db, nil
}
