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
		Clients:    make(map[string]*domain.Client),
		Broadcast:  make(chan *domain.Message, 256),
		Register:   make(chan *domain.Client),
		Unregister: make(chan *domain.Client),
	}
}

func NewHubService(hub *domain.Hub) domain.HubService {
	return &hubService{hub: hub}
}

func (h *hubService) Run() {
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
	h.hub.Clients[client.UserID] = client
	h.hub.Mu.Unlock()
	log.Printf("User %s connected. Total: %d", client.UserID, len(h.hub.Clients))
}

func (h *hubService) UnregisterClient(client *domain.Client) {
	h.hub.Mu.Lock()
	if _, ok := h.hub.Clients[client.UserID]; ok {
		delete(h.hub.Clients, client.UserID)
		close(client.Send)
	}
	h.hub.Mu.Unlock()
	log.Printf("User %s disconnected. Total: %d", client.UserID, len(h.hub.Clients))
}

func (h *hubService) BroadcastMessage(message *domain.Message) {
	h.hub.Broadcast <- message
}

func (h *hubService) SendToRecipients(message *domain.Message) {
	h.hub.Mu.RLock()
	defer h.hub.Mu.RUnlock()

	jsonData, err := json.Marshal(message.Data)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for _, recipientID := range message.Recipients {
		if client, ok := h.hub.Clients[recipientID]; ok {
			select {
			case client.Send <- jsonData:
				log.Printf("Message sent to user %s", recipientID)
			default:
				log.Printf("Channel full for user %s, closing connection", recipientID)
				close(client.Send)
				delete(h.hub.Clients, recipientID)
			}
		}
	}
}
