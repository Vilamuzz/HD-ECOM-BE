package repositories

import "app/domain/models"

func (r *appRepository) CreateTicketComment(comment *models.TicketComment) error {
	return r.Conn.Create(comment).Error
}

func (r *appRepository) GetTicketComments() ([]models.TicketComment, error) {
	var comments []models.TicketComment
	err := r.Conn.Find(&comments).Error
	return comments, err
}

func (r *appRepository) GetTicketCommentByID(id int) (*models.TicketComment, error) {
	var comment models.TicketComment
	err := r.Conn.First(&comment, id).Error
	return &comment, err
}

func (r *appRepository) GetTicketCommentsByTicketID(ticketID int) ([]models.TicketComment, error) {
	var comments []models.TicketComment
	err := r.Conn.Where("id_ticket = ?", ticketID).Find(&comments).Error
	return comments, err
}

func (r *appRepository) UpdateTicketComment(comment *models.TicketComment) error {
	return r.Conn.Save(comment).Error
}

func (r *appRepository) DeleteTicketComment(id int) error {
	return r.Conn.Delete(&models.TicketComment{}, id).Error
}