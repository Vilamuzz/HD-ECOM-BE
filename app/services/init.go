package services

import "app/domain"

type appService struct {
	repo domain.AppRepository
}

func NewAppService(repo domain.AppRepository) domain.AppService {
	return &appService{repo: repo}
}
