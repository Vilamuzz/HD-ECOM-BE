package models

import "time"

type TicketAssignment struct {
	IDAssignment      int       `json:"id_assignment" gorm:"column:id_assignment;primaryKey"`
	IDTicket          int       `json:"id_ticket" gorm:"column:id_ticket"`
	IDAdmin           int       `json:"id_admin" gorm:"column:id_admin"`
	TanggalDitugaskan time.Time `json:"tanggal_ditugaskan" gorm:"column:tanggal_ditugaskan;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket *Ticket `json:"ticket" gorm:"foreignKey:IDTicket"`
	Admin  *User   `json:"admin" gorm:"foreignKey:IDAdmin"`
}
