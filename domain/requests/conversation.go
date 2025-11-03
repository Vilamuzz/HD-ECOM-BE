package requests

type CreateConversationRequest struct {
	AgentID *int64 `json:"agent_id" example:"123"`
}
