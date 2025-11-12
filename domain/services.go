package domain

import (
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
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
	ServeWebSocket(ctx *gin.Context) helpers.Response

	GetConversations(claim jwt_helpers.Claims) helpers.Response
	CreateCustomerConversation(ctx context.Context, claim jwt_helpers.Claims) helpers.Response
	GetMessageHistory(conversationID uint64, limit int, cursor string) helpers.Response
}
