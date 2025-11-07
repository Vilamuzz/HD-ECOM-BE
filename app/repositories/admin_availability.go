package repositories

import "app/domain/models"

func (r *appRepository) CreateAdminAvailability(adminAvailability *models.AdminAvailability) error {
	return r.Conn.Create(adminAvailability).Error
}

func (r *appRepository) UpdateAdminAvailability(adminAvailability *models.AdminAvailability) error {
	return r.Conn.Save(adminAvailability).Error
}
