package repositories

import (
	"app/domain/models"
)

func (r *appRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	err := r.Conn.First(&user, id).Error
	return &user, err
}

func (r *appRepository) CreateUser(user *models.User) error {
	return r.Conn.Create(user).Error
}
