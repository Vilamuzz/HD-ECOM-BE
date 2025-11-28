package models

type TicketCategory struct {
	ID           int    `json:"id_category" gorm:"column:id_category;primaryKey"`
	NamaCategory string `json:"nama_category" gorm:"column:nama_category;type:varchar(100)"`

	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:CategoryID;"`
}
