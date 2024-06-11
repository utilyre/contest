package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	_ "github.com/lib/pq"
)

var (
	ErrInvalidOption = errors.New("invalid option")
)

type options struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string
}

type Option func(*options) error

func WithHost(host string) Option {
	return func(o *options) error {
		o.host = host
		return nil
	}
}

func WithPort(port string) Option {
	return func(o *options) error {
		o.port = port
		return nil
	}
}

func WithUser(user string) Option {
	return func(o *options) error {
		o.user = user
		return nil
	}
}

func WithPassword(password string) Option {
	return func(o *options) error {
		o.password = password
		return nil
	}
}

func WithDBName(dbname string) Option {
	return func(o *options) error {
		o.dbname = dbname
		return nil
	}
}

func WithSSLMode(sslmode string) Option {
	return func(o *options) error {
		if !slices.Contains([]string{
			"require",
			"verify-full",
			"verify-ca",
			"disable",
		}, sslmode) {
			return fmt.Errorf("sslmode: string '%s': %w", sslmode, ErrInvalidOption)
		}

		o.sslmode = sslmode
		return nil
	}
}

func New(ctx context.Context, opts ...Option) (*sql.DB, error) {
	o := &options{
		host:     "127.0.0.1",
		port:     "5432",
		user:     "postgres",
		password: "",
		dbname:   "postgres",
		sslmode:  "require",
	}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		o.host, o.port,
		o.user, o.password,
		o.dbname, o.sslmode,
	))
	if err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	return db, nil
}
