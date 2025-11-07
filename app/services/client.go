package services

import (
	"encoding/json"
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

type IncomingMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// StartClientPumps starts both read and write pumps for a client
func (s *appService) StartClientPumps(client *domain.Client) {
	go writePump(client)
	go s.readPump(client)
}

func (s *appService) readPump(c *domain.Client) {
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
		c.Conn.ReadJSON(&msg)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		s.handleMessage(c, &msg)
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

func (s *appService) handleMessage(c *domain.Client, msg *IncomingMessage) {
	if msg.Type != "send_message" {
		sendErrorToClient(c, "Unknown message type "+msg.Type)
		return
	}
	s.handleSendMessage(c, msg.Payload)
}

// Add error response helper
func sendErrorToClient(c *domain.Client, errorMsg string) {
	errorResponse := map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	}
	jsonData, _ := json.Marshal(errorResponse)
	select {
	case c.Send <- jsonData:
	case <-time.After(5 * time.Second):
		log.Println("Timeout sending message")
	}
}

func (s *appService) handleSendMessage(c *domain.Client, payload map[string]interface{}) {
	var conversationID uint64
	var err error

	_, ok := payload["conversation_id"]

	userID, err := strconv.ParseUint(c.UserID, 10, 64)
	if err != nil {
		sendErrorToClient(c, "Invalid user ID")
		return
	}

	// Find or create conversation
	if !ok || conversationID == 0 {
		existingConv, err := c.Repository.FindActiveConversationForCustomer(userID)
		if err == nil && existingConv != nil {
			conversationID = existingConv.ID
		} else {
			admin, err := c.Repository.GetAdminWithLeastConversations()
			var adminID uint64
			if err == nil && admin != nil {
				adminID = admin.ID
			}

			conversation := &models.Conversation{
				UserID:        userID,
				AdminID:       adminID,
				LastMessageAt: time.Now(),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			savedConv, err := c.Repository.CreateConversation(conversation)
			if err != nil {
				sendErrorToClient(c, "Failed to create conversation")
				return
			}

			conversationID = savedConv.ID

			// If admin was assigned, join them to the conversation
			if adminID != 0 {
				s.notifyAdminOfNewConversation(adminID, conversationID)
			}
		}

		convIDStr := strconv.FormatUint(conversationID, 10)
		s.JoinConversation(c, convIDStr)

		// Send conversation_id back to client so they can use it in future messages
		convResponse := map[string]interface{}{
			"type":            "conversation_created",
			"conversation_id": conversationID,
		}
		jsonData, _ := json.Marshal(convResponse)
		select {
		case c.Send <- jsonData:
		case <-time.After(5 * time.Second):
			log.Println("Timeout sending message")
		}
	} else {
		convIDStr := strconv.FormatUint(conversationID, 10)
		if !c.ConversationIDs[convIDStr] {
			s.JoinConversation(c, convIDStr)
		}
	}

	// Extract message text
	text, ok := payload["text"].(string)
	if !ok || text == "" {
		sendErrorToClient(c, "Message text is required")
		return
	}

	// Create chat message
	chatMessage := &models.ChatMessage{
		ConversationID: conversationID,
		SenderID:       userID,
		MessageText:    text,
		CreatedAt:      time.Now(),
	}

	savedMessage, err := c.Repository.SaveChatMessage(chatMessage)
	if err != nil {
		sendErrorToClient(c, "Failed to save message")
		return
	}

	// Update conversation's last_message_at
	err = c.Repository.UpdateConversationLastMessage(conversationID)
	if err != nil {
		sendErrorToClient(c, "Failed to update conversation")
	}

	// Get conversation participants
	recipients, err := c.Repository.GetConversationParticipants(conversationID)
	if err != nil {
		sendErrorToClient(c, "Failed to get conversation participants")
		return
	}

	recipientIDs := make([]string, len(recipients))
	for i, id := range recipients {
		recipientIDs[i] = strconv.FormatUint(id, 10)
	}

	// Broadcast to conversation participants only
	c.Hub.Broadcast <- &domain.Message{
		ConversationID: strconv.FormatUint(conversationID, 10),
		Data:           savedMessage,
	}
}

// Add helper function to notify admin
func (s *appService) notifyAdminOfNewConversation(adminID uint64, conversationID uint64) {
	s.hub.Mu.RLock()
	defer s.hub.Mu.RUnlock()

	adminIDStr := strconv.FormatUint(adminID, 10)
	convIDStr := strconv.FormatUint(conversationID, 10)

	// Find admin's client in any conversation
	for _, clients := range s.hub.Conversations {
		if adminClient, exists := clients[adminIDStr]; exists {
			// Join admin to new conversation
			s.JoinConversation(adminClient, convIDStr)

			// Notify admin
			notification := map[string]interface{}{
				"type":            "new_conversation_assigned",
				"conversation_id": conversationID,
				"message":         "A new conversation has been assigned to you",
			}
			jsonData, _ := json.Marshal(notification)
			select {
			case adminClient.Send <- jsonData:
			default:
			}
			break
		}
	}
}
