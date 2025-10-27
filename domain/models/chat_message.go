package models

import "time"

type ChatMessage struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	MessageText    string    `json:"message_text"`
	AttachmentURL  string    `json:"attachment_url,omitempty"`
	IsRead         bool      `json:"is_read"`
	ReadAt         time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
