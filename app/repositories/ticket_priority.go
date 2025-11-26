package repositories

import "app/domain/models"

func (r *appRepository) CreateTicketPriority(priority *models.TicketPriority) error {
	return r.Conn.Create(priority).Error
}

func (r *appRepository) GetTicketPriorities() ([]models.TicketPriority, error) {
	var priorities []models.TicketPriority
	err := r.Conn.Find(&priorities).Error
	return priorities, err
}

func (r *appRepository) GetTicketPriorityByID(id int) (*models.TicketPriority, error) {
	var priority models.TicketPriority
	err := r.Conn.First(&priority, id).Error
	return &priority, err
}

func (r *appRepository) UpdateTicketPriority(priority *models.TicketPriority) error {
	return r.Conn.Save(priority).Error
}

func (r *appRepository) DeleteTicketPriority(id int) error {
	return r.Conn.Delete(&models.TicketPriority{}, id).Error
}
