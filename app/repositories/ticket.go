package repositories

import (
	"strconv"

	"app/domain/models"
)

func (r *appRepository) CreateTicket(ticket *models.Ticket) error {
	return r.Conn.Create(ticket).Error
}

func (r *appRepository) GetTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.Conn.Preload("User").Preload("Category").Preload("Priority").Preload("Status").Find(&tickets).Error
	return tickets, err
}

func (r *appRepository) GetTicketsPaginated(limit, offset int) ([]models.Ticket, int, error) {
	var tickets []models.Ticket
	var total int64

	// Get total count
	if err := r.Conn.Model(&models.Ticket{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.Conn.Preload("User").Preload("Category").Preload("Priority").Preload("Status").
		Limit(limit).Offset(offset).Find(&tickets).Error

	return tickets, int(total), err
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

func (r *appRepository) GetTicketsByUserID(userID int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.Conn.Where("id_user = ?", userID).Preload("User").Preload("Category").Preload("Priority").Preload("Status").Find(&tickets).Error
	return tickets, err
}

// Add this function for cursor-based pagination
func (r *appRepository) GetTicketsCursor(limit int, cursor string, tipePengaduan string, statusID, priorityID, categoryID int) ([]models.Ticket, string, error) {
	var tickets []models.Ticket

	db := r.Conn.Preload("User").Preload("Category").Preload("Priority").Preload("Status")

	// Apply filters at database level
	if tipePengaduan != "" {
		db = db.Where("tipe_pengaduan = ?", tipePengaduan)
	}
	if statusID > 0 {
		db = db.Where("status_id = ?", statusID)
	}
	if priorityID > 0 {
		db = db.Where("priority_id = ?", priorityID)
	}
	if categoryID > 0 {
		db = db.Where("category_id = ?", categoryID)
	}

	if cursor != "" {
		// cursor is last seen ticket ID (assuming descending order)
		if lastID, err := strconv.Atoi(cursor); err == nil {
			db = db.Where("id_ticket < ?", lastID)
		}
	}

	err := db.Order("id_ticket desc").Limit(limit + 1).Find(&tickets).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(tickets) > limit {
		nextCursor = strconv.Itoa(tickets[limit].ID)
		tickets = tickets[:limit]
	}

	return tickets, nextCursor, nil
}
