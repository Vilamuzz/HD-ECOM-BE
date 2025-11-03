package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"app/domain"
	"app/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

type JWTClaims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// ServeWebSocket handles WebSocket connection upgrades
func ServeWebSocket(hub *domain.Hub, repo domain.AppRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (already validated by Auth middleware)
		currentUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		user, ok := currentUser.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
			return
		}

		// Upgrade to WebSocket
		conn, err := DefaultUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		// Create client
		client := &domain.Client{
			Hub:             hub,
			Conn:            &WebSocketConnectionWrapper{Conn: conn},
			Send:            make(chan []byte, 256),
			UserID:          strconv.FormatInt(user.ID, 10),
			Name:            user.Username,
			Repository:      repo,
			ConversationIDs: make(map[string]bool),
		}

		// Register client
		hub.Register <- client

		// Send connection success and setup conversations
		go sendInitialData(client, user, repo, hub)

		// Start pumps
		StartClientPumps(client)
	}
}

// sendInitialData sends connection confirmation and loads conversations
func sendInitialData(client *domain.Client, user models.User, repo domain.AppRepository, hub *domain.Hub) {
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
		conversations, err := repo.GetUserConversations(userID)
		if err != nil {
			log.Printf("[WS] Error loading conversations for admin/agent %s: %v", client.UserID, err)
			return
		}

		log.Printf("[WS] Admin/Agent %s loading %d conversations", client.UserID, len(conversations))

		// Join all their conversations
		hubSvc := NewHubService(hub)
		for _, conv := range conversations {
			convID := strconv.FormatInt(conv.ID, 10)
			hubSvc.JoinConversation(client, convID)

			// Load and send message history for each conversation
			messages, _ := repo.GetChatMessages(conv.ID)
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
		existingConv, err := repo.FindActiveConversationForCustomer(userID)

		if err != nil {
			log.Printf("[WS] No active conversation found for customer %s, will create on first message", client.UserID)
			return
		}

		log.Printf("[WS] Found active conversation %d for customer %s", existingConv.ID, client.UserID)

		// Join the conversation
		hubSvc := NewHubService(hub)
		convID := strconv.FormatInt(existingConv.ID, 10)
		hubSvc.JoinConversation(client, convID)

		// Send conversation info
		convResponse := map[string]interface{}{
			"type":            "conversation_loaded",
			"conversation_id": existingConv.ID,
			"status":          existingConv.Status,
		}
		convData, _ := json.Marshal(convResponse)
		client.Send <- convData

		// Load and send message history
		messages, err := repo.GetChatMessages(existingConv.ID)
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
