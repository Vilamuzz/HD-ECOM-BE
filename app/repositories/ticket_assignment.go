package repositories

import "app/domain/models"

func (r *appRepository) CreateTicketAssignment(assignment *models.TicketAssignment) error {
	return r.Conn.Create(assignment).Error
}

func (r *appRepository) GetTicketAssignments() ([]models.TicketAssignment, error) {
	var assignments []models.TicketAssignment
	err := r.Conn.Preload("Ticket").Preload("Admin").Preload("Priority").Find(&assignments).Error
	return assignments, err
}

func (r *appRepository) GetTicketAssignmentByID(id int) (*models.TicketAssignment, error) {
	var assignment models.TicketAssignment
	err := r.Conn.Preload("Ticket").Preload("Admin").Preload("Priority").First(&assignment, id).Error
	return &assignment, err
}

func (r *appRepository) UpdateTicketAssignment(assignment *models.TicketAssignment) error {
	return r.Conn.Save(assignment).Error
}

func (r *appRepository) DeleteTicketAssignment(id int) error {
	return r.Conn.Delete(&models.TicketAssignment{}, id).Error
}

func (r *appRepository) UpdateTicketPriorityByID(ticketID int, priorityID int) error {
	return r.Conn.Model(&models.Ticket{}).Where("id_ticket = ?", ticketID).Update("priority_id", priorityID).Error
}
