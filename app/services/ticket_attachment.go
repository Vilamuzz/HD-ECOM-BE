package services

import "app/domain/models"

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
	return s.repo.DeleteTicketAttachment(id)
}