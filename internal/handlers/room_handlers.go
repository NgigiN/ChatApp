package handlers

import (
	"strconv"

	"chat_app/internal/services"

	"github.com/gin-gonic/gin"
)

type RoomHandlers struct {
	roomService services.RoomService
	userService services.UserService
}

func NewRoomHandlers(roomService services.RoomService, userService services.UserService) *RoomHandlers {
	return &RoomHandlers{
		roomService: roomService,
		userService: userService,
	}
}

// CreateRoom creates a new room
func (h *RoomHandlers) CreateRoom(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsPrivate   bool   `json:"is_private"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	room, err := h.roomService.CreateRoom(c.Request.Context(), userIDInt, req.Name, req.Description, req.IsPrivate)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	CreatedResponse(c, room, "Room created successfully")
}

// JoinRoom allows a user to join a room
func (h *RoomHandlers) JoinRoom(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	err = h.roomService.JoinRoom(c.Request.Context(), roomID, userIDInt)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, nil, "Successfully joined room")
}

// LeaveRoom allows a user to leave a room
func (h *RoomHandlers) LeaveRoom(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	err = h.roomService.LeaveRoom(c.Request.Context(), roomID, userIDInt)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, nil, "Successfully left room")
}

// GetUserRooms returns all rooms the user is a member of
func (h *RoomHandlers) GetUserRooms(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	rooms, err := h.roomService.GetUserRooms(c.Request.Context(), userIDInt)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, rooms, "User rooms retrieved successfully")
}

// GetRoom returns room details
func (h *RoomHandlers) GetRoom(c *gin.Context) {
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

	SuccessResponse(c, room, "Room details retrieved successfully")
}

// UpdateRoom allows room creator to update room details
func (h *RoomHandlers) UpdateRoom(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	var req struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		IsPrivate   *bool  `json:"is_private,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsPrivate != nil {
		updates["is_private"] = *req.IsPrivate
	}

	room, err := h.roomService.UpdateRoom(c.Request.Context(), roomID, userIDInt, updates)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, room, "Room updated successfully")
}

// DeleteRoom allows room creator to delete a room
func (h *RoomHandlers) DeleteRoom(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int)

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	err = h.roomService.DeleteRoom(c.Request.Context(), roomID, userIDInt)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, nil, "Room deleted successfully")
}

// GetRoomMembers returns all members of a room
func (h *RoomHandlers) GetRoomMembers(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		ValidationErrorResponse(c, "Invalid room ID", err.Error())
		return
	}

	members, err := h.roomService.GetRoomMembers(c.Request.Context(), roomID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, members, "Room members retrieved successfully")
}
