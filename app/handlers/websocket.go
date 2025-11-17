package handlers

import (
	"github.com/gin-gonic/gin"
)

func (r *appRoute) WebSocketRoute(path string) {
	api := r.Route.Group(path)
	api.GET("", r.Middleware.Auth(), r.ServeWebSocket)
}

func (r *appRoute) ServeWebSocket(c *gin.Context) {
	r.Service.ServeWebSocket(c)
}
