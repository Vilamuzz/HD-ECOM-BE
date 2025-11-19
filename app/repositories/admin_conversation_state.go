package repositories

import "app/domain/models"

func (r *appRepository) CreateAdminConversationState(adminID uint8, conversationID uint64) error {
	state := &models.AdminConversationState{
		AdminID:        adminID,
		ConversationID: conversationID,
		UnreadCount:    0,
		LastMessageID:  0,
	}
	return r.Conn.Create(state).Error
}

func (r *appRepository) GetAdminConversationState(adminID uint8, conversationID uint64) (*models.AdminConversationState, error) {
	var state models.AdminConversationState
	err := r.Conn.Where("admin_id = ? AND conversation_id = ?", adminID, conversationID).First(&state).Error
	return &state, err
}

func (r *appRepository) IncrementUnreadCount(state *models.AdminConversationState) error {
	state.UnreadCount++
	return r.Conn.Save(state).Error
}

func (r *appRepository) ResetState(state *models.AdminConversationState, lastMessageID uint64) error {
	state.UnreadCount = 0
	state.LastMessageID = lastMessageID
	return r.Conn.Save(state).Error
}
