package models

import "time"

type Conversation struct {
	ID            uint64             `json:"id" gorm:"primaryKey;"`
	CustomerID    uint64             `json:"customer_id" gorm:"not null;index"`
	AdminID       uint64             `json:"admin_id" gorm:"not null;index"`
	Status        ConversationStatus `json:"status" gorm:"default:'open'"`
	LastMessageAt time.Time          `json:"last_message_at"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`

	Messages []Message `json:"messages" gorm:"foreignKey:ConversationID"`
	Customer User      `json:"customer" gorm:"foreignKey:CustomerID"`
	Admin    User      `json:"admin" gorm:"foreignKey:AdminID"`
}

type ConversationStatus string

const (
	StatusOpen   ConversationStatus = "open"
	StatusClosed ConversationStatus = "closed"
)
