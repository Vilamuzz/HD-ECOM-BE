package requests

type CreateTicketAssignmentRequest struct {
	TicketID          int    `json:"id_ticket" example:"1" description:"ID of the ticket to assign"`
	AdminID           int    `json:"id_admin" example:"1" description:"ID of the admin user who will handle the ticket"`
	PriorityID        int    `json:"id_priority" example:"3" description:"Priority level: 1=Low, 2=Medium, 3=High, 4=Critical"`
	TanggalDitugaskan string `json:"tanggal_ditugaskan" example:"2025-12-18T10:00:00Z" description:"Assignment date in ISO format"`
}

// TicketAssignmentResponse is a clean response DTO for TicketAssignment
// @Description TicketAssignmentResponse represents a ticket assignment with related ticket and admin info (no null fields)
type TicketAssignmentResponse struct {
	AssignmentID      int                 `json:"id_assignment"`
	TicketID          int                 `json:"id_ticket"`
	AdminID           int                 `json:"id_admin"`
	PriorityID        int                 `json:"id_priority"`
	TanggalDitugaskan string              `json:"tanggal_ditugaskan"`
	Ticket            *TicketResponse     `json:"ticket,omitempty"`
	Admin             *UserSimpleResponse `json:"admin,omitempty"`
	Priority          *PriorityResponse   `json:"priority,omitempty"`
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

// PriorityResponse is a minimal priority info for assignment
// @Description PriorityResponse represents priority info for assignment
type PriorityResponse struct {
	IDPriority   int    `json:"id_priority"`
	NamaPriority string `json:"nama_priority"`
}
