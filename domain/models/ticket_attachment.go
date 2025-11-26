package models

import "time"

type TicketAttachment struct {
	ID         int       `json:"id_attachment" gorm:"column:id_attachment;primaryKey;autoIncrement"`
	TicketID   int       `json:"id_ticket" gorm:"column:id_ticket;not null;index"`
	FilePath   string    `json:"file_path" gorm:"column:file_path;not null"`
	UploadedAt time.Time `json:"uploaded_at" gorm:"column:uploaded_at;default:CURRENT_TIMESTAMP"`

	// Relasi - specify references
	Ticket *Ticket `json:"ticket,omitempty" gorm:"foreignKey:TicketID;"`
}
