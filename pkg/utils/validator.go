package utils

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
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

func ValidatorErrors(err error) string {
	fields := make([]string, 0)
	for _, err := range err.(validator.ValidationErrors) {
		fields = append(fields, fmt.Sprintf("%s: %s", err.Field(), err.ActualTag()))
	}
	return strings.Join(fields, "\n")
}
