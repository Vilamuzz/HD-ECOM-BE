package domain

import "app/domain/models"

type AppRepository interface {
	// User operations
	GetUserByID(id int64) (*models.User, error)
	CreateUser(user *models.User) error

	// Conversation operations
	GetConversation(id int64) (*models.Conversation, error)
	GetConversationParticipants(conversationID int64) ([]int64, error)
	CreateConversation(conversation *models.Conversation) (*models.Conversation, error)
	GetAllConversations() ([]models.Conversation, error)
	GetUserConversations(userID int64) ([]models.Conversation, error)
	AssignConversationToAgent(conversationID int64, agentID int64) error
	FindActiveConversationForCustomer(customerID int64) (*models.Conversation, error)
	UpdateConversationLastMessage(conversationID int64) error

	// Chat message operations
	SaveChatMessage(message *models.ChatMessage) (*models.ChatMessage, error)
	GetChatMessages(conversationID int64) ([]models.ChatMessage, error)
	GetChatMessagesByConversationID(conversationID string) ([]models.ChatMessage, error)
}
