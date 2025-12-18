package services

import "app/domain/models"

func (s *appService) CreateTicketAssignment(assignment *models.TicketAssignment) error {
	// Create the assignment
	if err := s.repo.CreateTicketAssignment(assignment); err != nil {
		return err
	}

	// Update ticket priority if priority is specified
	if assignment.PriorityID > 0 {
		if err := s.repo.UpdateTicketPriorityByID(assignment.TicketID, assignment.PriorityID); err != nil {
			return err
		}
	}

	return nil
}

func (s *appService) GetTicketAssignments() ([]models.TicketAssignment, error) {
	return s.repo.GetTicketAssignments()
}

func (s *appService) GetTicketAssignmentByID(id int) (*models.TicketAssignment, error) {
	return s.repo.GetTicketAssignmentByID(id)
}

func (s *appService) UpdateTicketAssignment(assignment *models.TicketAssignment) error {
	// Update the assignment
	if err := s.repo.UpdateTicketAssignment(assignment); err != nil {
		return err
	}

	// Update ticket priority if priority is specified
	if assignment.PriorityID > 0 {
		if err := s.repo.UpdateTicketPriorityByID(assignment.TicketID, assignment.PriorityID); err != nil {
			return err
		}
	}

	return nil
}

func (s *appService) DeleteTicketAssignment(id int) error {
	return s.repo.DeleteTicketAssignment(id)
}
