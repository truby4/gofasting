package auth

import "testing"

func TestAuthenticate(t *testing.T) {
	db := newTestDB(t)
	s := New(db)

	err := s.Register("a@example.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	// ensure validation fails first
	_, err = s.Authenticate("a@example.com", "")
	assertValidationError(t, err, ErrValidation)
}

func assertValidationError(t testing.TB, got, want error) {
	t.Helper()

	if got.Error() != want.Error() {
		t.Errorf("got error %q want %q", got, want)
	}
}

func TestServiceRegister(t *testing.T) {
	db := newTestDB(t)
	s := New(db)

	err := s.Register("a@example.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	u, err := s.store.getByEmail("a@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if u.Email != "a@example.com" {
		t.Fatalf("got email %q", u.Email)
	}

	if string(u.HashedPassword) == "password123" {
		t.Fatal("password was stored unhashed")
	}
}
