package repositories

import (
	"app/domain/models"
	"context"
	"time"
)

func (r *appRepository) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	err := r.Conn.WithContext(ctx).Create(conversation).Error
	return err
}

func (r *appRepository) GetAdminConversations(adminID uint64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.Conn.Where("admin_id = ?", adminID).Order("last_message_at DESC").Find(&conversations).Error
	return conversations, err
}

func (r *appRepository) CloseConversation(ctx context.Context, conversationID uint64) error {
	return r.Conn.WithContext(ctx).Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("status", "closed").Error
}

func (r *appRepository) ReopenConversation(ctx context.Context, conversationID uint64) error {
	return r.Conn.WithContext(ctx).Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("status", models.StatusOpen).Error
}

func (r *appRepository) GetCustomerConversations(customerID uint64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.Conn.Where("customer_id = ?", customerID).
		Order("last_message_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func (r *appRepository) UpdateConversationLastMessage(conversationID uint64) error {
	return r.Conn.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_at", time.Now()).Error
}

func (r *appRepository) GetConversationByID(conversationID uint64) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.Conn.Where("id = ?", conversationID).First(&conversation).Error
	return &conversation, err
}
