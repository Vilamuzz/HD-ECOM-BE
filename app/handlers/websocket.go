package handlers

import (
	"app/app/websocket"
)

func (r *appRoute) WebSocketRoute(path string) {
	api := r.Route.Group(path)
	api.GET("", r.Middleware.Auth(), websocket.ServeWebSocket(r.Hub, r.Repository))
}
