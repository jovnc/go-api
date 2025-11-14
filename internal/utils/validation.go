package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate validates the request body
func Validate(i interface{}) error {
	err := validate.Struct(i)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error())
		}
		return errors.New(strings.Join(errorMessages, ", "))
	}
	return nil
}
