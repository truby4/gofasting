package auth

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type input struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=8"`
}

var validate = validator.New()

func (i *input) validateInput() error {
	err := validate.Struct(i)
	if err != nil {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			verrs := make(map[string]string)
			for _, e := range validateErrs {
				verrs[strings.ToLower(e.Field())] = validationMessage(e)
			}
			return ValidationError{Fields: verrs}
		}
		return err
	}
	return nil
}

func validationMessage(e validator.FieldError) string {
	switch e.Field() {
	case "Email":
		switch e.Tag() {
		case "required":
			return "email is required"
		case "email":
			return "please enter a valid email"
		}

	case "Password":
		switch e.Tag() {
		case "required":
			return "password is required"
		case "gte":
			return "please use 8 or more characters for password"
		}
	}

	return e.Error()
}
