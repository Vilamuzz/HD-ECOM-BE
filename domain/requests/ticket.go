package requests

type TicketCreateRequest struct {
	KodeTiket  string `json:"kode_tiket" example:"TCKT-001"`
	UserID     uint64 `json:"id_user" example:"1"`
	Judul      string `json:"judul" example:"Tidak bisa login"`
	Deskripsi  string `json:"deskripsi" example:"Saya tidak bisa login ke akun saya."`
	CategoryID int    `json:"id_category" example:"1"`
	PriorityID int    `json:"id_priority" example:"1"`
	StatusID   int    `json:"id_status" example:"1"`
	// Allowed: pelanggan, penjual
	TipePengaduan string `json:"tipe_pengaduan" example:"pelanggan" enums:"pelanggan,penjual"`
}

type TicketResponse struct {
	ID                int    `json:"id_ticket"`
	KodeTiket         string `json:"kode_tiket"`
	UserID            uint64 `json:"id_user"`
	Judul             string `json:"judul"`
	Deskripsi         string `json:"deskripsi"`
	CategoryID        int    `json:"id_category"`
	PriorityID        int    `json:"id_priority"`
	StatusID          int    `json:"id_status"`
	TipePengaduan     string `json:"tipe_pengaduan"`
	TanggalDibuat     string `json:"tanggal_dibuat"`
	TanggalDiperbarui string `json:"tanggal_diperbarui"`
}
