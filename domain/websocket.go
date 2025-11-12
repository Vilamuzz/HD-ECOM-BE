package domain

import (
	"sync"
	"time"
)

type Hub struct {
	Conversations map[uint64]map[uint64]*Client
	Broadcast     chan *Message
	Register      chan *Client
	Unregister    chan *Client
	Mu            sync.RWMutex
}

type Message struct {
	ConversationID uint64      `json:"conversation_id"`
	Recipients     []string    `json:"-"`
	Type           string      `json:"type"`
	Data           interface{} `json:"data"`
}

type Client struct {
	Hub             *Hub
	Conn            WebSocketConnection
	Send            chan []byte
	UserID          uint64
	Name            string
	Repository      AppRepository
	ConversationIDs map[uint64]bool
}

type WebSocketConnection interface {
	ReadJSON(v interface{}) error
	WriteMessage(messageType int, data []byte) error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	Close() error
}
