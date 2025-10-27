package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

type appMiddleware struct {
	publicKey string
}

type AppMiddleware interface {
	Validate() gin.HandlerFunc
}

func NewAppMiddleware() AppMiddleware {
	return &appMiddleware{
		publicKey: os.Getenv("JWT_PUBLIC_KEY"),
	}
}

func (m *appMiddleware) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Middleware logic to validate requests using m.publicKey
		c.Next()
	}
}
