package requests

import "time"

type CreateTicketCommentRequest struct {
	TicketID int    `json:"ticket_id" binding:"required"`
	IsiPesan string `json:"isi_pesan" binding:"required"`
}

type UpdateTicketCommentRequest struct {
	TicketID int    `json:"ticket_id" binding:"required"`
	IsiPesan string `json:"isi_pesan" binding:"required"`
}

type TicketCommentResponse struct {
	CommentID     int       `json:"comment_id"`
	TicketID      int       `json:"ticket_id"`
	UserID        int       `json:"user_id"`
	IsiPesan      string    `json:"isi_pesan"`
	TanggalDibuat time.Time `json:"tanggal_dibuat"`
}
