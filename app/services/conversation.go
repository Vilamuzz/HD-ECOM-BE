package services

import (
	"app/domain/models"
	"app/helpers"
	"context"
	"net/http"
	"strconv"
	"time"
)

func (s *appService) GetConversations(claim models.User) helpers.Response {
	var conversations []models.Conversation
	var err error
	if claim.Role == "admin" {
		conversations, err = s.repo.GetAdminConversations(claim.ID)

		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "Failed to get conversations", nil, nil)
		}

		type ConversationWithCustomer struct {
			models.Conversation
			CustomerName  string `json:"customer_name"`
			CustomerEmail string `json:"customer_email"`
			CustomerRole  string `json:"customer_role"`
		}

		list := make([]ConversationWithCustomer, 0, len(conversations))
		for _, conv := range conversations {
			customerID := conv.CustomerID

			customerName := ""
			customerEmail := ""
			customerRole := ""
			if customerID != 0 {
				u, err := s.repo.GetUserByID(customerID)
				if err == nil && u != nil {
					customerName = u.Username
					customerEmail = u.Email
					customerRole = string(u.Role)
				}
			}

			list = append(list, ConversationWithCustomer{
				Conversation:  conv,
				CustomerName:  customerName,
				CustomerEmail: customerEmail,
				CustomerRole:  customerRole,
			})
		}
		return helpers.NewResponse(http.StatusOK, "Successfully get conversation", nil, list)
	} else {
		conversations, err := s.repo.GetCustomerConversations(claim.ID)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "failed to get conversations", nil, nil)
		}
		if len(conversations) > 0 {
			return helpers.NewResponse(http.StatusOK, "conversation already exists", nil, conversations)
		}
		return helpers.NewResponse(http.StatusOK, "no conversations found", nil, conversations)
	}
}

func (s *appService) CreateCustomerConversation(ctx context.Context, claim models.User) helpers.Response {
	const createTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	admin, err := s.repo.GetAdminAvailabilityByAdminID()
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to get admin availability", nil, nil)
	}

	if admin == nil || admin.AdminID == 0 {
		return helpers.NewResponse(http.StatusInternalServerError, "no available admin found", nil, nil)
	}

	now := time.Now()
	createdConversation := &models.Conversation{
		CustomerID: claim.ID,
		AdminID:    admin.AdminID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateConversation(ctx, createdConversation); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to create conversation", nil, nil)
	}

	if err := s.repo.IncrementAdminConversationCount(admin.AdminID); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to update admin availability", nil, nil)
	}

	if err := s.repo.CreateAdminConversationState(admin.AdminID, createdConversation.ID); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to create admin conversation state", nil, nil)
	}

	// Notify admin if they have a websocket connection
	s.hub.Mu.RLock()
	adminClient := s.hub.Clients[uint64(admin.AdminID)]
	s.hub.Mu.RUnlock()
	if adminClient != nil {
		// Enrich payload with full conversation data, customer/admin info and admin state
		var adminUser *models.User
		var customerUser *models.User

		if u, err := s.repo.GetUserByID(uint64(admin.AdminID)); err == nil && u != nil {
			adminUser = u
		}
		if u, err := s.repo.GetUserByID(createdConversation.CustomerID); err == nil && u != nil {
			customerUser = u
		}

		adminState, _ := s.repo.GetAdminConversationState(admin.AdminID, createdConversation.ID)

		payload := map[string]interface{}{
			"conversation": createdConversation,
			"customer":     customerUser,
			"admin":        adminUser,
			"admin_state":  adminState,
		}
		s.sendDirect(adminClient, "conversation_created", payload)
	}

	return helpers.NewResponse(http.StatusCreated, "successfully created conversation", nil, createdConversation)
}

func (s *appService) CloseConversation(ctx context.Context, claim models.User, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	conversationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || conversationID == 0 {
		return helpers.NewResponse(http.StatusBadRequest, "invalid conversation ID", nil, nil)
	}

	// Close the conversation
	err = s.repo.CloseConversation(ctx, conversationID)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to close conversation", nil, nil)
	}

	// Soft delete all messages and set purge countdown (e.g., 30 days)
	purgeAfterDays := 30
	if err := s.repo.SoftDeleteConversationMessages(conversationID, purgeAfterDays); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to soft delete messages", nil, nil)
	}

	// Decrement admin conversation count
	if err := s.repo.DecrementAdminConversationCount(claim.ID); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to update admin availability", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "successfully closed conversation and messages scheduled for deletion", nil, nil)
}

// Add new method to reopen conversation
func (s *appService) ReopenConversation(ctx context.Context, conversationID uint64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Reopen conversation by setting status back to open
	err := s.repo.ReopenConversation(ctx, conversationID)
	if err != nil {
		return err
	}

	err = s.repo.ResetPurgeTimestamp(conversationID)
	if err != nil {
		return err
	}

	// Reset purge countdown but keep soft delete
	return nil
}

func (s *appService) GetConversationByID(conversationID uint64) (*models.Conversation, error) {
	return s.repo.GetConversationByID(conversationID)
}
