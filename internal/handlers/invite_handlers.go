package handlers

import (
	"strconv"

	"chat_app/internal/services"

	"github.com/gin-gonic/gin"
)

type InviteHandlers struct {
	roomService services.RoomService
	userService services.UserService
}

func NewInviteHandlers(roomService services.RoomService, userService services.UserService) *InviteHandlers {
	return &InviteHandlers{
		roomService: roomService,
		userService: userService,
	}
}

// InviteUser invites a user to a private room
func (h *InviteHandlers) InviteUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	// Check if the requesting user is the room creator
	room, err := h.roomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	if room.CreatedBy != userIDInt {
		ForbiddenResponse(c, "Only room creator can invite users")
		return
	}

	// Get the user to invite by username
	userToInvite, err := h.userService.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	// Add the user to the room
	err = h.roomService.JoinRoom(c.Request.Context(), roomID, userToInvite.ID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, nil, "User invited to room successfully")
}

// InviteMultipleUsers invites multiple users to a private room
func (h *InviteHandlers) InviteMultipleUsers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	var req struct {
		Usernames []string `json:"usernames" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	// Check if the requesting user is the room creator
	room, err := h.roomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	if room.CreatedBy != userIDInt {
		ForbiddenResponse(c, "Only room creator can invite users")
		return
	}

	var successCount int
	var failedInvites []string

	for _, username := range req.Usernames {
		// Get the user to invite by username
		userToInvite, err := h.userService.GetUserByUsername(c.Request.Context(), username)
		if err != nil {
			failedInvites = append(failedInvites, username)
			continue
		}

		// Add the user to the room
		err = h.roomService.JoinRoom(c.Request.Context(), roomID, userToInvite.ID)
		if err != nil {
			failedInvites = append(failedInvites, username)
			continue
		}

		successCount++
	}

	response := map[string]interface{}{
		"success_count": successCount,
		"failed_invites": failedInvites,
		"total_invited": len(req.Usernames),
	}

	SuccessResponse(c, response, "Bulk invite completed")
}

// GetInvitableUsers returns users that can be invited to the room
func (h *InviteHandlers) GetInvitableUsers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	// Check if the requesting user is the room creator
	room, err := h.roomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	if room.CreatedBy != userIDInt {
		ForbiddenResponse(c, "Only room creator can view invitable users")
		return
	}

	// Get all users (simplified - in production, you'd want pagination and search)
	users, err := h.userService.GetAllUsers(c.Request.Context(), 100, 0)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	// Get current room members
	members, err := h.roomService.GetRoomMembers(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	// Create a map of member IDs for quick lookup
	memberMap := make(map[int]bool)
	for _, member := range members {
		memberMap[member.ID] = true
	}

	// Filter out users who are already members
	var invitableUsers []map[string]interface{}
	for _, user := range users {
		if !memberMap[user.ID] {
			invitableUsers = append(invitableUsers, map[string]interface{}{
				"id": user.ID,
				"username": user.Username,
			})
		}
	}

	SuccessResponse(c, invitableUsers, "Invitable users retrieved successfully")
}
