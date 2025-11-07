package domain

import "app/domain/models"

func GetAllModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Conversation{},
		&models.ChatMessage{},
		&models.AdminAvailability{},
	}
}
