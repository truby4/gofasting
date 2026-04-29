package auth

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store store
}

func New(db *sql.DB) *Service {
	return &Service{
		store: store{
			db: db,
		},
	}
}

func (s *Service) Authenticate(email, password string) (int, error) {
	input := input{
		Email:    email,
		Password: password,
	}

	err := input.validateInput()
	if err != nil {
		return 0, err
	}

	u, err := s.store.getByEmail(email)
	if err != nil {
		if errors.Is(err, ErrNoRecordFound) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return u.ID, nil
}

func (s *Service) Register(email, password string) error {
	input := input{
		Email:    email,
		Password: password,
	}

	err := input.validateInput()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	return s.store.create(email, hashedPassword)
}

func (s *Service) Exists(userID int) (bool, error) {
	return s.store.exists(userID)
}
