package handlers

import (
	"app/app/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type appRoute struct {
	Service    domain.AppService
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
}

func App(service domain.AppService, route *gin.Engine, middleware middleware.AppMiddleware) {
	handler := &appRoute{
		Service:    service,
		Route:      route.Group("/api"),
		Middleware: middleware,
	}

	handler.WebSocketRoute("/ws")
	handler.ConversationRoute("/conversations")
	handler.MessageRoute("/conversations/:id/messages")
	handler.AdminConversationStatesRoute("/conversations/notifications")
	handler.TicketCategoryRoutes(handler.Route)
	handler.TicketPriorityRoutes(handler.Route)
	handler.TicketStatusRoutes(handler.Route)
	handler.TicketRoutes(handler.Route)
	handler.TicketAssignmentRoutes(handler.Route)
	handler.TicketAttachmentRoutes(handler.Route)
	handler.TicketCommentRoutes(handler.Route)
	handler.TicketLogRoutes(handler.Route)
	handler.Route.GET("/me", handler.Middleware.Auth(), handler.GetCurrentUser)
    handler.Route.GET("/users/support", handler.Middleware.Auth(), handler.GetSupportUsers)
}
