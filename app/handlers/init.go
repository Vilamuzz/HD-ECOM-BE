package handlers

import (
	"app/app/middleware"
	"app/app/websocket"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type appRoute struct {
	Service    domain.AppService
	Repository domain.AppRepository
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
	Hub        *domain.Hub
	HubService domain.HubService
}

func App(service domain.AppService, repo domain.AppRepository, route *gin.Engine, middleware middleware.AppMiddleware) {
	hub := websocket.NewHub()
	hubService := websocket.NewHubService(hub)
	go hubService.Run()

	handler := &appRoute{
		Service:    service,
		Repository: repo,
		Route:      route.Group("/api"),
		Middleware: middleware,
		Hub:        hub,
		HubService: hubService,
	}

	handler.WebSocketRoute("/ws")
	handler.ConversationRoute("/conversations")
}
