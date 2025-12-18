package models

import "time"

type TicketAssignment struct {
	ID                int       `json:"id_assignment" gorm:"column:id_assignment;primaryKey"`
	TicketID          int       `json:"id_ticket" gorm:"column:id_ticket;index"`
	AdminID           int       `json:"id_admin" gorm:"column:id_admin"`
	PriorityID        int       `json:"id_priority" gorm:"column:id_priority;index"`
	TanggalDitugaskan time.Time `json:"tanggal_ditugaskan" gorm:"column:tanggal_ditugaskan;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket   *Ticket         `json:"ticket,omitempty" gorm:"foreignKey:TicketID"`
	Admin    *User           `json:"admin" gorm:"foreignKey:AdminID"`
	Priority *TicketPriority `json:"priority,omitempty" gorm:"foreignKey:PriorityID"`
}
