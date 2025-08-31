package api

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func TransformValidationErrors(err error) []string {
	var errorMessages []string
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range fieldErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is a required field", e.Field()))
			case "alpha":
				errorMessages = append(errorMessages, fmt.Sprintf("%s can only contain non numeric strings", e.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s requires a min length of %s", e.Field(), e.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s exceeds the maximum length of %s", e.Field(), e.Param()))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is not a valid email address", e.Field()))
			case "strong":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is not strong enough", e.Field()))
			}
		}
	}

	return errorMessages
}

func IsPasswordStrong(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch), unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}
