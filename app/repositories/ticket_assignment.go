package repositories

import (
	"strconv"

	"app/domain/models"
)

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

func (r *appRepository) GetTicketAssignmentByTicketID(ticketID int) (*models.TicketAssignment, error) {
	var assignment models.TicketAssignment
	err := r.Conn.Preload("Ticket").Preload("Admin").Preload("Priority").Where("id_ticket = ?", ticketID).First(&assignment).Error // Fixed column name
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// New method for cursor-based pagination and filtering by admin ID and status
func (r *appRepository) GetTicketAssignmentsByAdminIDCursor(adminID int, limit int, cursor string, statusName string) ([]models.TicketAssignment, string, error) {
	var assignments []models.TicketAssignment

	db := r.Conn.Preload("Ticket").Preload("Admin").Preload("Priority").Where("id_admin = ?", adminID)

	// Join with tickets and ticket_statuses tables to filter by status name if provided
	if statusName != "" {
		db = db.Joins("JOIN tickets ON ticket_assignments.id_ticket = tickets.id_ticket"). // Fixed column name in JOIN
			Joins("JOIN ticket_statuses ON tickets.status_id = ticket_statuses.id_status").
			Where("ticket_statuses.nama_status = ?", statusName)
	}

	if cursor != "" {
		// Cursor is the last seen assignment ID (assuming descending order by ID)
		if lastID, err := strconv.Atoi(cursor); err == nil {
			db = db.Where("id_assignment < ?", lastID) // Fixed column name
		}
	}

	err := db.Order("id_assignment desc").Limit(limit + 1).Find(&assignments).Error // Fixed column name
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(assignments) > limit {
		nextCursor = strconv.Itoa(assignments[limit].ID)
		assignments = assignments[:limit]
	}

	return assignments, nextCursor, nil
}

// New method to count total assigned tickets by admin ID
func (r *appRepository) GetAssignedTicketCountByAdminID(adminID int) (int, error) {
	var count int64
	err := r.Conn.Model(&models.TicketAssignment{}).Where("id_admin = ?", adminID).Count(&count).Error
	return int(count), err
}

// New method to count assigned tickets by admin ID and status ID
func (r *appRepository) GetAssignedTicketCountByAdminIDAndStatus(adminID int, statusID int) (int, error) {
	var count int64
	err := r.Conn.Model(&models.TicketAssignment{}).
		Joins("JOIN tickets ON ticket_assignments.id_ticket = tickets.id_ticket").
		Where("ticket_assignments.id_admin = ? AND tickets.status_id = ?", adminID, statusID).
		Count(&count).Error
	return int(count), err
}
