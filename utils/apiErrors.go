package utils

import (
	"errors"
	"net/http"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	// General errors
	ErrInternalServer     = errors.New("INTERNAL_SERVER_ERROR")
	ErrNotFound           = errors.New("NOT_FOUND")
	ErrInvalidInput       = errors.New("INVALID_INPUT")
	ErrInvalidMethod      = errors.New("METHOD_NOT_ALLOWED")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")

	// Specific errors
	ErrDatabase     = errors.New("DATABASE_ERROR")
	ErrCache        = errors.New("CACHE_ERROR")
	ErrUnauthorized = errors.New("UNAUTHORIZED")
	ErrForbidden    = errors.New("FORBIDDEN")
)

// Predefined APIError objects with messages
var ErrorMap = map[error]APIError{
	ErrInternalServer:     {Code: "INTERNAL_SERVER_ERROR", Message: "An unexpected error occurred"},
	ErrNotFound:           {Code: "NOT_FOUND", Message: "The requested resource was not found"},
	ErrInvalidInput:       {Code: "INVALID_INPUT", Message: "The input provided is invalid"},
	ErrInvalidMethod:      {Code: "METHOD_NOT_ALLOWED", Message: "The HTTP method is not allowed for this endpoint"},
	ErrDatabase:           {Code: "DATABASE_ERROR", Message: "A database error occurred"},
	ErrCache:              {Code: "CACHE_ERROR", Message: "A cache error occurred"},
	ErrUnauthorized:       {Code: "UNAUTHORIZED", Message: "You are not authorized to access this resource"},
	ErrForbidden:          {Code: "FORBIDDEN", Message: "You don't have permission to perform this action"},
	ErrInvalidCredentials: {Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"},
}

// HTTP status codes for predefined errors
var StatusMap = map[error]int{
	ErrInternalServer:     http.StatusInternalServerError,
	ErrNotFound:           http.StatusNotFound,
	ErrInvalidInput:       http.StatusBadRequest,
	ErrInvalidMethod:      http.StatusMethodNotAllowed,
	ErrDatabase:           http.StatusInternalServerError,
	ErrCache:              http.StatusInternalServerError,
	ErrUnauthorized:       http.StatusUnauthorized,
	ErrForbidden:          http.StatusForbidden,
	ErrInvalidCredentials: http.StatusBadRequest,
}

func MapErrorToAPIError(err error) (*APIError, int) {
	if apiErr, found := ErrorMap[err]; found {
		return &apiErr, StatusMap[err]
	}
	// Default fallback for unexpected errors
	genericErr := APIError{
		Code:    "GENERIC",
		Message: err.Error(),
	}
	return &genericErr, http.StatusBadRequest
}
