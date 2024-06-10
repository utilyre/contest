package postgres

import (
	"context"
	"fmt"

	"github.com/utilyre/contest/gen/sqlc"
	"github.com/utilyre/contest/internal/app/domain"
)

type AccountRepo struct {
	queries *sqlc.Queries
}

func NewAccountRepo(db sqlc.DBTX) *AccountRepo {
	return &AccountRepo{queries: sqlc.New(db)}
}

func (ur *AccountRepo) CreateAccountReturningID(
	ctx context.Context,
	account *domain.Account,
) (int32, error) {
	id, err := ur.queries.CreateAccountReturningID(ctx, sqlc.CreateAccountReturningIDParams{
		Username: string(account.Username),
		Email:    string(account.Email),
		Password: []byte(account.Password),
	})
	if err != nil {
		return 0, fmt.Errorf("account: %w", err)
	}

	return id, nil
}

func (ur *AccountRepo) GetAccountByEmail(
	ctx context.Context,
	email domain.Email,
) (*domain.Account, error) {
	acc, err := ur.queries.GetAccountByEmail(ctx, string(email))
	if err != nil {
		return nil, fmt.Errorf("account: %w", err)
	}

	return &domain.Account{
		ID:        acc.ID,
		CreatedAt: acc.CreatedAt,
		Username:  domain.Username(acc.Username),
		Email:     domain.Email(acc.Email),
		Password:  domain.Password(acc.Password),
	}, nil
}

func (ur *AccountRepo) GetAccountPasswordByUsername(
	ctx context.Context,
	username domain.Username,
) (domain.Password, error) {
	hash, err := ur.queries.GetAccountPasswordByUsername(ctx, string(username))
	if err != nil {
		return nil, fmt.Errorf("account: %w", err)
	}

	return hash, nil
}
