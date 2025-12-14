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
	return r.Conn.Model(&models.Ticket{}).Where("id_ticket = ?", ticket.ID).Updates(map[string]interface{}{
		"kode_tiket":         ticket.KodeTiket,
		"id_user":            ticket.UserID,
		"judul":              ticket.Judul,
		"deskripsi":          ticket.Deskripsi,
		"category_id":        ticket.CategoryID,
		"priority_id":        ticket.PriorityID,
		"status_id":          ticket.StatusID,
		"tipe_pengaduan":     ticket.TipePengaduan,
		"tanggal_diperbarui": ticket.TanggalDiperbarui,
	}).Error
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

func (r *appRepository) GetOpenTicketCountsByType() (customerCount int, sellerCount int, err error) {
	var customerResult int64
	var sellerResult int64

	// Count tickets with tipe_pengaduan = customer and status_id = 1
	if err := r.Conn.Model(&models.Ticket{}).
		Where("tipe_pengaduan = ? AND status_id = ?", models.RoleCustomer, 1).
		Count(&customerResult).Error; err != nil {
		return 0, 0, err
	}

	// Count tickets with tipe_pengaduan = seller and status_id = 1
	if err := r.Conn.Model(&models.Ticket{}).
		Where("tipe_pengaduan = ? AND status_id = ?", models.RoleSeller, 1).
		Count(&sellerResult).Error; err != nil {
		return 0, 0, err
	}

	return int(customerResult), int(sellerResult), nil
}

func (r *appRepository) GetTicketStatistics() (total int, inProgress int, resolved int, priorityCounts map[int]int, err error) {
	var totalCount int64
	var inProgressCount int64
	var resolvedCount int64

	// Get total ticket count
	if err := r.Conn.Model(&models.Ticket{}).Count(&totalCount).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	// Get in progress tickets (status_id = 2)
	if err := r.Conn.Model(&models.Ticket{}).
		Where("status_id = ?", 2).
		Count(&inProgressCount).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	// Get resolved tickets (status_id = 3)
	if err := r.Conn.Model(&models.Ticket{}).
		Where("status_id = ?", 3).
		Count(&resolvedCount).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	// Get counts by priority (1, 2, 3, 4)
	priorityCounts = make(map[int]int)
	for priorityID := 1; priorityID <= 4; priorityID++ {
		var count int64
		if err := r.Conn.Model(&models.Ticket{}).
			Where("priority_id = ?", priorityID).
			Count(&count).Error; err != nil {
			return 0, 0, 0, nil, err
		}
		priorityCounts[priorityID] = int(count)
	}

	return int(totalCount), int(inProgressCount), int(resolvedCount), priorityCounts, nil
}
