package handlers

import (
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TestRoute(path string) {
	api := r.Route.Group(path)
	api.GET("/test", r.Middleware.Auth(), r.testHandler)
}

// Test Endpoint godoc
// @Summary      Test Endpoint
// @Description  Test JWT
// @Security 	 BearerAuth
// @Tags         users
// @Produce      json
// @Success      200  {object}   helpers.Response
// @Router       /test [get]
func (r *appRoute) testHandler(c *gin.Context) {
	response := helpers.NewResponse(http.StatusOK, "Test successful", nil, nil)
	c.JSON(http.StatusOK, response)
}
