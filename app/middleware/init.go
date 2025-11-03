package middleware

import (
	"app/domain"
	"os"

	"github.com/gin-gonic/gin"
)

type appMiddleware struct {
	publicKey  string
	repository domain.AppRepository
}

type AppMiddleware interface {
	Auth() gin.HandlerFunc
}

func NewAppMiddleware(repo domain.AppRepository) AppMiddleware {
	return &appMiddleware{
		publicKey:  os.Getenv("JWT_SECRET"),
		repository: repo,
	}
}
