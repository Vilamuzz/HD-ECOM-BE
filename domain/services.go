package domain

import (
	"app/domain/models"
	"app/helpers"
	"context"

	"github.com/gin-gonic/gin"
)

type AppService interface {
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *Message)
	SendToRecipients(message *Message)
	JoinConversation(client *Client, conversationID uint64)
	ServeWebSocket(ctx *gin.Context)
	GetConversations(claim models.User) helpers.Response
	CreateCustomerConversation(ctx context.Context, claim models.User) helpers.Response
	GetMessageHistory(conversationID uint64, limit int, cursor string) helpers.Response
}
