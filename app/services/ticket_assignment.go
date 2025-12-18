package services

import "app/domain/models"

func (s *appService) CreateTicketAssignment(assignment *models.TicketAssignment) error {
	return s.repo.CreateTicketAssignment(assignment)
}

func (s *appService) GetTicketAssignments() ([]models.TicketAssignment, error) {
	return s.repo.GetTicketAssignments()
}

func (s *appService) GetTicketAssignmentByID(id int) (*models.TicketAssignment, error) {
	return s.repo.GetTicketAssignmentByID(id)
}

func (s *appService) UpdateTicketAssignment(assignment *models.TicketAssignment) error {
	return s.repo.UpdateTicketAssignment(assignment)
}

func (s *appService) DeleteTicketAssignment(id int) error {
	return s.repo.DeleteTicketAssignment(id)
}
