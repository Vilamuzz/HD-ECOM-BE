package requests

import "app/domain/models"

type TicketCreateRequest struct {
	KodeTiket  string `json:"kode_tiket" example:"TCKT-001"`
	UserID     uint64 `json:"id_user" example:"1"`
	Judul      string `json:"judul" example:"Login Issues"`
	Deskripsi  string `json:"deskripsi" example:"Cannot access my account after password reset"`
	CategoryID int    `json:"id_category" example:"1" description:"1=Technical Issue, 2=Account Problem, 3=Payment Issue"`
	PriorityID int    `json:"id_priority" example:"2" description:"1=Low, 2=Medium, 3=High, 4=Critical"`
	StatusID   int    `json:"id_status" example:"1" description:"1=Open, 2=In Progress, 3=Resolved, 4=Closed"`
	// Note: tipe_pengaduan will be auto-filled based on authenticated user's role
	TipePengaduan models.UserRole `json:"tipe_pengaduan,omitempty" swaggerignore:"true"`
}

type TicketResponse struct {
	ID                int             `json:"id_ticket"`
	KodeTiket         string          `json:"kode_tiket"`
	UserID            uint64          `json:"id_user"`
	Judul             string          `json:"judul"`
	Deskripsi         string          `json:"deskripsi"`
	CategoryID        int             `json:"id_category"`
	PriorityID        int             `json:"id_priority"`
	StatusID          int             `json:"id_status"`
	TipePengaduan     models.UserRole `json:"tipe_pengaduan"`
	TanggalDibuat     string          `json:"tanggal_dibuat"`
	TanggalDiperbarui string          `json:"tanggal_diperbarui"`
}
