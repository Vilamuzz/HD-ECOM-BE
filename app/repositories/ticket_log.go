package repositories

import "app/domain/models"

func (r *appRepository) CreateTicketLog(log *models.TicketLog) error {
	return r.Conn.Create(log).Error
}

func (r *appRepository) GetTicketLogs() ([]models.TicketLog, error) {
	var logs []models.TicketLog
	err := r.Conn.Find(&logs).Error
	return logs, err
}

func (r *appRepository) GetTicketLogByID(id int) (*models.TicketLog, error) {
	var log models.TicketLog
	err := r.Conn.First(&log, id).Error
	return &log, err
}

func (r *appRepository) GetTicketLogsByTicketID(ticketID int) ([]models.TicketLog, error) {
	var logs []models.TicketLog
	err := r.Conn.Where("id_ticket = ?", ticketID).Find(&logs).Error
	return logs, err
}