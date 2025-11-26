package repositories

import (
	"app/domain/models"

	"gorm.io/gorm"
)

func (r *appRepository) CreateAdminAvailability(adminAvailability *models.AdminAvailability) error {
	return r.Conn.Create(adminAvailability).Error
}

func (r *appRepository) GetAdminAvailabilityByAdminID() (*models.AdminAvailability, error) {
	var availability models.AdminAvailability
	err := r.Conn.Order("current_conversations ASC").First(&availability).Error
	if err != nil {
		return nil, err
	}
	return &availability, nil
}

func (r *appRepository) IncrementAdminConversationCount(adminID uint64) error {
	return r.Conn.Model(&models.AdminAvailability{}).
		Where("admin_id = ?", adminID).
		UpdateColumn("current_conversations", gorm.Expr("current_conversations + ?", 1)).Error
}

func (r *appRepository) DecrementAdminConversationCount(adminID uint64) error {
	return r.Conn.Model(&models.AdminAvailability{}).
		Where("admin_id = ? AND current_conversations > 0", adminID).
		UpdateColumn("current_conversations", gorm.Expr("current_conversations - ?", 1)).Error
}
