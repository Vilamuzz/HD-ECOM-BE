package models

type TicketStatus struct {
	ID         int    `json:"id_status" gorm:"column:id_status;primaryKey"`
	NamaStatus string `json:"nama_status" gorm:"column:nama_status;type:varchar(50)"`

	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:StatusID;"`
}
