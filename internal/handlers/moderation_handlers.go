package handlers

import (
	"strconv"

	"chat_app/internal/services"

	"github.com/gin-gonic/gin"
)

type ModerationHandlers struct {
	roomService services.RoomService
	userService services.UserService
}

func NewModerationHandlers(roomService services.RoomService, userService services.UserService) *ModerationHandlers {
	return &ModerationHandlers{
		roomService: roomService,
		userService: userService,
	}
}

// RemoveUser removes a user from a room (only room creator can do this)
func (h *ModerationHandlers) RemoveUser(c *gin.Context) {
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

	// Get the user to remove by username
	userToRemove, err := h.userService.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	// Check if the requesting user is the room creator
	room, err := h.roomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	if room.CreatedBy != userIDInt {
		ForbiddenResponse(c, "Only room creator can remove users")
		return
	}

	// Remove the user from the room
	err = h.roomService.LeaveRoom(c.Request.Context(), roomID, userToRemove.ID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, nil, "User removed from room successfully")
}

// ResetRoom removes all members except the creator (for semester end)
func (h *ModerationHandlers) ResetRoom(c *gin.Context) {
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
		ForbiddenResponse(c, "Only room creator can reset room")
		return
	}

	// Get all room members
	members, err := h.roomService.GetRoomMembers(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	// Remove all members except the creator
	for _, member := range members {
		if member.ID != userIDInt {
			err = h.roomService.LeaveRoom(c.Request.Context(), roomID, member.ID)
			if err != nil {
				// Log error but continue with other members
				continue
			}
		}
	}

	SuccessResponse(c, nil, "Room reset successfully - all members removed except creator")
}

// GetRoomPermissions returns the user's permissions in a room
func (h *ModerationHandlers) GetRoomPermissions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	room, err := h.roomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	permissions := map[string]bool{
		"is_creator": room.CreatedBy == userIDInt,
		"can_moderate": room.CreatedBy == userIDInt,
		"can_remove_users": room.CreatedBy == userIDInt,
		"can_reset_room": room.CreatedBy == userIDInt,
	}

	SuccessResponse(c, permissions, "Room permissions retrieved successfully")
}
