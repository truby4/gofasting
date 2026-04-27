package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrNoRecordFound      = errors.New("no matching record found")
	ErrValidation         = errors.New("validation error")
)

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return ErrValidation.Error()
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}
