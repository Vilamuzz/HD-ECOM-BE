package models

type TicketPriority struct {
	IDPriority   int    `json:"id_priority" gorm:"column:id_priority;primaryKey"`
	NamaPriority string `json:"nama_priority" gorm:"column:nama_priority;type:varchar(50)"`

	// Relasi
	Tickets []Ticket `json:"tickets" gorm:"foreignKey:IDPriority"`
}
