package requests

import "time"

// CreateTicketCommentRequest is used for creating a new ticket comment
type CreateTicketCommentRequest struct {
	IDTicket      int       `json:"id_ticket"`
	IDUser        int       `json:"id_user"`
	IsiPesan      string    `json:"isi_pesan"`
	TanggalDibuat time.Time `json:"tanggal_dibuat"`
}

// UpdateTicketCommentRequest is used for updating a ticket comment
// All fields are required for update
// IDComment is taken from path param
// TanggalDibuat can be updated if needed
// If not, you can omit it from the request struct
// For now, keep it for consistency
// You can adjust as needed

// UpdateTicketCommentRequest is used for updating a ticket comment
// All fields except IDComment
// IDComment is taken from path param
// TanggalDibuat can be updated if needed
// If not, you can omit it from the request struct
// For now, keep it for consistency
// You can adjust as needed

// UpdateTicketCommentRequest is used for updating a ticket comment
// All fields except IDComment
// IDComment is taken from path param
// TanggalDibuat can be updated if needed
// If not, you can omit it from the request struct
// For now, keep it for consistency
// You can adjust as needed

type UpdateTicketCommentRequest struct {
	IDTicket      int       `json:"id_ticket"`
	IDUser        int       `json:"id_user"`
	IsiPesan      string    `json:"isi_pesan"`
	TanggalDibuat time.Time `json:"tanggal_dibuat"`
}

// TicketCommentResponse is used for returning a ticket comment in responses
// Only main fields, no relations

type TicketCommentResponse struct {
	IDComment     int       `json:"id_comment"`
	IDTicket      int       `json:"id_ticket"`
	IDUser        int       `json:"id_user"`
	IsiPesan      string    `json:"isi_pesan"`
	TanggalDibuat time.Time `json:"tanggal_dibuat"`
}
