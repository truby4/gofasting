package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrValidation         = errors.New("failed input validation")
)
