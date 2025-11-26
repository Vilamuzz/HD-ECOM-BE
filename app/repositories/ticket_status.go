package repositories

import "app/domain/models"

func (r *appRepository) CreateTicketStatus(status *models.TicketStatus) error {
	return r.Conn.Create(status).Error
}

func (r *appRepository) GetTicketStatuses() ([]models.TicketStatus, error) {
	var statuses []models.TicketStatus
	err := r.Conn.Find(&statuses).Error
	return statuses, err
}

func (r *appRepository) GetTicketStatusByID(id int) (*models.TicketStatus, error) {
	var status models.TicketStatus
	err := r.Conn.First(&status, id).Error
	return &status, err
}

func (r *appRepository) UpdateTicketStatus(status *models.TicketStatus) error {
	return r.Conn.Save(status).Error
}

func (r *appRepository) DeleteTicketStatus(id int) error {
	return r.Conn.Delete(&models.TicketStatus{}, id).Error
}
