package models

import "time"

type Conversation struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	CostumerID    int64     `json:"costumer_id" gorm:"not null"`
	AgentID       int64     `json:"agent_id" gorm:"default:null"`
	Status        string    `json:"status" gorm:"default:'open'"` // open, closed, pending
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
