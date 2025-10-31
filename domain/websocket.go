package domain

import (
	"sync"
	"time"
)

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.RWMutex
}

type Message struct {
	ConversationID string      `json:"conversation_id"`
	Recipients     []string    `json:"-"`
	Data           interface{} `json:"data"`
}

type Client struct {
	Hub        *Hub
	Conn       WebSocketConnection
	Send       chan []byte
	UserID     string
	Name       string
	Repository AppRepository
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

// HubService interface for hub operations
type HubService interface {
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *Message)
	SendToRecipients(message *Message)
}
