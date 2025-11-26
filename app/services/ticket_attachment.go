package services

import (
       "os"
       "app/domain/models"
)

func (s *appService) CreateTicketAttachment(attachment *models.TicketAttachment) error {
	return s.repo.CreateTicketAttachment(attachment)
}

func (s *appService) GetTicketAttachments() ([]models.TicketAttachment, error) {
	return s.repo.GetTicketAttachments()
}

func (s *appService) GetTicketAttachmentByID(id int) (*models.TicketAttachment, error) {
	return s.repo.GetTicketAttachmentByID(id)
}

func (s *appService) GetTicketAttachmentsByTicketID(ticketID int) ([]models.TicketAttachment, error) {
	return s.repo.GetTicketAttachmentsByTicketID(ticketID)
}

func (s *appService) UpdateTicketAttachment(attachment *models.TicketAttachment) error {
	return s.repo.UpdateTicketAttachment(attachment)
}

func (s *appService) DeleteTicketAttachment(id int) error {
       // Get the attachment first
       attachment, err := s.repo.GetTicketAttachmentByID(id)
       if err != nil {
	       return err
       }
       // Delete from DB
       err = s.repo.DeleteTicketAttachment(id)
       if err != nil {
	       return err
       }
       // Delete file from disk
       if attachment != nil && attachment.FilePath != "" {
	       _ = os.Remove(attachment.FilePath) // Ignore error if file doesn't exist
       }
       return nil
}