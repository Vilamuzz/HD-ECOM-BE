package repositories

import (
	"app/domain/models"
)

func (r *appRepository) CreateTicketAttachment(attachment *models.TicketAttachment) error {
	return r.Conn.Create(attachment).Error
}

func (r *appRepository) GetTicketAttachments() ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	err := r.Conn.Find(&attachments).Error
	return attachments, err
}

func (r *appRepository) GetTicketAttachmentByID(id int) (*models.TicketAttachment, error) {
	var attachment models.TicketAttachment
	err := r.Conn.First(&attachment, id).Error
	return &attachment, err
}

func (r *appRepository) GetTicketAttachmentsByTicketID(ticketID int) ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	err := r.Conn.Where("id_ticket = ?", ticketID).Find(&attachments).Error
	return attachments, err
}

func (r *appRepository) UpdateTicketAttachment(attachment *models.TicketAttachment) error {
	return r.Conn.Save(attachment).Error
}

func (r *appRepository) DeleteTicketAttachment(id int) error {
	return r.Conn.Delete(&models.TicketAttachment{}, id).Error
}