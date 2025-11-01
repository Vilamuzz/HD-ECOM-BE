package repositories

import (
	"app/domain/models"
	"log"
	"strconv"
	"time"
)

func (r *appRepository) GetConversation(id int64) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.Conn.First(&conversation, id).Error
	return &conversation, err
}

func (r *appRepository) GetConversationParticipants(conversationID int64) ([]int64, error) {
	var conversation models.Conversation
	err := r.Conn.First(&conversation, conversationID).Error
	if err != nil {
		return nil, err
	}

	participants := []int64{conversation.CostumerID}

	if conversation.AgentID != 0 {
		// If agent is assigned, include them
		participants = append(participants, conversation.AgentID)
	} else {
		// If no agent assigned, include all admins and agents
		var users []models.User
		err := r.Conn.Where("role IN ?", []string{"admin", "agent"}).Find(&users).Error
		if err == nil {
			for _, user := range users {
				participants = append(participants, user.ID)
			}
		}
	}

	log.Printf("[REPO] GetConversationParticipants - ConversationID: %d, CustomerID: %d, AgentID: %d, Participants: %v",
		conversationID, conversation.CostumerID, conversation.AgentID, participants)

	return participants, nil
}

func (r *appRepository) CreateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	log.Printf("[REPO] Creating conversation: %+v", conversation)
	err := r.Conn.Create(conversation).Error
	if err != nil {
		log.Printf("[REPO] Error creating conversation: %v", err)
		return conversation, err
	}
	log.Printf("[REPO] Successfully created conversation with ID: %d", conversation.ID)
	return conversation, err
}

func (r *appRepository) GetAllConversations() ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.Conn.Order("last_message_at DESC").Find(&conversations).Error
	log.Printf("[REPO] GetAllConversations found %d conversations", len(conversations))
	for _, conv := range conversations {
		log.Printf("[REPO] Conversation: ID=%d, CustomerID=%d, AgentID=%d, Status=%s",
			conv.ID, conv.CostumerID, conv.AgentID, conv.Status)
	}
	return conversations, err
}

func (r *appRepository) GetUserConversations(userID int64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.Conn.Where("costumer_id = ? OR agent_id = ?", userID, userID).
		Order("last_message_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func (r *appRepository) AssignConversationToAgent(conversationID int64, agentID int64) error {
	return r.Conn.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("agent_id", agentID).Error
}

func (r *appRepository) GetChatMessagesByConversationID(conversationID string) ([]models.ChatMessage, error) {
	id, err := strconv.ParseInt(conversationID, 10, 64)
	if err != nil {
		return nil, err
	}
	return r.GetChatMessages(id)
}

func (r *appRepository) FindActiveConversationForCustomer(customerID int64) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.Conn.Where("costumer_id = ? AND status = ?", customerID, "open").
		Order("last_message_at DESC").
		First(&conversation).Error

	if err != nil {
		log.Printf("[REPO] FindActiveConversationForCustomer: No active conversation found for customer %d: %v", customerID, err)
		return nil, err
	}

	log.Printf("[REPO] FindActiveConversationForCustomer: Found conversation %d for customer %d", conversation.ID, customerID)
	return &conversation, nil
}

func (r *appRepository) UpdateConversationLastMessage(conversationID int64) error {
	return r.Conn.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_at", time.Now()).Error
}
