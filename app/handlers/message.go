package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) MessageRoute(path string) {
	api := r.Route.Group(path)
	api.Use(r.Middleware.Auth())
	{
		api.GET("", r.GetMessageHistory)
	}
}

// GetConversationMessages godoc
// @Summary      Get conversation messages
// @Description  Get all messages for a specific conversation
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Param        id   path      int  true  "Conversation ID"
// @Success      200  {object}   helpers.Response{data=[]models.Message}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations/{id}/messages [get]
func (r *appRoute) GetMessageHistory(c *gin.Context) {
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid conversation ID", nil, nil))
		return
	}

	// Get user from context to check if admin
	claim, _ := c.MustGet("userData").(models.User)
	isAdmin := claim.Role == models.RoleAdmin

	// Authorization check - verify user can access this conversation
	conversation, err := r.Service.GetConversationByID(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewResponse(http.StatusNotFound, "Conversation not found", nil, nil))
		return
	}

	// Check if user has permission to access this conversation
	if !isAdmin && conversation.CustomerID != claim.ID {
		c.JSON(http.StatusForbidden, helpers.NewResponse(http.StatusForbidden, "Access denied - not your conversation", nil, nil))
		return
	}

	// parse query params for cursor pagination
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	cursor := c.Query("cursor")

	response := r.Service.GetMessageHistory(conversationID, limit, cursor, isAdmin)
	c.JSON(response.Status, response)
}
