package app

import (
	"database/sql"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/dev.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, err
	}

	err = initDB(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initDB(db *sql.DB) error {
	schema, err := os.ReadFile("./data/schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	return err
}
