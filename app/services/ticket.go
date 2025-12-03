package services

import (
	"app/domain/models"
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

func (s *appService) CreateTicket(ticket *models.Ticket) error {
	// Generate kode_tiket if not provided
	if ticket.KodeTiket == "" {
		ticket.KodeTiket = s.generateTicketCode(ticket.Judul, ticket.Deskripsi)
	}
	return s.repo.CreateTicket(ticket)
}

func (s *appService) generateTicketCode(title, description string) string {
	// Combine title + description + current time
	input := fmt.Sprintf("%s%s%d", title, description, time.Now().UnixNano())

	// Generate MD5 hash
	hash := md5.Sum([]byte(input))
	hashString := fmt.Sprintf("%x", hash)

	// Take first 8 characters and convert to uppercase
	shortHash := strings.ToUpper(hashString[:8])

	return fmt.Sprintf("TCKT-%s", shortHash)
}

func (s *appService) GetTickets() ([]models.Ticket, error) {
	return s.repo.GetTickets()
}

func (s *appService) GetTicketByID(id int) (*models.Ticket, error) {
	return s.repo.GetTicketByID(id)
}

func (s *appService) UpdateTicket(ticket *models.Ticket) error {
	return s.repo.UpdateTicket(ticket)
}

func (s *appService) DeleteTicket(id int) error {
	return s.repo.DeleteTicket(id)
}

func (s *appService) GetTicketsByUserID(userID int) ([]models.Ticket, error) {
	return s.repo.GetTicketsByUserID(userID)
}
