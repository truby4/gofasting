package web

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddFieldErrors(errors map[string]string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	for key, message := range errors {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}
