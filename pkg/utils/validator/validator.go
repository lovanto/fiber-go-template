package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err == nil {
			return true // valid UUID
		}
		return false // invalid UUID
	})
	return validate
}

func ValidatorErrors(err error) string {
	if err == nil {
		return ""
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ""
	}
	fields := make([]string, 0)
	for _, ferr := range verrs {
		fields = append(fields, fmt.Sprintf("%s: %s", ferr.Field(), ferr.ActualTag()))
	}
	return strings.Join(fields, "\n")
}
