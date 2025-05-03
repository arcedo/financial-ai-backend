package utils

import (
	"errors"
	"net/http"
	"net/mail"
	"slices"
	"strings"
)

func ValidateStringField(field string, fieldName string) error {
	if strings.TrimSpace(field) == "" {
		return errors.New(fieldName + " cannot be empty")
	}
	return nil
}

func ValidateEmail(email string) error {
	if err := ValidateStringField(email, "email"); err != nil {
		return err
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if err := ValidateStringField(password, "password"); err != nil {
		return err
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func ValidateMethods(r *http.Request, allowedMethods []string) error {
	if slices.Contains(allowedMethods, r.Method) {
		return nil
	}
	return ErrInvalidMethod
}
