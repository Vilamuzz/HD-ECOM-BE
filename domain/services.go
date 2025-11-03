package domain

import (
	"app/helpers"

	"github.com/gin-gonic/gin"
)

type AppService interface {
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *Message)
	SendToRecipients(message *Message)
	JoinConversation(client *Client, conversationID string)
	ServeWebSocket(ctx *gin.Context) helpers.Response
}
