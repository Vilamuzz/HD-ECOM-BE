package repositories

import "app/domain/models"

func (r *appRepository) CreateTicket(ticket *models.Ticket) error {
	return r.Conn.Create(ticket).Error
}

func (r *appRepository) GetTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.Conn.Preload("User").Preload("Category").Preload("Priority").Preload("Status").Find(&tickets).Error
	return tickets, err
}

func (r *appRepository) GetTicketByID(id int) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.Conn.Preload("User").Preload("Category").Preload("Priority").Preload("Status").First(&ticket, id).Error
	return &ticket, err
}

func (r *appRepository) UpdateTicket(ticket *models.Ticket) error {
	return r.Conn.Save(ticket).Error
}

func (r *appRepository) DeleteTicket(id int) error {
	return r.Conn.Delete(&models.Ticket{}, id).Error
}
