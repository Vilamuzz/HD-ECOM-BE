package domain

import (
	"app/domain/models"
	"context"
	"mime/multipart"
	"time"
)

type AppRepository interface {
	// User operations
	GetUserByID(id uint64) (*models.User, error)
	CreateUser(user *models.User) error
	GetUsersByRole(role models.UserRole) ([]models.User, error)

	// Conversation operations
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversationByID(conversationID uint64) (*models.Conversation, error)
	GetAdminConversations(adminID uint64) ([]models.Conversation, error)
	GetCustomerConversations(userID uint64) ([]models.Conversation, error)
	UpdateConversationLastMessage(conversationID uint64) error
	CloseConversation(ctx context.Context, conversationID uint64) error
	ReopenConversation(ctx context.Context, conversationID uint64) error

	// Chat message operations
	GetMessageHistory(conversationID uint64, limit int, cursor string) ([]models.Message, string, error)
	GetMessageHistoryForAdmin(conversationID uint64, limit int, cursor string) ([]models.Message, string, error)
	SaveMessage(message *models.Message) (*models.Message, error)
	SoftDeleteConversationMessages(conversationID uint64, purgeAfterDays int) error
	ResetPurgeTimestamp(conversationID uint64) error
	PermanentlyDeleteExpiredMessages() error

	// Admin availability operations
	GetAdminAvailabilityByAdminID() (*models.AdminAvailability, error)
	CreateAdminAvailability(adminAvailability *models.AdminAvailability) error
	IncrementAdminConversationCount(adminID uint64) error
	DecrementAdminConversationCount(adminID uint64) error

	// Admin conversation state operations
	CreateAdminConversationState(adminID uint64, conversationID uint64) error
	GetAdminConversationState(adminID uint64, conversationID uint64) (*models.AdminConversationState, error)
	GetAdminConversationStatesByAdminID(adminID uint64) ([]models.AdminConversationState, error)
	IncrementUnreadCount(state *models.AdminConversationState) error
	ResetState(state *models.AdminConversationState, lastMessageID uint64) error

	// Ticket notifications
	GetOpenTicketCountsByType() (customerCount int, sellerCount int, err error)
	GetTicketStatistics() (total int, inProgress int, resolved int, priorityCounts map[int]int, err error)

	// Ticket Category
	CreateTicketCategory(category *models.TicketCategory) error
	GetTicketCategories() ([]models.TicketCategory, error)
	GetTicketCategoryByID(id int) (*models.TicketCategory, error)
	UpdateTicketCategory(category *models.TicketCategory) error
	DeleteTicketCategory(id int) error
	GetTicketsPaginated(limit, offset int) ([]models.Ticket, int, error)

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
	GetTicketsCursor(limit int, cursor string, tipePengaduan string, statusID, priorityID, categoryID int) ([]models.Ticket, string, error)

	// Ticket Assignment
	CreateTicketAssignment(assignment *models.TicketAssignment) error
	GetTicketAssignments() ([]models.TicketAssignment, error)
	GetTicketAssignmentByID(id int) (*models.TicketAssignment, error)
	GetTicketAssignmentByTicketID(ticketID int) (*models.TicketAssignment, error)
	UpdateTicketAssignment(assignment *models.TicketAssignment) error
	DeleteTicketAssignment(id int) error
	GetTicketAssignmentsByAdminIDCursor(adminID int, limit int, cursor string, statusName string) ([]models.TicketAssignment, string, error)
	GetAssignedTicketCountByAdminID(adminID int) (int, error)
	GetAssignedTicketCountByAdminIDAndStatus(adminID int, statusID int) (int, error)

	// Ticket Attachment
	CreateTicketAttachment(attachment *models.TicketAttachment) error
	GetTicketAttachments() ([]models.TicketAttachment, error)
	GetTicketAttachmentByID(id int) (*models.TicketAttachment, error)
	GetTicketAttachmentsByTicketID(ticketID int) ([]models.TicketAttachment, error)
	UpdateTicketAttachment(attachment *models.TicketAttachment) error
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

type S3Repository interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
	GetFileURL(ctx context.Context, filePath string, expiry time.Duration) (string, error)
}
