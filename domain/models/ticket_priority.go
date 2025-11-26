package models

type TicketPriority struct {
	ID           int    `json:"id_priority" gorm:"column:id_priority;primaryKey"`
	NamaPriority string `json:"nama_priority" gorm:"column:nama_priority;type:varchar(50)"`

	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:PriorityID;"`
}
