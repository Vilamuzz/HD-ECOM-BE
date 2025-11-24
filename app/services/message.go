package services

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
)

func (s *appService) GetMessageHistory(conversationID uint64, limit int, cursor string, isAdmin bool) helpers.Response {
	var messages []models.Message
	var nextCursor string
	var err error

	// Admin sees all messages including soft-deleted, customers only see non-deleted
	if isAdmin {
		messages, nextCursor, err = s.repo.GetMessageHistoryForAdmin(conversationID, limit, cursor)
	} else {
		messages, nextCursor, err = s.repo.GetMessageHistory(conversationID, limit, cursor)
	}

	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Failed to get messages", nil, nil)
	}

	// Return messages and next_cursor in response data
	data := map[string]interface{}{
		"messages":    messages,
		"next_cursor": nextCursor,
		"limit":       limit,
	}

	return helpers.NewResponse(http.StatusOK, "Successfully retrieved messages", nil, data)
}
