package models

type TicketStatus struct {
	IDStatus   int    `json:"id_status" gorm:"column:id_status;primaryKey"`
	NamaStatus string `json:"nama_status" gorm:"column:nama_status;type:varchar(50)"`

	// Relasi
	Tickets []Ticket `json:"tickets" gorm:"foreignKey:IDStatus;references:IDStatus"`
}
