package models

import "time"

type Ticket struct {
	ID                int       `json:"id_ticket" gorm:"column:id_ticket;primaryKey;autoIncrement"`
	KodeTiket         string    `json:"kode_tiket" gorm:"column:kode_tiket;unique"`
	UserID            uint64    `json:"user_id" gorm:"column:id_user;not null;index"`
	Judul             string    `json:"judul"`
	Deskripsi         string    `json:"deskripsi" gorm:"type:text"`
	CategoryID        int       `json:"category_id" gorm:"not null;index"`
	PriorityID        int       `json:"priority_id" gorm:"not null;index"`
	StatusID          int       `json:"status_id" gorm:"not null;index"`
	TipePengaduan     UserRole  `json:"tipe_pengaduan" gorm:"type:varchar(50);check:tipe_pengaduan IN ('admin', 'seller', 'customer')"`
	TanggalDibuat     time.Time `json:"tanggal_dibuat" gorm:"default:CURRENT_TIMESTAMP"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui" gorm:"default:CURRENT_TIMESTAMP"`

	// Relasi - Add references to match custom column names
	User        User               `gorm:"foreignKey:UserID"`
	Category    *TicketCategory    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Priority    *TicketPriority    `json:"priority,omitempty" gorm:"foreignKey:PriorityID"`
	Status      *TicketStatus      `json:"status,omitempty" gorm:"foreignKey:StatusID"`
	Comments    []TicketComment    `json:"comments,omitempty" gorm:"foreignKey:TicketID"`
	Attachments []TicketAttachment `json:"attachments,omitempty" gorm:"foreignKey:TicketID"`
	Assignments []TicketAssignment `json:"assignments,omitempty" gorm:"foreignKey:TicketID"`
	Logs        []TicketLog        `json:"logs,omitempty" gorm:"foreignKey:TicketID"`
}
