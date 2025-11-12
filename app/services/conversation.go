package services

import (
	"app/domain/models"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"net/http"
	"time"
)

func (s *appService) GetConversations(claim jwt_helpers.Claims) helpers.Response {
	var conversations []models.Conversation
	var err error
	if claim.Role != 0 {
		conversations, err = s.repo.GetAdminConversations(uint8(claim.UserID))

		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "Failed to get conversations", nil, nil)
		}

		type ConversationWithCustomer struct {
			Conversation models.Conversation `json:"conversation"`
			CustomerID   uint64              `json:"customer_id"`
			CustomerName string              `json:"customer_name"`
		}

		list := make([]ConversationWithCustomer, 0, len(conversations))
		for _, conv := range conversations {
			customerID := conv.UserID

			customerName := ""
			if customerID != 0 {
				u, err := s.repo.GetUserByID(customerID)
				if err == nil && u != nil {
					customerName = u.Username
				}
			}

			list = append(list, ConversationWithCustomer{
				Conversation: conv,
				CustomerID:   customerID,
				CustomerName: customerName,
			})
		}
		return helpers.NewResponse(http.StatusOK, "Successfully get conversation", nil, list)
	} else {
		conversations, err := s.repo.GetCustomerConversations(claim.UserID)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "failed to get conversations", nil, nil)
		}
		if len(conversations) > 0 {
			return helpers.NewResponse(http.StatusOK, "conversation already exists", nil, conversations)
		}
		return helpers.NewResponse(http.StatusOK, "no conversations found", nil, conversations)
	}
}

func (s *appService) CreateCustomerConversation(ctx context.Context, claim jwt_helpers.Claims) helpers.Response {
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
		UserID:    claim.UserID,
		AdminID:   admin.AdminID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateConversation(ctx, createdConversation); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to create conversation", nil, nil)
	}

	if err := s.repo.IncrementAdminConversationCount(admin.AdminID); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "failed to update admin availability", nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "successfully created conversation", nil, createdConversation)
}
