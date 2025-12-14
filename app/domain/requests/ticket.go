package requests

// ...existing code...

type TicketListMeta struct {
	NextCursor string `json:"next_cursor"`
	Limit      int    `json:"limit"`
}

type TicketListResponse struct {
	Data []TicketResponse `json:"data"`
	Meta TicketListMeta   `json:"meta"`
}

// ...existing code...