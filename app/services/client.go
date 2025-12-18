package services

import (
	"context"
	"encoding/json"
	"fmt"
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
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			// Connection closed or error - exit gracefully
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %d: %v", c.UserID, err)
			}
			return
		}
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
	switch msg.Type {
	case "send_message":
		s.handleSendMessage(c, msg.Payload)
	case "subscribe":
		s.handleSubscribe(c, msg.Payload)
	case "unsubscribe":
		s.handleUnsubscribe(c, msg.Payload)
	default:
		sendErrorToClient(c, "Unknown message type "+msg.Type)
	}
}

func (s *appService) handleSubscribe(c *domain.Client, payload map[string]interface{}) {
	convID, err := parseConversationID(payload["conversation_id"])
	log.Println("Subscribing to conversation ID:", convID)
	if err != nil {
		sendErrorToClient(c, "Invalid conversation ID")
		return
	}

	// Authorization check - verify user can subscribe to this conversation
	conversation, err := c.Repository.GetConversationByID(convID)
	if err != nil {
		sendErrorToClient(c, "Conversation not found")
		return
	}

	user, err := c.Repository.GetUserByID(c.UserID)
	if err != nil {
		sendErrorToClient(c, "Failed to get user")
		return
	}

	// Check if user has permission (admin or conversation participant)
	isAdmin := user.Role == models.RoleAdmin
	isParticipant := conversation.CustomerID == c.UserID || conversation.AdminID == c.UserID

	if !isAdmin && !isParticipant {
		sendErrorToClient(c, "Access denied - not your conversation")
		return
	}

	if c.ConversationIDs == nil {
		c.ConversationIDs = make(map[uint64]bool)
	}

	s.JoinConversation(c, convID)

	// Reset unread count if subscriber is admin
	if user.Role == models.RoleAdmin {

		if state, sErr := c.Repository.GetAdminConversationState(user.ID, convID); sErr == nil && state != nil {
			// Get latest message ID (if any) to set LastMessageID
			msgs, _, mErr := c.Repository.GetMessageHistoryForAdmin(convID, 1, "")
			var lastID uint64
			if mErr == nil && len(msgs) > 0 {
				lastID = msgs[0].ID
			}
			if rErr := c.Repository.ResetState(state, lastID); rErr != nil {
				log.Printf("Failed to reset state for admin %d conv %d: %v", uint8(user.ID), convID, rErr)
			}
		}
	}

	resp := map[string]interface{}{
		"type":    "subscribed",
		"payload": map[string]uint64{"conversation_id": convID},
	}
	jsonData, _ := json.Marshal(resp)
	select {
	case c.Send <- jsonData:
	case <-time.After(5 * time.Second):
		log.Println("Timeout sending subscribe response")
	}
}

func (s *appService) handleUnsubscribe(c *domain.Client, payload map[string]interface{}) {
	convID, err := parseConversationID(payload["conversation_id"])
	if err != nil {
		sendErrorToClient(c, "Invalid conversation ID")
		return
	}

	// Remove client from conversation in hub
	s.hub.Mu.Lock()
	if clients, exists := s.hub.Conversations[convID]; exists {
		delete(clients, c.UserID)
		if len(clients) == 0 {
			delete(s.hub.Conversations, convID)
		}
	}
	delete(c.ConversationIDs, convID)
	s.hub.Mu.Unlock()

	// Acknowledge unsubscription
	resp := map[string]interface{}{
		"type":            "unsubscribed",
		"conversation_id": convID,
	}
	log.Println("Unsubscribing from conversation ID:", convID)
	jsonData, _ := json.Marshal(resp)
	select {
	case c.Send <- jsonData:
	case <-time.After(5 * time.Second):
		log.Println("Timeout sending unsubscribe response")
	}
}

// Update sendErrorToClient to use a timeout select like other sends
func sendErrorToClient(c *domain.Client, errorMsg string) {
	errorResponse := map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	}
	jsonData, _ := json.Marshal(errorResponse)
	select {
	case c.Send <- jsonData:
	case <-time.After(5 * time.Second):
		log.Printf("Timeout sending error message to client %d: %s", c.UserID, errorMsg)
	}
}

func parseConversationID(v interface{}) (uint64, error) {
	if v == nil {
		return 0, fmt.Errorf("conversation_id missing")
	}
	switch t := v.(type) {
	case float64:
		return uint64(t), nil
	case float32:
		return uint64(t), nil
	case int:
		return uint64(t), nil
	case int64:
		return uint64(t), nil
	case uint64:
		return t, nil
	case string:
		if t == "" {
			return 0, fmt.Errorf("conversation_id empty")
		}
		u, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			return 0, err
		}
		return u, nil
	default:
		return 0, fmt.Errorf("invalid conversation_id type %T", v)
	}
}

func (s *appService) handleSendMessage(c *domain.Client, payload map[string]interface{}) {
	conversationID, err := parseConversationID(payload["conversation_id"])
	if err != nil {
		sendErrorToClient(c, "Invalid conversation ID")
		return
	}

	// Check if conversation exists and user has access
	conversation, err := c.Repository.GetConversationByID(conversationID)
	if err != nil {
		sendErrorToClient(c, "Conversation not found")
		return
	}

	// Authorization check - verify user can send to this conversation
	user, err := c.Repository.GetUserByID(c.UserID)
	if err != nil {
		sendErrorToClient(c, "Failed to get user")
		return
	}

	// Check if user has permission (admin or conversation participant)
	isAdmin := user.Role == models.RoleAdmin
	isParticipant := conversation.CustomerID == c.UserID || conversation.AdminID == c.UserID

	if !isAdmin && !isParticipant {
		sendErrorToClient(c, "Access denied - not your conversation")
		return
	}

	// require subscription to prevent unauthorized sends / avoid silent broadcasts
	if c.ConversationIDs == nil || !c.ConversationIDs[conversationID] {
		sendErrorToClient(c, "You are not subscribed to this conversation")
		return
	}

	// If conversation is closed and reopen it
	if conversation.Status == models.StatusClosed {
		ctx := context.Background()
		if err := s.ReopenConversation(ctx, conversationID); err != nil {
			sendErrorToClient(c, "Failed to reopen conversation")
			return
		}
	}

	userID := c.UserID

	// Extract message text
	text, ok := payload["text"].(string)
	if !ok || text == "" {
		sendErrorToClient(c, "Message text is required")
		return
	}

	// Create message (no deleted_at or purge_at for new messages)
	Message := &models.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		MessageText:    text,
		CreatedAt:      time.Now(),
		DeletedAt:      nil, // New message is not deleted
		PurgeAt:        nil, // No purge scheduled
	}

	savedMessage, err := c.Repository.SaveMessage(Message)
	if err != nil {
		sendErrorToClient(c, "Failed to save message")
		return
	}

	// Update conversation's last_message_at
	err = c.Repository.UpdateConversationLastMessage(conversationID)
	if err != nil {
		sendErrorToClient(c, "Failed to update conversation")
	}

	// Handle admin notification for unread messages
	go s.handleAdminNotification(c, conversationID, savedMessage.ID)

	// Broadcast to conversation participants only (hub routes to connected clients)
	c.Hub.Broadcast <- &domain.Message{
		ConversationID: conversationID,
		Data:           savedMessage,
		Type:           "new_message",
	}
}

// handleAdminNotification increments unread count if sender is not admin and admin is not in the room
func (s *appService) handleAdminNotification(c *domain.Client, conversationID uint64, messageID uint64) {
	sender, err := c.Repository.GetUserByID(c.UserID)
	if err != nil || sender == nil || sender.Role == models.RoleAdmin {
		return
	}

	conversation, err := c.Repository.GetConversationByID(conversationID)
	if err != nil || conversation == nil {
		log.Printf("Failed to get conversation %d for notification: %v", conversationID, err)
		return
	}

	adminID := conversation.AdminID

	// Check if admin is in room
	s.hub.Mu.RLock()
	clientsInConv := s.hub.Conversations[conversationID]
	_, adminInRoom := clientsInConv[uint64(adminID)]
	// Also check global connection
	adminClient := s.hub.Clients[uint64(adminID)]
	s.hub.Mu.RUnlock()

	// Update unread count only if not in room
	var unreadAfter uint
	if !adminInRoom {
		state, err := c.Repository.GetAdminConversationState(adminID, conversationID)
		if err != nil {
			log.Printf("Failed to get admin conversation state for admin %d, conversation %d: %v", adminID, conversationID, err)
			return
		}
		if state != nil {
			if err := c.Repository.IncrementUnreadCount(state); err != nil {
				log.Printf("Failed to increment unread count for admin %d, conversation %d: %v", adminID, conversationID, err)
			} else {
				unreadAfter = state.UnreadCount
				log.Printf("Incremented unread count for admin %d in conversation %d (new message: %d)", adminID, conversationID, messageID)
			}
		}
	}

	// Send live notification frame if admin has websocket connection
	if adminClient != nil {
		notification := map[string]interface{}{
			"conversation_id": conversationID,
			"message_id":      messageID,
			"unread_count":    unreadAfter,
			"sender_id":       c.UserID,
			"in_room":         adminInRoom,
		}
		s.sendDirect(adminClient, "admin_notification", notification)
	}
}

func (s *appService) GetTicketNotifications() map[string]interface{} {
	customerCount, sellerCount, err := s.repo.GetOpenTicketCountsByType()
	if err != nil {
		return map[string]interface{}{
			"customer_open_tickets": 0,
			"seller_open_tickets":   0,
			"total_open_tickets":    0,
			"error":                 err.Error(),
		}
	}

	return map[string]interface{}{
		"customer_open_tickets": customerCount,
		"seller_open_tickets":   sellerCount,
		"total_open_tickets":    customerCount + sellerCount,
	}
}
