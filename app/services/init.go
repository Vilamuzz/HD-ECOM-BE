package services

import "app/domain"

type appService struct {
	repo domain.AppRepository
	hub  *domain.Hub
}

func NewAppService(repo domain.AppRepository, hub *domain.Hub) domain.AppService {
	return &appService{repo: repo, hub: hub}
}
