package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type options struct {
	dsn string
}

type Option func(*options) error

func WithDSN(dsn string) Option {
	return func(o *options) error {
		o.dsn = dsn
		return nil
	}
}

func New(ctx context.Context, opts ...Option) (*sql.DB, error) {
	o := &options{dsn: ""}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("postgres", o.dsn)
	if err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	return db, nil
}
