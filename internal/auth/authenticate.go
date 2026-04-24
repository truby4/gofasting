package auth

import (
	"errors"

	"github.com/truby4/go-fasting/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func validateAuthenticationForm(f *AuthenticationForm) {
	f.CheckField(validator.NotBlank(f.Email), "email", "This field cannot be blank")
	f.CheckField(validator.Matches(f.Email, validator.EmailRX), "email", "This field must be a valid email address")
	f.CheckField(validator.NotBlank(f.Password), "password", "This field cannot be blank")
	f.CheckField(validator.MaxBytes(f.Password, 72), "password", "This field must not be more than 72 bytes long")
}

func (s *Service) Authenticate(f *AuthenticationForm) (int, error) {
	validateAuthenticationForm(f)

	if !f.Valid() {
		return 0, ErrValidation
	}

	u, err := s.store.GetByEmail(f.Email)
	if err != nil {
		// Best error? Don't want to say if password wrong or not?
		return 0, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(f.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return u.ID, nil
}
