package models

import "time"

type TicketLog struct {
	IDLog     int       `json:"id_log" gorm:"column:id_log;primaryKey"`
	IDTicket  int       `json:"id_ticket" gorm:"column:id_ticket"`
	Aktivitas string    `json:"aktivitas" gorm:"column:aktivitas"`
	IDUser    int       `json:"id_user" gorm:"column:id_user"`
	Waktu     time.Time `json:"waktu" gorm:"column:waktu;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket *Ticket `json:"ticket" gorm:"foreignKey:IDTicket"`
	User   *User   `json:"user" gorm:"foreignKey:IDUser"`
}
