package models

import "time"

type Message struct {
	ID             uint64    `json:"id"`
	ConversationID uint64    `json:"conversation_id"`
	SenderID       uint64    `json:"sender_id"`
	MessageText    string    `json:"message_text"`
	AttachmentURL  string    `json:"attachment_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
