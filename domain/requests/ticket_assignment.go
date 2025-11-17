package requests

type CreateTicketAssignmentRequest struct {
	IDTicket          int    `json:"id_ticket" example:"1"`
	IDAdmin           int    `json:"id_admin" example:"2"`
	TanggalDitugaskan string `json:"tanggal_ditugaskan" example:"2025-11-07T10:00:00Z"`
}

// TicketAssignmentResponse is a clean response DTO for TicketAssignment
// @Description TicketAssignmentResponse represents a ticket assignment with related ticket and admin info (no null fields)
type TicketAssignmentResponse struct {
	IDAssignment      int                 `json:"id_assignment"`
	IDTicket          int                 `json:"id_ticket"`
	IDAdmin           int                 `json:"id_admin"`
	TanggalDitugaskan string              `json:"tanggal_ditugaskan"`
	Ticket            *TicketResponse     `json:"ticket,omitempty"`
	Admin             *UserSimpleResponse `json:"admin,omitempty"`
}

// UserSimpleResponse is a minimal user info for assignment admin
// @Description UserSimpleResponse represents a minimal user info for assignment admin
// Only the fields needed for assignment context
// (id_user, nama, email, username, role, created_at)
type UserSimpleResponse struct {
	IDUser    int    `json:"id_user"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
