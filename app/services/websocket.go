package services

import (
	"app/domain"
	"app/domain/models"
	"app/helpers"
	"encoding/json"
	"net/http"

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
		UserID:          user.ID,
		Name:            user.Username,
		Repository:      s.repo,
		ConversationIDs: make(map[uint64]bool),
	}

	// Register client
	s.hub.Register <- client

	// Send connection success and setup conversations
	go s.sendInitialData(client)

	// Start pumps
	s.StartClientPumps(client)

	return helpers.NewResponse(http.StatusOK, "WebSocket connection established", nil, nil)
}

// sendInitialData sends connection confirmation and loads conversations
func (s *appService) sendInitialData(client *domain.Client) {
	connectResponse := map[string]interface{}{
		"type":    "connected",
		"user_id": client.UserID,
		"message": "Successfully connected to WebSocket",
	}
	connectData, _ := json.Marshal(connectResponse)
	client.Send <- connectData
}
