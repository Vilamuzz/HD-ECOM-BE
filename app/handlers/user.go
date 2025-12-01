package handlers

import (
    "app/helpers"
    "net/http"
    "github.com/gin-gonic/gin"
)

// GetCurrentUser godoc
// @Summary      Get current user info
// @Description  Get information about the currently authenticated user
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.Response{data=object}
// @Failure      401 {object} helpers.Response
// @Router       /me [get]
func (r *appRoute) GetCurrentUser(c *gin.Context) {
    user, exists := c.Get("userData")
    if !exists {
        c.JSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Unauthorized", nil, nil))
        return
    }
    c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Current user info", nil, user))
}