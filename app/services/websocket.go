package services

import (
	"app/domain"
	"app/domain/models"
	"app/helpers"
	"context"
	"encoding/json"
	"errors"
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
		return true
	},
}

// ServeWebSocket handles WebSocket connection upgrades
func (s *appService) ServeWebSocket(ctx *gin.Context) helpers.Response {

	user, ok := ctx.Value("userData").(models.User)
	if !ok {
		return helpers.NewResponse(http.StatusUnauthorized, "Unauthorized", nil, nil)
	}

	// Upgrade to WebSocket
	conn, err := DefaultUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "WebSocket upgrade failed", nil, nil)
	}

	// Create client
	client := &domain.Client{
		Hub:             s.hub,
		Conn:            &WebSocketConnectionWrapper{Conn: conn},
		Send:            make(chan []byte, 256),
		UserID:          strconv.FormatUint(user.ID, 10),
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
	connectResponse := map[string]interface{}{
		"type":    "connected",
		"user_id": client.UserID,
		"message": "Successfully connected to WebSocket",
	}
	connectData, _ := json.Marshal(connectResponse)
	client.Send <- connectData

	userID, _ := strconv.ParseUint(client.UserID, 10, 64)

	if user.Role == "admin" {
		conversations, err := s.repo.GetUserConversations(userID)
		if err != nil && errors.Is(err, context.Canceled) {
			return
		}

		for _, conv := range conversations {
			convID := strconv.FormatUint(conv.ID, 10)
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

	} else {
		existingConv, err := s.repo.FindActiveConversationForCustomer(userID)

		if err != nil {
			return
		}

		convID := strconv.FormatUint(existingConv.ID, 10)
		s.JoinConversation(client, convID)

		// Send conversation info
		convResponse := map[string]interface{}{
			"type":            "conversation_loaded",
			"conversation_id": existingConv.ID,
		}
		convData, _ := json.Marshal(convResponse)
		client.Send <- convData

		// Load and send message history
		messages, err := s.repo.GetChatMessages(existingConv.ID)
		if err != nil {
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
		}
	}
}
