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
		api.POST("", r.CreateCustomerConversation)
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
