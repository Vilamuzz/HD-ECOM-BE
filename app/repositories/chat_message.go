package repositories

import "app/domain/models"

func (r *appRepository) SaveChatMessage(message *models.ChatMessage) (*models.ChatMessage, error) {
	err := r.Conn.Create(message).Error
	return message, err
}

func (r *appRepository) GetChatMessages(conversationID int64) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.Conn.Where("conversation_id = ?", conversationID).Find(&messages).Error
	return messages, err
}
