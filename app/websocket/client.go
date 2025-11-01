package websocket

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
		log.Printf("[WS] ReadPump closing for user %s", c.UserID)
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		log.Printf("[WS] Pong received from user %s", c.UserID)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg IncomingMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] WebSocket unexpected close error for user %s: %v", c.UserID, err)
			} else {
				log.Printf("[WS] WebSocket read error for user %s: %v", c.UserID, err)
			}
			break
		}

		log.Printf("[WS] Received message from user %s: Type=%s, Payload=%+v", c.UserID, msg.Type, msg.Payload)

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
				log.Printf("[WS] Error writing message to user %s: %v", c.UserID, err)
				return
			}

			log.Printf("[WS] Message written to user %s", c.UserID)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleMessage(c *domain.Client, msg *IncomingMessage) {
	log.Printf("[WS] Handling message type: %s for user %s", msg.Type, c.UserID)

	switch msg.Type {
	case "send_message":
		handleSendMessage(c, msg.Payload)
	default:
		log.Printf("[WS] Unknown message type: %s from user %s", msg.Type, c.UserID)
	}
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
		log.Printf("[WS] Error sent to client %s: %s", c.UserID, errorMsg)
	default:
		log.Printf("[WS] Could not send error to client %s", c.UserID)
	}
}

func handleSendMessage(c *domain.Client, payload map[string]interface{}) {
	log.Printf("[WS] handleSendMessage - User: %s, Payload: %+v", c.UserID, payload)

	var conversationID int64
	var err error

	// Check if conversation_id is provided
	convIDRaw, hasConvID := payload["conversation_id"]
	log.Printf("[WS] Checking conversation_id - hasConvID: %v, raw value: %+v (type: %T)", hasConvID, convIDRaw, convIDRaw)

	if hasConvID {
		// Try to parse conversation_id from different types
		switch v := convIDRaw.(type) {
		case string:
			log.Printf("[WS] conversation_id is string: '%s'", v)
			if v != "" {
				conversationID, err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Printf("[WS] Error parsing string conversation_id '%s': %v", v, err)
					sendErrorToClient(c, "Invalid conversation_id format")
					return
				}
				log.Printf("[WS] Parsed conversation_id from string: %d", conversationID)
			} else {
				log.Println("[WS] conversation_id is empty string")
				hasConvID = false
			}
		case float64:
			conversationID = int64(v)
			log.Printf("[WS] conversation_id is number: %d", conversationID)
		case int64:
			conversationID = v
			log.Printf("[WS] conversation_id is int64: %d", conversationID)
		case int:
			conversationID = int64(v)
			log.Printf("[WS] conversation_id is int: %d", conversationID)
		default:
			log.Printf("[WS] conversation_id has unexpected type: %T, value: %+v", v, v)
			hasConvID = false
		}
	}

	userID, err := strconv.ParseInt(c.UserID, 10, 64)
	if err != nil {
		log.Printf("[WS] Error parsing user ID '%s': %v", c.UserID, err)
		sendErrorToClient(c, "Invalid user ID")
		return
	}

	// Find or create conversation
	if !hasConvID || conversationID == 0 {
		log.Printf("[WS] No conversation_id provided, finding or creating active conversation for user %d", userID)

		// Try to find existing open conversation for this customer
		existingConv, err := c.Repository.FindActiveConversationForCustomer(userID)
		if err == nil && existingConv != nil {
			conversationID = existingConv.ID
			log.Printf("[WS] Found existing active conversation %d for user %d", conversationID, userID)
		} else {
			// Create new conversation
			log.Printf("[WS] Creating new conversation for user %d", userID)
			conversation := &models.Conversation{
				CostumerID:    userID,
				Status:        "open",
				LastMessageAt: time.Now(),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			savedConv, err := c.Repository.CreateConversation(conversation)
			if err != nil {
				log.Printf("[WS] Error creating conversation: %v", err)
				sendErrorToClient(c, "Failed to create conversation")
				return
			}

			conversationID = savedConv.ID
			log.Printf("[WS] Created new conversation %d for user %d", conversationID, userID)
		}

		// Join the newly created conversation
		convIDStr := strconv.FormatInt(conversationID, 10)
		hubSvc := NewHubService(c.Hub)
		hubSvc.JoinConversation(c, convIDStr)

		// Send conversation_id back to client so they can use it in future messages
		convResponse := map[string]interface{}{
			"type":            "conversation_created",
			"conversation_id": conversationID,
		}
		jsonData, _ := json.Marshal(convResponse)
		select {
		case c.Send <- jsonData:
			log.Printf("[WS] Sent conversation_id %d to client %s", conversationID, c.UserID)
		default:
			log.Printf("[WS] Could not send conversation_id to client %s", c.UserID)
		}
	} else {
		// Ensure user is in the conversation
		convIDStr := strconv.FormatInt(conversationID, 10)
		if !c.ConversationIDs[convIDStr] {
			hubSvc := NewHubService(c.Hub)
			hubSvc.JoinConversation(c, convIDStr)
		}
	}

	// Extract message text
	text, ok := payload["text"].(string)
	if !ok || text == "" {
		log.Printf("[WS] Missing or invalid text in payload: %+v", payload)
		sendErrorToClient(c, "Message text is required")
		return
	}
	log.Printf("[WS] Message text: '%s'", text)

	// Create chat message
	chatMessage := &models.ChatMessage{
		ConversationID: conversationID,
		SenderID:       userID,
		MessageText:    text,
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	log.Printf("[WS] Saving chat message: ConversationID=%d, SenderID=%d, Text='%s'",
		chatMessage.ConversationID, chatMessage.SenderID, chatMessage.MessageText)

	savedMessage, err := c.Repository.SaveChatMessage(chatMessage)
	if err != nil {
		log.Printf("[WS] Error saving message: %v", err)
		sendErrorToClient(c, "Failed to save message")
		return
	}

	log.Printf("[WS] Message saved with ID: %d", savedMessage.ID)

	// Update conversation's last_message_at
	err = c.Repository.UpdateConversationLastMessage(conversationID)
	if err != nil {
		log.Printf("[WS] Warning: Failed to update conversation last_message_at: %v", err)
	}

	// Get conversation participants
	recipients, err := c.Repository.GetConversationParticipants(conversationID)
	if err != nil {
		log.Printf("[WS] Error getting conversation participants for conversation %d: %v", conversationID, err)
		sendErrorToClient(c, "Failed to get conversation participants")
		return
	}

	log.Printf("[WS] Conversation %d participants: %+v", conversationID, recipients)

	recipientIDs := make([]string, len(recipients))
	for i, id := range recipients {
		recipientIDs[i] = strconv.FormatInt(id, 10)
	}

	log.Printf("[WS] Broadcasting message to recipients: %+v", recipientIDs)

	// Broadcast to conversation participants only
	c.Hub.Broadcast <- &domain.Message{
		ConversationID: strconv.FormatInt(conversationID, 10),
		Data:           savedMessage,
	}

	log.Printf("[WS] Message broadcast complete for conversation %d", conversationID)
}
