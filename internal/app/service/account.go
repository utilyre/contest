package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/utilyre/contest/internal/app"
	"github.com/utilyre/contest/internal/app/domain"
	"github.com/utilyre/contest/internal/app/port"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	accountRepo port.AccountRepo
}

func NewAccountService(accountRepo port.AccountRepo) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

type AccountRegisterParams struct {
	Username string
	Email    string
	Password string
}

type AccountRegisterOutput struct {
	Account *domain.Account
}

func (as *AccountService) Register(
	ctx context.Context,
	params *AccountRegisterParams,
) (*AccountRegisterOutput, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt: password hash: %w", err)
	}

	username, err := domain.NewUsername(params.Username)
	if err != nil {
		return nil, fmt.Errorf("domain: %w", err)
	}

	email, err := domain.NewEmail(params.Email)
	if err != nil {
		return nil, fmt.Errorf("domain: %w", err)
	}

	password, err := domain.NewPassword(hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("domain: %w", err)
	}

	account := &domain.Account{
		Username: username,
		Email:    email,
		Password: password,
	}

	id, err := as.accountRepo.CreateAccountReturningID(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	account.ID = id

	return &AccountRegisterOutput{Account: account}, nil
}

type AccountLoginParams struct {
	Username string
	Password string
}

type AccountLoginOutput struct {
	Token string
}

func (as *AccountService) Login(
	ctx context.Context,
	params *AccountLoginParams,
) (*AccountLoginOutput, error) {
	username, err := domain.NewUsername(params.Username)
	if err != nil {
		return nil, fmt.Errorf("domain: %w", err)
	}

	hashedPassword, err := as.accountRepo.GetAccountPasswordByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrAccountNotFound
		}

		return nil, fmt.Errorf("repo: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(params.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, app.ErrAccountNotFound
		}

		return nil, fmt.Errorf("bcrypt: password comparison: %w", err)
	}

	token, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, newAccountClaims(params.Username)).
		SignedString([]byte("secret"))
	if err != nil {
		return nil, fmt.Errorf("jwt: token signing: %w", err)
	}

	return &AccountLoginOutput{
		Token: token,
	}, nil
}

const accessTokenExpiryTime = 5 * time.Minute

type accountClaims struct {
	jwt.RegisteredClaims
	Username string
}

func newAccountClaims(username string) *accountClaims {
	return &accountClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiryTime)),
		},
		Username: username,
	}
}
