package models

import "time"

type AdminConversationState struct {
	AdminID        uint64     `json:"admin_id" gorm:"primaryKey;autoIncrement:false"`
	ConversationID uint64    `json:"conversation_id" gorm:"primaryKey;autoIncrement:false"`
	UnreadCount    uint      `json:"unread_count"`
	LastMessageID  uint64    `json:"last_message_id"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
