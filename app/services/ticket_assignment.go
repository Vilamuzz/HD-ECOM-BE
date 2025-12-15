package services

import (
	"app/domain/models"
	"errors"
	"fmt"
	"time"
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

    // Get the ticket first to ensure it exists
    ticket, err := s.repo.GetTicketByID(assignment.TicketID)
    if err != nil {
        return fmt.Errorf("ticket not found: %v", err)
    }

    // Update ticket status to "In Progress" (ID 2) BEFORE creating assignment
    ticket.StatusID = 2
    ticket.TanggalDiperbarui = time.Now()
    
    if err := s.repo.UpdateTicket(ticket); err != nil {
        return fmt.Errorf("failed to update ticket status: %v", err)
    }

    // Create the assignment after ticket is updated
    if err := s.repo.CreateTicketAssignment(assignment); err != nil {
        return fmt.Errorf("failed to create ticket assignment: %v", err)
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
    // Validate that the admin being assigned has the support role
    admin, err := s.repo.GetUserByID(uint64(assignment.AdminID))
    if err != nil {
        return errors.New("admin not found")
    }

    if admin.Role != models.RoleSupport {
        return errors.New("only users with support role can be assigned to tickets")
    }

    // Get the ticket first to ensure it exists
    ticket, err := s.repo.GetTicketByID(assignment.TicketID)
    if err != nil {
        return fmt.Errorf("ticket not found: %v", err)
    }

    // Update ticket status to "In Progress" (ID 2) BEFORE updating assignment
    ticket.StatusID = 2
    ticket.TanggalDiperbarui = time.Now()
    
    if err := s.repo.UpdateTicket(ticket); err != nil {
        return fmt.Errorf("failed to update ticket status: %v", err)
    }

    // Update the assignment after ticket is updated
    if err := s.repo.UpdateTicketAssignment(assignment); err != nil {
        return fmt.Errorf("failed to update ticket assignment: %v", err)
    }

    return nil
}

func (s *appService) DeleteTicketAssignment(id int) error {
    return s.repo.DeleteTicketAssignment(id)
}

func (s *appService) GetTicketAssignmentByTicketID(ticketID int) (*models.TicketAssignment, error) {
    return s.repo.GetTicketAssignmentByTicketID(ticketID)
}

// New method for cursor-based pagination and status filtering
func (s *appService) GetTicketAssignmentsByAdminIDCursor(adminID int, limit int, cursor string, statusName string) ([]models.TicketAssignment, string, error) {
    return s.repo.GetTicketAssignmentsByAdminIDCursor(adminID, limit, cursor, statusName)
}

// New method to get total assigned ticket count by admin ID
func (s *appService) GetAssignedTicketCountByAdminID(adminID int) (int, error) {
    return s.repo.GetAssignedTicketCountByAdminID(adminID)
}

// New method to get assigned ticket count by admin ID and status ID
func (s *appService) GetAssignedTicketCountByAdminIDAndStatus(adminID int, statusID int) (int, error) {
    return s.repo.GetAssignedTicketCountByAdminIDAndStatus(adminID, statusID)
}
