package services

import "app/domain/models"

func (s *appService) CreateTicket(ticket *models.Ticket) error {
	return s.repo.CreateTicket(ticket)
}

func (s *appService) GetTickets() ([]models.Ticket, error) {
	return s.repo.GetTickets()
}

func (s *appService) GetTicketByID(id int) (*models.Ticket, error) {
	return s.repo.GetTicketByID(id)
}

func (s *appService) UpdateTicket(ticket *models.Ticket) error {
	return s.repo.UpdateTicket(ticket)
}

func (s *appService) DeleteTicket(id int) error {
	return s.repo.DeleteTicket(id)
}
