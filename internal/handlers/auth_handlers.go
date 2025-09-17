package handlers

import (
	"chat_app/internal/models"
	"chat_app/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	authService services.AuthService
}

func NewAuthHandlers(authService services.AuthService) *AuthHandlers {
	return &AuthHandlers{authService: authService}
}

func (h *AuthHandlers) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	CreatedResponse(c, response, "User registered successfully")
}

func (h *AuthHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, response, "Login successful")
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	// Get token from context (set by auth middleware)
	token := c.GetString("token")
	if token == "" {
		UnauthorizedResponse(c, "No token provided")
		return
	}

	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, nil, "Logout successful")
}

func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, response, "Token refreshed successfully")
}
