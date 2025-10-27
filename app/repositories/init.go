package repositories

import (
	"app/domain"

	"gorm.io/gorm"
)

type appRepository struct {
	Conn *gorm.DB
}

func NewAppRepository(conn *gorm.DB) domain.AppRepository {
	return &appRepository{
		Conn: conn,
	}
}
