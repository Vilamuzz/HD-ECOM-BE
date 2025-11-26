package domain

import (
	"app/domain/models"
	"app/helpers"
	"context"

	"github.com/gin-gonic/gin"
)

type AppService interface {
	// WebSocket management
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *Message)
	SendToRecipients(message *Message)
	JoinConversation(client *Client, conversationID uint64)
	ServeWebSocket(ctx *gin.Context)

	// Conversation management
	GetConversations(claim models.User) helpers.Response
	CreateCustomerConversation(ctx context.Context, claim models.User) helpers.Response
	CloseConversation(ctx context.Context, claim models.User, id string) helpers.Response
	ReopenConversation(ctx context.Context, conversationID uint64) error

	// Message management
	GetMessageHistory(conversationID uint64, limit int, cursor string, isAdmin bool) helpers.Response

	// Admin state management
	GetAdminConversationState(adminID uint64, conversationID uint64) (*models.AdminConversationState, error)
	GetAdminListConversationStates(claim models.User) helpers.Response
}
