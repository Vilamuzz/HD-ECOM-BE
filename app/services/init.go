package services

import (
	"app/domain"
	"time"
)

type appService struct {
	repo    domain.AppRepository
	hub     *domain.Hub
	timeout time.Duration
	s3Repo  domain.S3Repository
}

type DBInjection struct {
	Repo   domain.AppRepository
	S3Repo domain.S3Repository
}

func NewAppService(repoInjection DBInjection, hub *domain.Hub, timeout time.Duration) domain.AppService {
	return &appService{
		repo:    repoInjection.Repo,
		s3Repo:  repoInjection.S3Repo,
		hub:     hub,
		timeout: timeout,
	}
}
