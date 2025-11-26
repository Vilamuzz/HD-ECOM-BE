package services

import "app/domain/models"

func (s *appService) CreateTicketLog(log *models.TicketLog) error {
	return s.repo.CreateTicketLog(log)
}

func (s *appService) GetTicketLogs() ([]models.TicketLog, error) {
	return s.repo.GetTicketLogs()
}

func (s *appService) GetTicketLogByID(id int) (*models.TicketLog, error) {
	return s.repo.GetTicketLogByID(id)
}

func (s *appService) GetTicketLogsByTicketID(ticketID int) ([]models.TicketLog, error) {
	return s.repo.GetTicketLogsByTicketID(ticketID)
}