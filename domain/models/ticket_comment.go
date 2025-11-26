package models

import "time"

type TicketComment struct {
	IDComment     int       `json:"id_comment" gorm:"column:id_comment;primaryKey"`
	IDTicket      int       `json:"id_ticket" gorm:"column:id_ticket"`
	IDUser        int       `json:"id_user" gorm:"column:id_user"`
	IsiPesan      string    `json:"isi_pesan" gorm:"column:isi_pesan"`
	TanggalDibuat time.Time `json:"tanggal_dibuat" gorm:"column:tanggal_dibuat;default:CURRENT_TIMESTAMP"`

	// Relasi
	Ticket *Ticket `json:"ticket" gorm:"foreignKey:IDTicket"`
	User   *User   `json:"user" gorm:"foreignKey:IDUser"`
}
