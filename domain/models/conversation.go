package models

import "time"

type Conversation struct {
	ID            uint64             `json:"id" gorm:"primaryKey"`
	CustomerID    uint64             `json:"customer_id" gorm:"not null"`
	AdminID       uint8              `json:"admin_id" gorm:"not null"`
	Status        ConversationStatus `json:"status" gorm:"default:'open'"`
	LastMessageAt time.Time          `json:"last_message_at"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type ConversationStatus string

const (
	StatusOpen   ConversationStatus = "open"
	StatusClosed ConversationStatus = "closed"
)
