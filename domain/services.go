package domain

import (
	"app/domain/models"
	"app/helpers"
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type AppService interface {
	// User management
	GetSupportUsers() helpers.Response

	// WebSocket management
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *Message)
	SendToRecipients(message *Message)
	JoinConversation(client *Client, conversationID uint64)
	ServeWebSocket(ctx *gin.Context)

	// Conversation management
	GetConversations(claim models.User) helpers.Response
	GetConversationByID(conversationID uint64) (*models.Conversation, error)
	CreateCustomerConversation(ctx context.Context, claim models.User) helpers.Response
	CloseConversation(ctx context.Context, claim models.User, id string) helpers.Response
	ReopenConversation(ctx context.Context, conversationID uint64) error

	// Message management
	GetMessageHistory(conversationID uint64, limit int, cursor string, isAdmin bool) helpers.Response

	// Admin state management
	GetAdminConversationState(adminID uint64, conversationID uint64) (*models.AdminConversationState, error)
	GetAdminListConversationStates(claim models.User) helpers.Response
	// Ticket Category
	CreateTicketCategory(category *models.TicketCategory) error
	GetTicketCategories() ([]models.TicketCategory, error)
	GetTicketCategoryByID(id int) (*models.TicketCategory, error)
	UpdateTicketCategory(category *models.TicketCategory) error
	DeleteTicketCategory(id int) error

	// Ticket Priority
	CreateTicketPriority(priority *models.TicketPriority) error
	GetTicketPriorities() ([]models.TicketPriority, error)
	GetTicketPriorityByID(id int) (*models.TicketPriority, error)
	UpdateTicketPriority(priority *models.TicketPriority) error
	DeleteTicketPriority(id int) error

	// Ticket Status
	CreateTicketStatus(status *models.TicketStatus) error
	GetTicketStatuses() ([]models.TicketStatus, error)
	GetTicketStatusByID(id int) (*models.TicketStatus, error)
	UpdateTicketStatus(status *models.TicketStatus) error
	DeleteTicketStatus(id int) error

	// Ticket
	CreateTicket(ticket *models.Ticket) error
	GetTickets() ([]models.Ticket, error)
	GetTicketByID(id int) (*models.Ticket, error)
	GetTicketsByUserID(userID int) ([]models.Ticket, error)
	UpdateTicket(ticket *models.Ticket) error
	DeleteTicket(id int) error

	// Ticket Assignment
	CreateTicketAssignment(assignment *models.TicketAssignment) error
	GetTicketAssignments() ([]models.TicketAssignment, error)
	GetTicketAssignmentByID(id int) (*models.TicketAssignment, error)
	UpdateTicketAssignment(assignment *models.TicketAssignment) error
	DeleteTicketAssignment(id int) error

	// Ticket Attachment
	CreateTicketAttachment(ticketID int, file *multipart.FileHeader) (*models.TicketAttachment, error)
	GetTicketAttachments() ([]models.TicketAttachment, error)
	GetTicketAttachmentByID(id int) (*models.TicketAttachment, string, error)
	GetTicketAttachmentsByTicketID(ticketID int) ([]models.TicketAttachment, error)
	UpdateTicketAttachment(id int, ticketID *int, file *multipart.FileHeader) (*models.TicketAttachment, error)
	DeleteTicketAttachment(id int) error

	// Ticket Comment
	CreateTicketComment(comment *models.TicketComment) error
	GetTicketComments() ([]models.TicketComment, error)
	GetTicketCommentByID(id int) (*models.TicketComment, error)
	GetTicketCommentsByTicketID(ticketID int) ([]models.TicketComment, error)
	UpdateTicketComment(comment *models.TicketComment) error
	DeleteTicketComment(id int) error

	// Ticket Log
	CreateTicketLog(log *models.TicketLog) error
	GetTicketLogs() ([]models.TicketLog, error)
	GetTicketLogByID(id int) (*models.TicketLog, error)
	GetTicketLogsByTicketID(ticketID int) ([]models.TicketLog, error)
}
