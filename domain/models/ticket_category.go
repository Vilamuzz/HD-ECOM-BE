package models

type TicketCategory struct {
	IDCategory   int    `json:"id_category" gorm:"column:id_category;primaryKey"`
	NamaCategory string `json:"nama_category" gorm:"column:nama_category;type:varchar(100)"`

	// Relasi
	Tickets []Ticket `json:"tickets" gorm:"foreignKey:IDCategory;references:IDCategory"`
}
