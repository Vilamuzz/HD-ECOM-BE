package services

import (
	"app/helpers"
	"net/http"
)

func (s *appService) GetMessageHistory(conversationID uint64, limit int, cursor string) helpers.Response {
	messages, nextCursor, err := s.repo.GetMessageHistory(conversationID, limit, cursor)
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
