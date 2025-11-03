package models

type AgentAvailability struct {
	AgentID              string `json:"agent_id"`
	IsAvailable          bool   `json:"is_available"`
	MaxConversations     int    `json:"max_conversations"`
	CurrentConversations int    `json:"current_conversations"`
	UpdatedAt            int64  `json:"updated_at"`
}
