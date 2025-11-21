package services

import (
	"app/domain"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketConnectionWrapper struct {
	*websocket.Conn
}

func (w *WebSocketConnectionWrapper) SetReadDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebSocketConnectionWrapper) SetWriteDeadline(t time.Time) error {
	return w.Conn.SetWriteDeadline(t)
}

func NewHub() *domain.Hub {
	return &domain.Hub{
		Conversations: make(map[uint64]map[uint64]*domain.Client),
		Clients:       make(map[uint64]*domain.Client), // added
		Broadcast:     make(chan *domain.Message, 256),
		Register:      make(chan *domain.Client),
		Unregister:    make(chan *domain.Client),
	}
}

func (s *appService) Run() {
	for {
		select {
		case client := <-s.hub.Register:
			s.RegisterClient(client)

		case client := <-s.hub.Unregister:
			s.UnregisterClient(client)

		case message := <-s.hub.Broadcast:
			s.SendToRecipients(message)
		}
	}
}

func (s *appService) RegisterClient(client *domain.Client) {
	s.hub.Mu.Lock()
	defer s.hub.Mu.Unlock()

	if client.ConversationIDs == nil {
		client.ConversationIDs = make(map[uint64]bool)
	}
	// track globally
	s.hub.Clients[client.UserID] = client
}

func (s *appService) UnregisterClient(client *domain.Client) {
	s.hub.Mu.Lock()
	defer s.hub.Mu.Unlock()

	for convID := range client.ConversationIDs {
		if clients, exists := s.hub.Conversations[convID]; exists {
			delete(clients, client.UserID)
			if len(clients) == 0 {
				delete(s.hub.Conversations, convID)
			}
		}
	}
	delete(s.hub.Clients, client.UserID) // remove global reference

	close(client.Send)
}

func (s *appService) BroadcastMessage(message *domain.Message) {
	s.hub.Broadcast <- message
}

func (s *appService) SendToRecipients(message *domain.Message) {
	s.hub.Mu.RLock()
	defer s.hub.Mu.RUnlock()

	response := map[string]interface{}{
		"type":    message.Type,
		"payload": message.Data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return
	}

	// Get clients in this specific conversation
	clients, exists := s.hub.Conversations[message.ConversationID]
	if !exists || len(clients) == 0 {
		return
	}

	sentCount := 0
	for _, client := range clients {
		select {
		case client.Send <- jsonData:
			sentCount++
		default:
			go func(c *domain.Client) {
				s.hub.Unregister <- c
			}(client)
		}
	}
}

// Add helper method to join a conversation
func (s *appService) JoinConversation(client *domain.Client, conversationID uint64) {
	s.hub.Mu.Lock()
	defer s.hub.Mu.Unlock()

	if s.hub.Conversations[conversationID] == nil {
		s.hub.Conversations[conversationID] = make(map[uint64]*domain.Client)
	}

	s.hub.Conversations[conversationID][client.UserID] = client
	client.ConversationIDs[conversationID] = true
}

// helper to send direct JSON to a single client (non-conversation scoped)
func (s *appService) sendDirect(client *domain.Client, frameType string, payload interface{}) {
	resp := map[string]interface{}{
		"type":    frameType,
		"payload": payload,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	select {
	case client.Send <- data:
	case <-time.After(3 * time.Second):
	}
}
