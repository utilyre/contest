package app

import "errors"

var (
	ErrAccountDup      = errors.New("account already exists")
	ErrAccountNotFound = errors.New("account not found")
)
