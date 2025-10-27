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

	handler.TestRoute("/")
}
