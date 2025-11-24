package services

import (
	"app/domain"
	"time"
)

type appService struct {
	repo    domain.AppRepository
	hub     *domain.Hub
	timeout time.Duration
}

type DBInjection struct {
	Repo domain.AppRepository
}

func NewAppService(repoInjection DBInjection, hub *domain.Hub, timeout time.Duration) domain.AppService {
	return &appService{
		repo:    repoInjection.Repo,
		hub:     hub,
		timeout: timeout,
	}
}
