package services

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
)

func (s *appService) GetAdminConversationState(adminID uint64, conversationID uint64) (*models.AdminConversationState, error) {
	return s.repo.GetAdminConversationState(adminID, conversationID)
}

func (s *appService) GetAdminListConversationStates(claim models.User) helpers.Response {
	states, err := s.repo.GetAdminConversationStatesByAdminID(claim.ID)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Failed to retrieve conversation states", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Successfully retrieved conversation states", nil, states)
}
