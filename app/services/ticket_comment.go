package services

import (
	"app/domain/models"
	"fmt"
	"log"
	"time"
)

func (s *appService) CreateTicketComment(comment *models.TicketComment) error {
	// Get the ticket first to ensure it exists
	ticket, err := s.repo.GetTicketByID(comment.TicketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %v", err)
	}

	// Get ticket owner (user who created the ticket)
	user, err := s.repo.GetUserByID(ticket.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Create the comment first
	if err := s.repo.CreateTicketComment(comment); err != nil {
		return fmt.Errorf("failed to create comment: %v", err)
	}

	// Update ticket status to 3 (resolved/commented)
	ticket.StatusID = 3
	ticket.TanggalDiperbarui = time.Now()
	
	if err := s.repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to update ticket status: %v", err)
	}

	// Send email notification asynchronously
	go func() {
		emailErr := s.SendTicketCommentEmail(
			user.Email,
			user.Username,
			ticket.KodeTiket,
			ticket.Judul,
			comment.IsiPesan,
		)
		if emailErr != nil {
			log.Printf("Failed to send email notification for ticket #%s: %v", ticket.KodeTiket, emailErr)
		} else {
			log.Printf("Email notification sent successfully for ticket #%s to %s", ticket.KodeTiket, user.Email)
		}
	}()

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