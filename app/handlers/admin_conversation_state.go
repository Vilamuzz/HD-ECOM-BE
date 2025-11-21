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

func (r *appRoute) GetAdminListConversationStates(c *gin.Context) {
	claim, _ := c.MustGet("userData").(models.User)

	response := r.Service.GetAdminListConversationStates(claim)
	c.JSON(response.Status, response)
}
