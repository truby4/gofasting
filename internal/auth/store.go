package auth

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/glebarez/go-sqlite"
)

type userRecord struct {
	ID             int
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type store struct {
	db *sql.DB
}

func (s *store) create(email string, hashedPassword []byte) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO users (email, hashed_password, created, updated)
    VALUES(?, ?, ?, ?)`

	_, err := s.db.Exec(query, email, string(hashedPassword), now, now)
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

func (s *store) getByEmail(email string) (*userRecord, error) {
	var u userRecord
	var createdAt string
	var updatedAt string

	query := "SELECT id, email, hashed_password, created, updated FROM users WHERE email = ?"

	err := s.db.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.HashedPassword, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecordFound
		} else {
			return nil, err
		}
	}

	u.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}

	u.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *store) exists(id int) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := s.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}
