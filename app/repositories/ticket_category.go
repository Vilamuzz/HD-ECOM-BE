package repositories

import "app/domain/models"

func (r *appRepository) CreateTicketCategory(category *models.TicketCategory) error {
	return r.Conn.Create(category).Error
}

func (r *appRepository) GetTicketCategories() ([]models.TicketCategory, error) {
	var categories []models.TicketCategory
	err := r.Conn.Find(&categories).Error
	return categories, err
}

func (r *appRepository) GetTicketCategoryByID(id int) (*models.TicketCategory, error) {
	var category models.TicketCategory
	err := r.Conn.First(&category, id).Error
	return &category, err
}

func (r *appRepository) UpdateTicketCategory(category *models.TicketCategory) error {
	return r.Conn.Save(category).Error
}

func (r *appRepository) DeleteTicketCategory(id int) error {
	return r.Conn.Delete(&models.TicketCategory{}, id).Error
}
