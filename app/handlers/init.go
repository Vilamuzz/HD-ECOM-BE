package handlers

import (
	"app/app/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type appRoute struct {
	Service    domain.AppService
	Repository domain.AppRepository
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
}

func App(service domain.AppService, repo domain.AppRepository, route *gin.Engine, middleware middleware.AppMiddleware) {
	handler := &appRoute{
		Service:    service,
		Repository: repo,
		Route:      route.Group("/api"),
		Middleware: middleware,
	}

	handler.WebSocketRoute("/ws")
	handler.ConversationRoute("/conversations")
}
