package websocket

import (
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
			Hub:        hub,
			Conn:       &WebSocketConnectionWrapper{Conn: conn},
			Send:       make(chan []byte, 256),
			UserID:     strconv.FormatInt(user.ID, 10),
			Name:       user.Username,
			Repository: repo,
		}

		// Register client
		hub.Register <- client

		// Start pumps
		StartClientPumps(client)
	}
}
