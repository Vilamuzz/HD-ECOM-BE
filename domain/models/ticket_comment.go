package models

import "time"

type TicketComment struct {
	ID            int       `json:"id_comment" gorm:"column:id_comment;primaryKey"`
	TicketID      int       `json:"id_ticket" gorm:"column:id_ticket;not null;index"`
	UserID        int       `json:"id_user" gorm:"index"`
	IsiPesan      string    `json:"isi_pesan" gorm:"column:isi_pesan"`
	TanggalDibuat time.Time `json:"tanggal_dibuat" gorm:"column:tanggal_dibuat;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket *Ticket `json:"ticket,omitempty" gorm:"foreignKey:TicketID;"`
	User   *User   `json:"user" gorm:"foreignKey:UserID"`
}
