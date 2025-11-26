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
	GetConversationByID(conversationID uint64) (*models.Conversation, error)
	GetAdminConversations(adminID uint64) ([]models.Conversation, error)
	GetCustomerConversations(userID uint64) ([]models.Conversation, error)
	UpdateConversationLastMessage(conversationID uint64) error
	CloseConversation(ctx context.Context, conversationID uint64) error
	ReopenConversation(ctx context.Context, conversationID uint64) error

	// Chat message operations
	GetMessageHistory(conversationID uint64, limit int, cursor string) ([]models.Message, string, error)
	GetMessageHistoryForAdmin(conversationID uint64, limit int, cursor string) ([]models.Message, string, error)
	SaveMessage(message *models.Message) (*models.Message, error)
	SoftDeleteConversationMessages(conversationID uint64, purgeAfterDays int) error
	ResetPurgeTimestamp(conversationID uint64) error
	PermanentlyDeleteExpiredMessages() error

	// Admin availability operations
	GetAdminAvailabilityByAdminID() (*models.AdminAvailability, error)
	CreateAdminAvailability(adminAvailability *models.AdminAvailability) error
	IncrementAdminConversationCount(adminID uint64) error
	DecrementAdminConversationCount(adminID uint64) error

	// Admin conversation state operations
	CreateAdminConversationState(adminID uint64, conversationID uint64) error
	GetAdminConversationState(adminID uint64, conversationID uint64) (*models.AdminConversationState, error)
	GetAdminConversationStatesByAdminID(adminID uint64) ([]models.AdminConversationState, error)
	IncrementUnreadCount(state *models.AdminConversationState) error
	ResetState(state *models.AdminConversationState, lastMessageID uint64) error
}
