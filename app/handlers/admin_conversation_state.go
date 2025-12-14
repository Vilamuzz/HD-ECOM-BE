package handlers

import (
	"app/domain/models"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) AdminConversationStatesRoute(path string) {
	api := r.Route.Group(path)
	api.Use(r.Middleware.Auth())
	{
		api.GET("", r.GetAdminListConversationStates)
	}
}

// GetAdminListConversationStates godoc
// @Summary      Get admin list conversation states and ticket notifications
// @Description  Get a list of conversation states for admin with open ticket counts by type
// @Security 	 BearerAuth
// @Tags         admin_conversation_states
// @Produce      json
// @Success      200  {object}   helpers.Response{data=map[string]interface{}}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations/notifications [get]
func (r *appRoute) GetAdminListConversationStates(c *gin.Context) {
	claim, _ := c.MustGet("userData").(models.User)

	response := r.Service.GetAdminListConversationStates(claim)
	c.JSON(response.Status, response)
}
