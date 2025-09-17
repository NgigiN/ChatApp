package handlers

import (
	"net/http"
	"time"

	"chat_app/pkg/errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	Version   string      `json:"version"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

const API_VERSION = "v1"

func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func CreatedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func ErrorResponse(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.NewInternalError("Internal server error", err)
	}

	c.JSON(appErr.HTTPStatus, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func ValidationErrorResponse(c *gin.Context, message string, details string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: details,
		},
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func ForbiddenResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "FORBIDDEN",
			Message: message,
		},
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "NOT_FOUND",
			Message: message,
		},
		Timestamp: getCurrentTimestamp(),
		Version:   API_VERSION,
	})
}

func getCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
