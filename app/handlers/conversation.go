package handlers

import (
	"app/domain/models"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) ConversationRoute(path string) {
	api := r.Route.Group(path)
	api.Use(r.Middleware.Auth())
	{
		api.GET("", r.GetConversations)
		api.GET("/:id/messages", r.GetConversationMessages)
	}
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
	userData, _ := c.Get("userData")
	user := userData.(models.User)

	response := r.Service.GetConversations(user)
	c.JSON(response.Status, response)
}

// GetConversationMessages godoc
// @Summary      Get conversation messages
// @Description  Get all messages for a specific conversation
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Param        id   path      int  true  "Conversation ID"
// @Success      200  {object}   helpers.Response{data=[]models.ChatMessage}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations/{id}/messages [get]
func (r *appRoute) GetConversationMessages(c *gin.Context) {
	conversationID := c.Param("id")

	response := r.Service.GetConversationMessages(conversationID)
	c.JSON(response.Status, response)
}
