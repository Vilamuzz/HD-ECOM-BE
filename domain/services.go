package domain

import "app/domain/models"

type AppService interface {
	// Ticket Category
	CreateTicketCategory(category *models.TicketCategory) error
	GetTicketCategories() ([]models.TicketCategory, error)
	GetTicketCategoryByID(id int) (*models.TicketCategory, error)
	UpdateTicketCategory(category *models.TicketCategory) error
	DeleteTicketCategory(id int) error
}
