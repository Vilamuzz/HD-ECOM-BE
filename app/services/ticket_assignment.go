package services

import (
	"app/domain/models"
	"errors"
)



func (s *appService) CreateTicketAssignment(assignment *models.TicketAssignment) error {
    // Validate that the admin being assigned has the support role
    admin, err := s.repo.GetUserByID(uint64(assignment.AdminID))
    if err != nil {
        return errors.New("admin not found")
    }

    if admin.Role != models.RoleSupport {
        return errors.New("only users with support role can be assigned to tickets")
    }

    return s.repo.CreateTicketAssignment(assignment)
}

func (s *appService) GetTicketAssignments() ([]models.TicketAssignment, error) {
	return s.repo.GetTicketAssignments()
}

func (s *appService) GetTicketAssignmentByID(id int) (*models.TicketAssignment, error) {
	return s.repo.GetTicketAssignmentByID(id)
}

func (s *appService) UpdateTicketAssignment(assignment *models.TicketAssignment) error {
    // Validate that the admin being assigned has the support role
    admin, err := s.repo.GetUserByID(uint64(assignment.AdminID))
    if err != nil {
        return errors.New("admin not found")
    }

    if admin.Role != models.RoleSupport {
        return errors.New("only users with support role can be assigned to tickets")
    }

    return s.repo.UpdateTicketAssignment(assignment)
}

func (s *appService) DeleteTicketAssignment(id int) error {
	return s.repo.DeleteTicketAssignment(id)
}
