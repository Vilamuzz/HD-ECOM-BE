package domain

import "app/domain/models"

type AppRepository interface {
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
}
