package store

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/glebarez/go-sqlite"
)

type UserStore struct {
	DB *sql.DB
}

type User struct {
	ID             int
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type CreateUserInput struct {
	Email          string
	HashedPassword []byte
}

func (s *UserStore) Insert(input CreateUserInput) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO users (email, hashed_password, created)
    VALUES(?, ?, ?)`

	_, err := s.DB.Exec(query, input.Email, string(input.HashedPassword), now)
	if err != nil {
		if sqliteError, ok := errors.AsType[*sqlite.Error](err); ok {
			if sqliteError.Code() == 2067 && strings.Contains(sqliteError.Error(), "UNIQUE") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserStore) GetByEmail(email string) (*User, error) {
	var u User
	var createdAt string

	query := "SELECT id, email, hashed_password, created FROM users WHERE email = ?"

	err := m.DB.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.HashedPassword, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	u.Created, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *UserStore) Exists(id int) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}
