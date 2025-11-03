package services

import (
	"app/domain"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketConnectionWrapper wraps gorilla websocket to implement the interface
type WebSocketConnectionWrapper struct {
	*websocket.Conn
}

func (w *WebSocketConnectionWrapper) SetReadDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebSocketConnectionWrapper) SetWriteDeadline(t time.Time) error {
	return w.Conn.SetWriteDeadline(t)
}

// hubService implements the HubService interface
type hubService struct {
	hub *domain.Hub
}

func NewHub() *domain.Hub {
	return &domain.Hub{
		Conversations: make(map[string]map[string]*domain.Client),
		Broadcast:     make(chan *domain.Message, 256),
		Register:      make(chan *domain.Client),
		Unregister:    make(chan *domain.Client),
	}
}

func (s *appService) Run() {
	log.Println("[HUB] Hub service started")
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

	// Initialize ConversationIDs map if needed
	if client.ConversationIDs == nil {
		client.ConversationIDs = make(map[string]bool)
	}

	log.Printf("[HUB] User %s (Name: %s) connected. Total unique clients: %d",
		client.UserID, client.Name, s.countUniqueClients())
}

func (s *appService) UnregisterClient(client *domain.Client) {
	s.hub.Mu.Lock()
	defer s.hub.Mu.Unlock()

	// Remove client from all conversations they're in
	for convID := range client.ConversationIDs {
		if clients, exists := s.hub.Conversations[convID]; exists {
			delete(clients, client.UserID)
			if len(clients) == 0 {
				delete(s.hub.Conversations, convID)
			}
		}
	}

	close(client.Send)
	log.Printf("[HUB] User %s disconnected. Total unique clients: %d",
		client.UserID, s.countUniqueClients())
}

func (s *appService) BroadcastMessage(message *domain.Message) {
	s.hub.Broadcast <- message
}

func (s *appService) SendToRecipients(message *domain.Message) {
	s.hub.Mu.RLock()
	defer s.hub.Mu.RUnlock()

	log.Printf("[HUB] Broadcasting message - ConversationID: %s, Data: %+v",
		message.ConversationID, message.Data)

	jsonData, err := json.Marshal(message.Data)
	if err != nil {
		log.Printf("[HUB] Error marshaling message: %v", err)
		return
	}

	// Get clients in this specific conversation
	clients, exists := s.hub.Conversations[message.ConversationID]
	if !exists || len(clients) == 0 {
		log.Printf("[HUB] No clients connected to conversation %s", message.ConversationID)
		return
	}

	sentCount := 0
	for userID, client := range clients {
		select {
		case client.Send <- jsonData:
			sentCount++
			log.Printf("[HUB] Message sent to user %s in conversation %s", userID, message.ConversationID)
		default:
			log.Printf("[HUB] Channel full for user %s, closing connection", userID)
			go func(c *domain.Client) {
				s.hub.Unregister <- c
			}(client)
		}
	}

	log.Printf("[HUB] Message broadcast complete: %d/%d recipients received message",
		sentCount, len(clients))
}

func (s *appService) countUniqueClients() int {
	uniqueUsers := make(map[string]bool)
	for _, clients := range s.hub.Conversations {
		for userID := range clients {
			uniqueUsers[userID] = true
		}
	}
	return len(uniqueUsers)
}

// Add helper method to join a conversation
func (s *appService) JoinConversation(client *domain.Client, conversationID string) {
	s.hub.Mu.Lock()
	defer s.hub.Mu.Unlock()

	if s.hub.Conversations[conversationID] == nil {
		s.hub.Conversations[conversationID] = make(map[string]*domain.Client)
	}

	s.hub.Conversations[conversationID][client.UserID] = client
	client.ConversationIDs[conversationID] = true

	log.Printf("[HUB] User %s joined conversation %s. Conversation has %d participants",
		client.UserID, conversationID, len(s.hub.Conversations[conversationID]))
}
