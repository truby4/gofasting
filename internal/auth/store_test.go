package auth

import (
	"database/sql"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			hashed_password TEXT NOT NULL,
			created TEXT NOT NULL,
			updated TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestCreate(t *testing.T) {
	db := newTestDB(t)
	s := store{db: db}

	err := s.create("a@example.com", []byte("password"))
	if err != nil {
		t.Fatal(err)
	}
	err = s.create("a@example.com", []byte("password2"))
	assertError(t, err, ErrDuplicateEmail)
}

func TestGetByEmail(t *testing.T) {
	db := newTestDB(t)
	s := store{db: db}

	err := s.create("a@example.com", []byte("hashed-password"))
	if err != nil {
		t.Fatal(err)
	}

	u, err := s.getByEmail("a@example.com")
	if err != nil {
		t.Fatal(err)
	}

	assertStrings(t, u.Email, "a@example.com")
	assertStrings(t, string(u.HashedPassword), "hashed-password")

	_, err = s.getByEmail("b@example.com")
	assertError(t, err, ErrNoRecordFound)
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}
