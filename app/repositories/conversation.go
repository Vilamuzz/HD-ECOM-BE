package repositories

import (
	"app/domain/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (r *appRepository) GetConversationParticipants(conversationID uint64) ([]uint64, error) {
	var conversation models.Conversation
	err := r.Conn.First(&conversation, conversationID).Error
	if err != nil {
		return nil, err
	}
	participants := []uint64{conversation.UserID}
	if conversation.AdminID != 0 {
		participants = append(participants, conversation.AdminID)
	} else {
		var users []models.User
		err := r.Conn.Where("role IN ?", []string{"admin"}).Find(&users).Error
		if err == nil {
			for _, user := range users {
				participants = append(participants, user.ID)
			}
		}
	}
	return participants, nil
}

func (r *appRepository) CreateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	err := r.Conn.Create(conversation).Error
	return conversation, err
}

func (r *appRepository) GetAdminConversations(adminID uint64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.Conn.Where("admin_id = ?", adminID).Order("last_message_at DESC").Find(&conversations).Error
	return conversations, err
}

func (r *appRepository) GetUserConversations(userID uint64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.Conn.Where("customer_id = ? OR admin_id = ?", userID, userID).
		Order("last_message_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func (r *appRepository) GetChatMessagesByConversationID(conversationID string) ([]models.ChatMessage, error) {
	id, err := strconv.ParseUint(conversationID, 10, 64)
	if err != nil {
		return nil, err
	}
	return r.GetChatMessages(id)
}

func (r *appRepository) FindActiveConversationForCustomer(userID uint64) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.Conn.Where("user_id = ?", userID).
		Order("last_message_at DESC").
		First(&conversation).Error

	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *appRepository) UpdateConversationLastMessage(conversationID uint64) error {
	return r.Conn.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_at", time.Now()).Error
}

func (r *appRepository) GetAdminConversationCount(adminID uint64) (int, error) {
	var count int64
	err := r.Conn.Model(&models.Conversation{}).
		Where("admin_id = ?", adminID).
		Count(&count).Error
	return int(count), err
}

func (r *appRepository) GetAdminWithLeastConversations() (*models.User, error) {
	var admins []models.User
	err := r.Conn.Where("role = ?", "admin").Find(&admins).Error
	if err != nil {
		return nil, err
	}

	if len(admins) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Find admin with least conversations
	var selectedAdmin *models.User
	minCount := int(^uint(0) >> 1)

	for i := range admins {
		count, err := r.GetAdminConversationCount(admins[i].ID)
		if err != nil {
			continue
		}

		if count < minCount {
			minCount = count
			selectedAdmin = &admins[i]
		}
	}

	return selectedAdmin, nil
}
