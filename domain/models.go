package domain

import "app/domain/models"

func GetAllModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Conversation{},
		&models.ChatMessage{},
		// Master tables first
		&models.TicketCategory{},
		&models.TicketPriority{},
		&models.TicketStatus{},
		// Then transaction tables
		&models.Ticket{},
		&models.TicketComment{},
		&models.TicketAttachment{},
		&models.TicketAssignment{},
		&models.TicketLog{},
	}
}
