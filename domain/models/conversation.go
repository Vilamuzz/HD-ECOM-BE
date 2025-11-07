package models

import "time"

type Conversation struct {
	ID            uint64     `json:"id" gorm:"primaryKey"`
	UserID        uint64     `json:"user_id" gorm:"not null"`
	AdminID       uint64     `json:"admin_id" gorm:"default:null"`
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
