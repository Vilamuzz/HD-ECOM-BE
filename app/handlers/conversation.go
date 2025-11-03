package handlers

import (
	"app/domain/models"
	"app/domain/requests"
	"app/helpers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) ConversationRoute(path string) {
	api := r.Route.Group(path)
	api.Use(r.Middleware.Auth())
	{
		api.POST("", r.CreateConversation)
		api.GET("", r.GetConversations)
		api.GET("/:id/messages", r.GetConversationMessages)
	}
}

// CreateConversation godoc
// @Summary      Create a new conversation
// @Description  Create a new conversation with optional agent assignment
// @Security 	 BearerAuth
// @Tags         conversations
// @Accept       json
// @Produce      json
// @Param        conversation  body      requests.CreateConversationRequest  true  "Conversation Data"
// @Success      201  {object}   helpers.Response{data=models.Conversation}
// @Failure      400  {object}   helpers.Response
// @Failure      500  {object}   helpers.Response
// @Router       /conversations [post]
func (r *appRoute) CreateConversation(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)

	var req requests.CreateConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid request", nil, nil))
		return
	}

	conversation := &models.Conversation{
		CostumerID:    user.ID,
		Status:        "open",
		LastMessageAt: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if req.AgentID != nil {
		conversation.AgentID = *req.AgentID
	}

	savedConversation, err := r.Repository.CreateConversation(conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create conversation", nil, nil))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewResponse(http.StatusCreated, "Conversation created", nil, savedConversation))
}

// GetConversations godoc
// @Summary      Get conversations
// @Description  Get all conversations for the current user (admin/agent sees all, customer sees their own)
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Success      200  {object}   helpers.Response{data=[]models.Conversation}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations [get]
func (r *appRoute) GetConversations(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)

	log.Printf("[HANDLER] GetConversations - User: ID=%d, Role=%s", user.ID, user.Role)

	var conversations []models.Conversation
	var err error

	if user.Role == "admin" || user.Role == "agent" {
		// Admin/Agent sees all conversations or assigned conversations
		log.Println("[HANDLER] Fetching all conversations (admin/agent)")
		conversations, err = r.Repository.GetAllConversations()
	} else {
		// Customer sees only their conversations
		log.Printf("[HANDLER] Fetching conversations for user %d (customer)", user.ID)
		conversations, err = r.Repository.GetUserConversations(user.ID)
	}

	if err != nil {
		log.Printf("[HANDLER] Error fetching conversations: %v", err)
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to get conversations", nil, nil))
		return
	}

	log.Printf("[HANDLER] Returning %d conversations", len(conversations))
	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Success", nil, conversations))
}

// GetConversationMessages godoc
// @Summary      Get conversation messages
// @Description  Get all messages for a specific conversation
// @Security 	 BearerAuth
// @Tags         conversations
// @Produce      json
// @Param        id   path      int  true  "Conversation ID"
// @Success      200  {object}   helpers.Response{data=[]models.ChatMessage}
// @Failure      500  {object}   helpers.Response
// @Router       /conversations/{id}/messages [get]
func (r *appRoute) GetConversationMessages(c *gin.Context) {
	conversationID := c.Param("id")

	messages, err := r.Repository.GetChatMessagesByConversationID(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to get messages", nil, nil))
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Success", nil, messages))
}
