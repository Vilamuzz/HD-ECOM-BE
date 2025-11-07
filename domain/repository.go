package domain

import (
	"app/domain/models"
)

type AppRepository interface {
	// User operations
	GetUserByID(id int64) (*models.User, error)
	CreateUser(user *models.User) error

	// Conversation operations
	GetConversationParticipants(conversationID uint64) ([]uint64, error)
	CreateConversation(conversation *models.Conversation) (*models.Conversation, error)
	GetAdminConversations(adminID uint64) ([]models.Conversation, error)
	GetUserConversations(userID uint64) ([]models.Conversation, error)
	FindActiveConversationForCustomer(customerID uint64) (*models.Conversation, error)
	UpdateConversationLastMessage(conversationID uint64) error

	// Chat message operations
	SaveChatMessage(message *models.ChatMessage) (*models.ChatMessage, error)
	GetChatMessages(conversationID uint64) ([]models.ChatMessage, error)
	GetChatMessagesByConversationID(conversationID string) ([]models.ChatMessage, error)

	// Admin availability operations
	GetAdminConversationCount(adminID uint64) (int, error)
	GetAdminWithLeastConversations() (*models.User, error)
	CreateAdminAvailability(adminAvailability *models.AdminAvailability) error
	UpdateAdminAvailability(adminAvailability *models.AdminAvailability) error
}
