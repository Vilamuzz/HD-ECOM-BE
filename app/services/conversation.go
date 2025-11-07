package services

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
)

func (s *appService) GetConversations(user models.User) helpers.Response {
	var conversations []models.Conversation
	var err error

	if user.Role == "admin" {
		conversations, err = s.repo.GetAdminConversations(user.ID)
	} else {
		conversations, err = s.repo.GetUserConversations(user.ID)
	}

	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Failed to get conversations", nil, nil)
	}
	return helpers.NewResponse(http.StatusOK, "Successfuly get conversation", nil, conversations)
}

func (s *appService) GetConversationMessages(conversationID string) helpers.Response {
	messages, err := s.repo.GetChatMessagesByConversationID(conversationID)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Failed to get messages", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Successfully retrieved messages", nil, messages)
}
