package services

import (
	"app/domain/models"
	"fmt"
	"time"
)

func (s *appService) CreateTicketComment(comment *models.TicketComment) error {
	// Get the ticket first to ensure it exists
	ticket, err := s.repo.GetTicketByID(comment.TicketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %v", err)
	}

	// Create the comment first
	if err := s.repo.CreateTicketComment(comment); err != nil {
		return fmt.Errorf("failed to create comment: %v", err)
	}

	// Update ticket status to 3 (status after comment is added)
	ticket.StatusID = 3
	ticket.TanggalDiperbarui = time.Now()
	
	if err := s.repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to update ticket status: %v", err)
	}

	return nil
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