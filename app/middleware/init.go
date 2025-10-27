package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

type appMiddleware struct {
	publicKey string
}

type AppMiddleware interface {
	Auth() gin.HandlerFunc
}

func NewAppMiddleware() AppMiddleware {
	return &appMiddleware{
		publicKey: os.Getenv("JWT_SECRET"),
	}
}
