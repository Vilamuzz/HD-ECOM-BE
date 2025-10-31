package repositories

import (
	"app/domain/models"
	"log"
	"strconv"
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
		participants = append(participants, conversation.AgentID)
	}

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
