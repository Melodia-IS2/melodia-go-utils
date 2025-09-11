package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code     string
	Title    string
	Message  string
	HTTPCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(resource string, identifier any) *AppError {
	if identifier == nil {
		return &AppError{
			Code:     "NOT_FOUND",
			Title:    fmt.Sprintf("%s Not Found", resource),
			Message:  fmt.Sprintf("The %s was not found", resource),
			HTTPCode: http.StatusNotFound,
		}
	}
	return &AppError{
		Code:     "NOT_FOUND",
		Title:    fmt.Sprintf("%s Not Found", resource),
		Message:  fmt.Sprintf("The %s with the identifier %v was not found", resource, identifier),
		HTTPCode: http.StatusNotFound,
	}
}

func NewValidationError(msg string) *AppError {
	return &AppError{
		Code:     "VALIDATION_ERROR",
		Title:    "Validation Error",
		Message:  msg,
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func NewConflictError(resource string, identifier any) *AppError {
	if identifier == nil {
		return &AppError{
			Code:     "CONFLICT",
			Title:    fmt.Sprintf("%s Conflict", resource),
			Message:  fmt.Sprintf("The %s already exists", resource),
			HTTPCode: http.StatusConflict,
		}
	}
	return &AppError{
		Code:     "CONFLICT",
		Title:    fmt.Sprintf("%s Conflict", resource),
		Message:  fmt.Sprintf("The %s with the identifier %v already exists", resource, identifier),
		HTTPCode: http.StatusConflict,
	}
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{
		Code:     "BAD_REQUEST",
		Title:    "Bad Request",
		Message:  msg,
		HTTPCode: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(msg string) *AppError {
	return &AppError{
		Code:     "UNAUTHORIZED",
		Title:    "Unauthorized",
		Message:  msg,
		HTTPCode: http.StatusUnauthorized,
	}
}

func NewTokenExpiredError(msg string) *AppError {
	return &AppError{
		Code:     "TOKEN_EXPIRED",
		Title:    "Token Expired",
		Message:  msg,
		HTTPCode: http.StatusUnauthorized,
	}
}
