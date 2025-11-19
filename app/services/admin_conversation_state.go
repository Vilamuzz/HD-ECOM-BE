package services

import "app/domain/models"

func (s *appService) GetAdminConversationState(adminID uint8, conversationID uint64) (*models.AdminConversationState, error) {
	return s.repo.GetAdminConversationState(adminID, conversationID)
}
