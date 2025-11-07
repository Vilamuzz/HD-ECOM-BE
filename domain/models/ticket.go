package models

import "time"

type Ticket struct {
	IDTicket          int       `json:"id_ticket" gorm:"column:id_ticket;primaryKey"`
	KodeTiket         string    `json:"kode_tiket" gorm:"column:kode_tiket"`
	IDUser            int       `json:"id_user" gorm:"column:id_user"`
	Judul             string    `json:"judul" gorm:"column:judul"`
	Deskripsi         string    `json:"deskripsi" gorm:"column:deskripsi"`
	IDCategory        int       `json:"id_category" gorm:"column:id_category"`
	IDPriority        int       `json:"id_priority" gorm:"column:id_priority"`
	IDStatus          int       `json:"id_status" gorm:"column:id_status"`
	TipePengaduan     string    `json:"tipe_pengaduan" gorm:"column:tipe_pengaduan;type:varchar(50);check:tipe_pengaduan IN ('pelanggan', 'penjual')"`
	TanggalDibuat     time.Time `json:"tanggal_dibuat" gorm:"column:tanggal_dibuat;default:CURRENT_TIMESTAMP"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui" gorm:"column:tanggal_diperbarui;default:CURRENT_TIMESTAMP"`

	// Relasi
	User        *User              `json:"user" gorm:"foreignKey:IDUser"`
	Category    *TicketCategory    `json:"category" gorm:"foreignKey:IDCategory;references:IDCategory"`
	Priority    *TicketPriority    `json:"priority" gorm:"foreignKey:IDPriority"`
	Status      *TicketStatus      `json:"status" gorm:"foreignKey:IDStatus;references:IDStatus"`
	Comments    []TicketComment    `json:"comments" gorm:"foreignKey:IDTicket"`
	Attachments []TicketAttachment `json:"attachments" gorm:"foreignKey:IDTicket"`
	Assignments []TicketAssignment `json:"assignments" gorm:"foreignKey:IDTicket"`
	Logs        []TicketLog        `json:"logs" gorm:"foreignKey:IDTicket"`
}
