package models

import "time"

type Message struct {
	ID             uint64     `json:"id" gorm:"primaryKey"`
	ConversationID uint64     `json:"conversation_id" gorm:"not null;index"`
	SenderID       uint64     `json:"sender_id" gorm:"not null"`
	MessageText    string     `json:"message_text" gorm:"type:text"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	PurgeAt        *time.Time `json:"purge_at,omitempty" gorm:"index"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
