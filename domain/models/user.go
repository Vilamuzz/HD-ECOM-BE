package models

import "time"

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      UserRole  `json:"role" gorm:"default:'customer'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Conversations     []Conversation           `json:"conversations" gorm:"foreignKey:CustomerID"`
	Messages          []Message                `json:"messages" gorm:"foreignKey:SenderID"`
	AdminStates       []AdminConversationState `json:"admin_states" gorm:"foreignKey:AdminID"`
	AdminAvailability *AdminAvailability       `json:"admin_availability" gorm:"foreignKey:AdminID"`
	Tickets           []Ticket                 `json:"tickets,omitempty" gorm:"foreignKey:UserID"`
	Comments          []TicketComment          `json:"comments" gorm:"foreignKey:UserID"`
	AssignedTickets   []TicketAssignment       `json:"assigned_tickets" gorm:"foreignKey:AdminID"`
	Logs              []TicketLog              `json:"logs" gorm:"foreignKey:UserID"`
}

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleSeller   UserRole = "seller"
	RoleCustomer UserRole = "customer"
)
