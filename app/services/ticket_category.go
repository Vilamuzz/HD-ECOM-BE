package services

import "app/domain/models"

func (s *appService) CreateTicketCategory(category *models.TicketCategory) error {
	return s.repo.CreateTicketCategory(category)
}

func (s *appService) GetTicketCategories() ([]models.TicketCategory, error) {
	return s.repo.GetTicketCategories()
}

func (s *appService) GetTicketCategoryByID(id int) (*models.TicketCategory, error) {
	return s.repo.GetTicketCategoryByID(id)
}

func (s *appService) UpdateTicketCategory(category *models.TicketCategory) error {
	return s.repo.UpdateTicketCategory(category)
}

func (s *appService) DeleteTicketCategory(id int) error {
	return s.repo.DeleteTicketCategory(id)
}
