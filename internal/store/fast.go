package store

import (
	"database/sql"
)

type FastStore struct {
	DB *sql.DB
}
