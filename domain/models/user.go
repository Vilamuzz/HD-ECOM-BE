package models

import "time"

type User struct {
	IDUser    int       `json:"id_user" gorm:"column:id_user;primaryKey"`
	ID        int       `json:"id" gorm:"column:id"`
	Nama      string    `json:"nama" gorm:"column:nama"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex:idx_users_email"`
	Username  string    `json:"username" gorm:"column:username;uniqueIndex:idx_users_username"`
	Role      string    `json:"role" gorm:"column:role;type:varchar(50);check:role IN ('pelanggan', 'penjual', 'admin')"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`

	// Relasi
	Tickets         []Ticket           `json:"tickets" gorm:"foreignKey:IDUser"`
	Comments        []TicketComment    `json:"comments" gorm:"foreignKey:IDUser"`
	AssignedTickets []TicketAssignment `json:"assigned_tickets" gorm:"foreignKey:IDAdmin"`
	Logs            []TicketLog        `json:"logs" gorm:"foreignKey:IDUser"`
}
