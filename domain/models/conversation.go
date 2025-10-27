package models

import "time"

type Conversation struct {
	ID            int64     `json:"id"`
	CostumerID    int64     `json:"costumer_id"`
	AgentID       int64     `json:"agent_id"`
	Status        string    `json:"status"`
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
