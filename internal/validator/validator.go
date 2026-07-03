package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *Validator) Validate(s interface{}) []ValidationError {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(e.Field()),
			Message: translateTag(e),
		})
	}

	return errors
}

func translateTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "esta campo es obligatorio"
	case "email":
		return "debe ser un correo válido"
	case "min":
		return fmt.Sprintf("debe tener mínimo %s caracteres", e.Param())
	case "max":
		return fmt.Sprintf("debe tener máximo %s caracteres", e.Param())
	default:
		return "valor inválido"
	}
}
