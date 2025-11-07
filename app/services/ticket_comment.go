package services

import "app/domain/models"

func (s *appService) CreateTicketComment(comment *models.TicketComment) error {
	return s.repo.CreateTicketComment(comment)
}

func (s *appService) GetTicketComments() ([]models.TicketComment, error) {
	return s.repo.GetTicketComments()
}

func (s *appService) GetTicketCommentByID(id int) (*models.TicketComment, error) {
	return s.repo.GetTicketCommentByID(id)
}

func (s *appService) GetTicketCommentsByTicketID(ticketID int) ([]models.TicketComment, error) {
	return s.repo.GetTicketCommentsByTicketID(ticketID)
}

func (s *appService) UpdateTicketComment(comment *models.TicketComment) error {
	return s.repo.UpdateTicketComment(comment)
}

func (s *appService) DeleteTicketComment(id int) error {
	return s.repo.DeleteTicketComment(id)
}