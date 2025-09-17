package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrCodeConflict          ErrorCode = "CONFLICT"
	ErrCodeInternalError     ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabaseError     ErrorCode = "DATABASE_ERROR"
	ErrCodeValidationError   ErrorCode = "VALIDATION_ERROR"
	ErrCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
)

type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Cause      error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func NewInvalidInputError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Cause:      cause,
	}
}

func NewUnauthorizedError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
		Cause:      cause,
	}
}

func NewForbiddenError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
		Cause:      cause,
	}
}

func NewNotFoundError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
		Cause:      cause,
	}
}

func NewConflictError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
		Cause:      cause,
	}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeInternalError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Cause:      cause,
	}
}

func NewDatabaseError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Cause:      cause,
	}
}

func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeValidationError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Cause:      cause,
	}
}

func NewRateLimitError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeRateLimitExceeded,
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
		Cause:      cause,
	}
}
