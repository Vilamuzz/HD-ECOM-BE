package models

import "time"

type AdminAvailability struct {
	AdminID              uint64    `json:"admin_id"`
	CurrentConversations uint      `json:"current_conversations"`
	UpdatedAt            time.Time `json:"updated_at"`
}
