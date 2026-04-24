package store

import (
	"database/sql"

	"github.com/charmbracelet/log"
)

type Store struct {
	DB     *sql.DB
	logger *log.Logger
	Fasts  *FastStore
}

func New(logger *log.Logger) (*Store, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	return &Store{
		DB:     db,
		logger: logger.WithPrefix("STORE"),
		Fasts:  &FastStore{DB: db},
	}, nil
}

func (s *Store) Close() error {
	return s.DB.Close()
}
