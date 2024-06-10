package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	rgAlphanumeric = regexp.MustCompile(`^[\w]*$`)
	rgEmail        = regexp.MustCompile(`^[\w-+\.]+@([\w-]+\.)+[\w-]{2,4}$`)
)

var (
	ErrTooShort      = errors.New("too short")
	ErrTooLong       = errors.New("too long")
	ErrIllegalFormat = errors.New("illegal format")
)

type Account struct {
	ID        int32
	CreatedAt time.Time

	Username Username
	Email    Email
	Password Password
}

func NewAccount(
	username Username,
	email Email,
	password Password,
) *Account {
	return &Account{
		Username: username,
		Email:    email,
		Password: password,
	}
}

type Username string

func NewUsername(v string) (Username, error) {
	if len(v) < 3 {
		return "", fmt.Errorf("username: string '%s': %w", v, ErrTooShort)
	}
	if len(v) > 64 {
		return "", fmt.Errorf("username: string '%s': %w", v, ErrTooLong)
	}
	if !rgAlphanumeric.MatchString(v) {
		return "", fmt.Errorf("username: string '%s': %w", v, ErrIllegalFormat)
	}

	return Username(v), nil
}

func (u Username) String() string {
	return string(u)
}

type Email string

func NewEmail(v string) (Email, error) {
	if len(v) > 320 {
		return "", fmt.Errorf("email: string '%s': %w", v, ErrTooLong)
	}
	if !rgEmail.MatchString(v) {
		return "", fmt.Errorf("email: string '%s': %w", v, ErrIllegalFormat)
	}

	return Email(v), nil
}

func (e Email) String() string {
	return string(e)
}

type Password []byte

func NewPassword(v []byte) (Password, error) {
	if len(v) < 8 {
		return nil, fmt.Errorf("password: %w", ErrTooShort)
	}
	if len(v) > 1024 {
		return nil, fmt.Errorf("password: %w", ErrTooLong)
	}

	return Password(v), nil
}

func (p Password) String() string {
	return "xxx"
}
