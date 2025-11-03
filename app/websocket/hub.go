package websocket

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

func NewHubService(hub *domain.Hub) domain.HubService {
	return &hubService{hub: hub}
}

func (h *hubService) Run() {
	log.Println("[HUB] Hub service started")
	for {
		select {
		case client := <-h.hub.Register:
			h.RegisterClient(client)

		case client := <-h.hub.Unregister:
			h.UnregisterClient(client)

		case message := <-h.hub.Broadcast:
			h.SendToRecipients(message)
		}
	}
}

func (h *hubService) RegisterClient(client *domain.Client) {
	h.hub.Mu.Lock()
	defer h.hub.Mu.Unlock()

	// Initialize ConversationIDs map if needed
	if client.ConversationIDs == nil {
		client.ConversationIDs = make(map[string]bool)
	}

	log.Printf("[HUB] User %s (Name: %s) connected. Total unique clients: %d",
		client.UserID, client.Name, h.countUniqueClients())
}

func (h *hubService) UnregisterClient(client *domain.Client) {
	h.hub.Mu.Lock()
	defer h.hub.Mu.Unlock()

	// Remove client from all conversations they're in
	for convID := range client.ConversationIDs {
		if clients, exists := h.hub.Conversations[convID]; exists {
			delete(clients, client.UserID)
			if len(clients) == 0 {
				delete(h.hub.Conversations, convID)
			}
		}
	}

	close(client.Send)
	log.Printf("[HUB] User %s disconnected. Total unique clients: %d",
		client.UserID, h.countUniqueClients())
}

func (h *hubService) BroadcastMessage(message *domain.Message) {
	h.hub.Broadcast <- message
}

func (h *hubService) SendToRecipients(message *domain.Message) {
	h.hub.Mu.RLock()
	defer h.hub.Mu.RUnlock()

	log.Printf("[HUB] Broadcasting message - ConversationID: %s, Data: %+v",
		message.ConversationID, message.Data)

	jsonData, err := json.Marshal(message.Data)
	if err != nil {
		log.Printf("[HUB] Error marshaling message: %v", err)
		return
	}

	// Get clients in this specific conversation
	clients, exists := h.hub.Conversations[message.ConversationID]
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
				h.hub.Unregister <- c
			}(client)
		}
	}

	log.Printf("[HUB] Message broadcast complete: %d/%d recipients received message",
		sentCount, len(clients))
}

func (h *hubService) countUniqueClients() int {
	uniqueUsers := make(map[string]bool)
	for _, clients := range h.hub.Conversations {
		for userID := range clients {
			uniqueUsers[userID] = true
		}
	}
	return len(uniqueUsers)
}

// Add helper method to join a conversation
func (h *hubService) JoinConversation(client *domain.Client, conversationID string) {
	h.hub.Mu.Lock()
	defer h.hub.Mu.Unlock()

	if h.hub.Conversations[conversationID] == nil {
		h.hub.Conversations[conversationID] = make(map[string]*domain.Client)
	}

	h.hub.Conversations[conversationID][client.UserID] = client
	client.ConversationIDs[conversationID] = true

	log.Printf("[HUB] User %s joined conversation %s. Conversation has %d participants",
		client.UserID, conversationID, len(h.hub.Conversations[conversationID]))
}
