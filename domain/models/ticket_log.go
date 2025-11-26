package models

import "time"

type TicketLog struct {
	ID     int       `json:"id_log" gorm:"column:id_log;primaryKey"`
	TicketID  int       `json:"id_ticket" gorm:"column:id_ticket;index"`
	Aktivitas string    `json:"aktivitas" gorm:"column:aktivitas"`
	UserID    int       `json:"id_user" gorm:"column:id_user"`
	Waktu     time.Time `json:"waktu" gorm:"column:waktu;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket *Ticket `json:"ticket,omitempty" gorm:"foreignKey:TicketID;"`
	User   *User   `json:"user" gorm:"foreignKey:UserID"`
}
