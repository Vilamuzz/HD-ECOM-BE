package services

import "app/domain/models"

func (s *appService) CreateTicketStatus(status *models.TicketStatus) error {
	return s.repo.CreateTicketStatus(status)
}

func (s *appService) GetTicketStatuses() ([]models.TicketStatus, error) {
	return s.repo.GetTicketStatuses()
}

func (s *appService) GetTicketStatusByID(id int) (*models.TicketStatus, error) {
	return s.repo.GetTicketStatusByID(id)
}

func (s *appService) UpdateTicketStatus(status *models.TicketStatus) error {
	return s.repo.UpdateTicketStatus(status)
}

func (s *appService) DeleteTicketStatus(id int) error {
	return s.repo.DeleteTicketStatus(id)
}
