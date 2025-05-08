package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct
func Validate(v interface{}) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	// Format validation errors
	var errMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		switch err.Tag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("%s is required", field))
		case "email":
			errMessages = append(errMessages, fmt.Sprintf("%s must be a valid email", field))
		case "min":
			errMessages = append(errMessages, fmt.Sprintf("%s must be at least %s characters", field, err.Param()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("%s is invalid", field))
		}
	}

	return fmt.Errorf(strings.Join(errMessages, ", "))
}
