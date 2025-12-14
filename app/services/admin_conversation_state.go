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

	// Get ticket notification counts
	customerCount, sellerCount, ticketErr := s.repo.GetOpenTicketCountsByType()

	// Prepare response data
	responseData := map[string]interface{}{
		"conversation_states": states,
		"ticket_notifications": map[string]interface{}{
			"customer_open_tickets": customerCount,
			"seller_open_tickets":   sellerCount,
			"total_open_tickets":    customerCount + sellerCount,
		},
	}

	// Add error info if ticket count failed
	if ticketErr != nil {
		responseData["ticket_notifications"].(map[string]interface{})["error"] = ticketErr.Error()
	}

	return helpers.NewResponse(http.StatusOK, "Successfully retrieved conversation states and ticket notifications", nil, responseData)
}
