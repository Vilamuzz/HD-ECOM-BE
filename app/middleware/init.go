package middleware

import (
	"app/domain"
	"os"
)

type appMiddleware struct {
	publicKey  string
	repository domain.AppRepository
}


func NewAppMiddleware(repo domain.AppRepository) AppMiddleware {
	return &appMiddleware{
		publicKey:  os.Getenv("JWT_SECRET"),
		repository: repo,
	}
}
