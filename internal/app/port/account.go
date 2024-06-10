package port

import (
	"context"
	"github.com/utilyre/contest/internal/app/domain"
)

type AccountRepo interface {
	CreateAccountReturningID(
		ctx context.Context,
		account *domain.Account,
	) (int32, error)

	GetAccountByEmail(
		ctx context.Context,
		email domain.Email,
	) (*domain.Account, error)

	GetAccountPasswordByUsername(
		ctx context.Context,
		username domain.Username,
	) (domain.Password, error)
}
