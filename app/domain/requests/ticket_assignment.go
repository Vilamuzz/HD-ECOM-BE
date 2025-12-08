package requests

// TicketAssignmentResponse represents the response structure for a ticket assignment
type TicketAssignmentResponse struct {
	AssignmentID      int             `json:"id_assignment"`
	TicketID          int             `json:"id_ticket"`
	AdminID           int             `json:"id_admin"`
	TanggalDitugaskan string          `json:"tanggal_ditugaskan"`
	Ticket            *TicketResponse `json:"ticket,omitempty"`
}

// TicketResponse represents the response structure for a ticket
type TicketResponse struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}