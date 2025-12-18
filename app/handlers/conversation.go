package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) ConversationRoute(path string) {
	api := r.Route.Group(path)
	api.Use(r.Middleware.Auth())
	{
		api.GET("", r.GetConversations)
		api.POST("", r.CreateCustomerConversation)
		api.POST("/:id/close", r.CloseConversation)
	}
}

// CreateCustomerConversation godoc
// @Summary      Create conversation
// @Description  Create a new conversation between a customer and an admin
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Success      200  {object}   helpers.Response{data=models.Conversation}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations [post]
func (r *appRoute) CreateCustomerConversation(c *gin.Context) {
	ctx := c.Request.Context()
	claim, _ := c.MustGet("userData").(models.User)

	response := r.Service.CreateCustomerConversation(ctx, claim)
	c.JSON(response.Status, response)
}

// GetConversations godoc
// @Summary      Get conversations
// @Description  Get all conversations for the current user (admin sees all, customer sees their own)
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Success      200  {object}   helpers.Response{data=[]models.Conversation}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations [get]
func (r *appRoute) GetConversations(c *gin.Context) {
	claim, _ := c.MustGet("userData").(models.User)

	response := r.Service.GetConversations(claim)
	c.JSON(response.Status, response)
}

// CloseConversation godoc
// @Summary      Close conversation
// @Description  Close an existing conversation
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Success      200  {object}   helpers.Response{data=models.Conversation}
// @Failure      500  {object}   helpers.Response
// @Param        id   path      string  true  "Conversation ID"
// @Router       /conversations/{id}/close [post]
func (r *appRoute) CloseConversation(c *gin.Context) {
	claim, _ := c.MustGet("userData").(models.User)
	ctx := c.Request.Context()
	id := c.Param("id")

	// Parse conversation ID
	conversationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid conversation ID", nil, nil))
		return
	}

	// Authorization check - verify user can close this conversation
	conversation, err := r.Service.GetConversationByID(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewResponse(http.StatusNotFound, "Conversation not found", nil, nil))
		return
	}

	// Only admins or conversation participants can close conversations
	isAdmin := claim.Role == models.RoleAdmin
	isParticipant := conversation.CustomerID == claim.ID || conversation.AdminID == claim.ID

	if !isAdmin && !isParticipant {
		c.JSON(http.StatusForbidden, helpers.NewResponse(http.StatusForbidden, "Access denied - not your conversation", nil, nil))
		return
	}

	response := r.Service.CloseConversation(ctx, claim, id)
	c.JSON(response.Status, response)
}
