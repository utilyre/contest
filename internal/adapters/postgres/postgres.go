package postgres

import "database/sql"

func New() (*sql.DB, error) {
	return sql.Open("postgres", "postgresql://admin:admin@localhost:5432/contest?sslmode=disable")
}
