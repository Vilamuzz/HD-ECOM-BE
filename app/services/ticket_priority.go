package services

import "app/domain/models"

func (s *appService) CreateTicketPriority(priority *models.TicketPriority) error {
	return s.repo.CreateTicketPriority(priority)
}

func (s *appService) GetTicketPriorities() ([]models.TicketPriority, error) {
	return s.repo.GetTicketPriorities()
}

func (s *appService) GetTicketPriorityByID(id int) (*models.TicketPriority, error) {
	return s.repo.GetTicketPriorityByID(id)
}

func (s *appService) UpdateTicketPriority(priority *models.TicketPriority) error {
	return s.repo.UpdateTicketPriority(priority)
}

func (s *appService) DeleteTicketPriority(id int) error {
	return s.repo.DeleteTicketPriority(id)
}
