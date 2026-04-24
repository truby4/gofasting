package auth

import (
	"errors"

	"github.com/truby4/go-fasting/internal/store"
	"github.com/truby4/go-fasting/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func validateRegistrationForm(f *RegistrationForm) {
	f.CheckField(validator.NotBlank(f.Email), "email", "This field cannot be blank")
	f.CheckField(validator.Matches(f.Email, validator.EmailRX), "email", "This field must be a valid email address")
	f.CheckField(validator.NotBlank(f.Password), "password", "This field cannot be blank")
	f.CheckField(validator.MinChars(f.Password, 8), "password", "This field must be at least 8 characters long")
	f.CheckField(validator.MaxBytes(f.Password, 72), "password", "This field must not be more than 72 bytes long")
}

func (s *Service) Registration(f *RegistrationForm) error {
	validateRegistrationForm(f)

	if !f.Valid() {
		return ErrValidation
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 12)
	if err != nil {
		return err
	}

	input := store.CreateUserInput{
		Email:          f.Email,
		HashedPassword: hashedPassword,
	}

	err = s.store.Insert(input)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateEmail) {
			f.AddFieldError("email", "Email address is already in use")
			return ErrValidation
		}
		return err
	}
	return nil
}
