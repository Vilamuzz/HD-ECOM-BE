package websocket

import (
	"log"
	"strconv"
	"time"

	"app/domain"
	"app/domain/models"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 54 * time.Second
	maxMessageSize = 512 * 1024
)

type Client struct {
	Hub        *domain.Hub
	Conn       *websocket.Conn
	Send       chan []byte
	UserID     string
	Name       string
	Repository domain.AppRepository
}

type IncomingMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// StartClientPumps starts both read and write pumps for a client
func StartClientPumps(client *domain.Client) {
	go writePump(client)
	go readPump(client)
}

func readPump(c *domain.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg IncomingMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		handleMessage(c, &msg)
	}
}

func writePump(c *domain.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.TextMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleMessage(c *domain.Client, msg *IncomingMessage) {
	switch msg.Type {
	case "send_message":
		handleSendMessage(c, msg.Payload)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func handleSendMessage(c *domain.Client, payload map[string]interface{}) {
	var conversationID int64
	var err error

	// Check if conversation_id is provided
	if convIDStr, ok := payload["conversation_id"].(string); ok && convIDStr != "" {
		conversationID, err = strconv.ParseInt(convIDStr, 10, 64)
		if err != nil {
			log.Println("Invalid conversation_id format")
			return
		}
	} else {
		// Auto-create conversation if not provided
		userID, err := strconv.ParseInt(c.UserID, 10, 64)
		if err != nil {
			log.Printf("Invalid user ID format: %v", err)
			return
		}

		conversation := &models.Conversation{
			CostumerID:    userID,
			Status:        "open",
			LastMessageAt: time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		savedConv, err := c.Repository.CreateConversation(conversation)
		if err != nil {
			log.Printf("Error creating conversation: %v", err)
			return
		}

		conversationID = savedConv.ID
		log.Printf("Auto-created conversation %d for user %d", conversationID, userID)
	}

	text, ok := payload["text"].(string)
	if !ok || text == "" {
		log.Println("Missing text")
		return
	}

	userID, err := strconv.ParseInt(c.UserID, 10, 64)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		return
	}

	// Create chat message
	chatMessage := &models.ChatMessage{
		ConversationID: conversationID,
		SenderID:       userID,
		MessageText:    text,
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	savedMessage, err := c.Repository.SaveChatMessage(chatMessage)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	recipients, err := c.Repository.GetConversationParticipants(conversationID)
	if err != nil {
		log.Printf("Error getting conversation participants: %v", err)
		return
	}

	recipientIDs := make([]string, len(recipients))
	for i, id := range recipients {
		recipientIDs[i] = strconv.FormatInt(id, 10)
	}

	// Broadcast message
	c.Hub.Broadcast <- &domain.Message{
		ConversationID: strconv.FormatInt(conversationID, 10),
		Recipients:     recipientIDs,
		Data:           savedMessage,
	}
}
