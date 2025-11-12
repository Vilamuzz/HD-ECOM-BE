package domain

import (
	"app/domain/models"
	"context"
)

type AppRepository interface {
	// User operations
	GetUserByID(id uint64) (*models.User, error)
	CreateUser(user *models.User) error

	// Conversation operations
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetAdminConversations(adminID uint8) ([]models.Conversation, error)
	GetCustomerConversations(userID uint64) ([]models.Conversation, error)
	UpdateConversationLastMessage(conversationID uint64) error

	// Chat message operations
	// Modified for cursor pagination: limit and cursor input, returns next cursor
	GetMessageHistory(conversationID uint64, limit int, cursor string) ([]models.Message, string, error)
	SaveMessage(message *models.Message) (*models.Message, error)

	// Admin availability operations
	GetAdminAvailabilityByAdminID() (*models.AdminAvailability, error)
	CreateAdminAvailability(adminAvailability *models.AdminAvailability) error
	IncrementAdminConversationCount(adminID uint8) error
}
