package models

import "time"

type TicketAttachment struct {
	IDAttachment int       `json:"id_attachment" gorm:"column:id_attachment;primaryKey"`
	IDTicket     int       `json:"id_ticket" gorm:"column:id_ticket"`
	FilePath     string    `json:"file_path" gorm:"column:file_path"`
	UploadedAt   time.Time `json:"uploaded_at" gorm:"column:uploaded_at;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket *Ticket `json:"ticket" gorm:"foreignKey:IDTicket"`
}
