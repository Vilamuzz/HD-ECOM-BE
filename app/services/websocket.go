package services

import (
	"app/domain"
	"app/domain/models"
	"app/helpers"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// DefaultUpgrader provides a default WebSocket upgrader configuration
var DefaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// ServeWebSocket handles WebSocket connection upgrades
func (s *appService) ServeWebSocket(ctx *gin.Context) helpers.Response {

	// Get user from context (already validated by Auth middleware)
	user, ok := ctx.Value("currentUser").(models.User)
	if !ok {
		return helpers.NewResponse(http.StatusUnauthorized, "Unauthorized", nil, nil)
	}

	// Upgrade to WebSocket
	conn, err := DefaultUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return helpers.NewResponse(http.StatusInternalServerError, "WebSocket upgrade failed", nil, nil)
	}

	// Create client
	client := &domain.Client{
		Hub:             s.hub,
		Conn:            &WebSocketConnectionWrapper{Conn: conn},
		Send:            make(chan []byte, 256),
		UserID:          strconv.FormatInt(user.ID, 10),
		Name:            user.Username,
		Repository:      s.repo,
		ConversationIDs: make(map[string]bool),
	}

	// Register client
	s.hub.Register <- client

	// Send connection success and setup conversations
	go s.sendInitialData(client, user)

	// Start pumps
	s.StartClientPumps(client)

	return helpers.NewResponse(http.StatusOK, "WebSocket connection established", nil, nil)
}

// sendInitialData sends connection confirmation and loads conversations
func (s *appService) sendInitialData(client *domain.Client, user models.User) {
	// Send connection confirmation
	connectResponse := map[string]interface{}{
		"type":    "connected",
		"user_id": client.UserID,
		"message": "Successfully connected to WebSocket",
	}
	connectData, _ := json.Marshal(connectResponse)
	client.Send <- connectData

	log.Printf("[WS] Sent connection confirmation to user %s", client.UserID)

	userID, _ := strconv.ParseInt(client.UserID, 10, 64)

	if user.Role == "admin" {
		// Admin: Load all conversations they're involved in
		conversations, err := s.repo.GetUserConversations(userID)
		if err != nil {
			log.Printf("[WS] Error loading conversations for admin/agent %s: %v", client.UserID, err)
			return
		}

		log.Printf("[WS] Admin/Agent %s loading %d conversations", client.UserID, len(conversations))

		for _, conv := range conversations {
			convID := strconv.FormatInt(conv.ID, 10)
			s.JoinConversation(client, convID)

			// Load and send message history for each conversation
			messages, _ := s.repo.GetChatMessages(conv.ID)
			if len(messages) > 0 {
				historyResponse := map[string]interface{}{
					"type":            "message_history",
					"conversation_id": conv.ID,
					"messages":        messages,
					"count":           len(messages),
				}
				historyData, _ := json.Marshal(historyResponse)
				client.Send <- historyData
			}
		}

		log.Printf("[WS] Admin/Agent %s joined %d conversations", client.UserID, len(conversations))

	} else {
		// Customer: Find or create their active conversation
		existingConv, err := s.repo.FindActiveConversationForCustomer(userID)

		if err != nil {
			log.Printf("[WS] No active conversation found for customer %s, will create on first message", client.UserID)
			return
		}

		log.Printf("[WS] Found active conversation %d for customer %s", existingConv.ID, client.UserID)

		convID := strconv.FormatInt(existingConv.ID, 10)
		s.JoinConversation(client, convID)

		// Send conversation info
		convResponse := map[string]interface{}{
			"type":            "conversation_loaded",
			"conversation_id": existingConv.ID,
			"status":          existingConv.Status,
		}
		convData, _ := json.Marshal(convResponse)
		client.Send <- convData

		// Load and send message history
		messages, err := s.repo.GetChatMessages(existingConv.ID)
		if err != nil {
			log.Printf("[WS] Error loading messages for conversation %d: %v", existingConv.ID, err)
			return
		}

		if len(messages) > 0 {
			historyResponse := map[string]interface{}{
				"type":            "message_history",
				"conversation_id": existingConv.ID,
				"messages":        messages,
				"count":           len(messages),
			}
			historyData, _ := json.Marshal(historyResponse)
			client.Send <- historyData

			log.Printf("[WS] Sent %d messages to customer %s for conversation %d",
				len(messages), client.UserID, existingConv.ID)
		}
	}
}
