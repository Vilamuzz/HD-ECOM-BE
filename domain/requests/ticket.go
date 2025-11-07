package requests

type TicketCreateRequest struct {
	KodeTiket      string `json:"kode_tiket" example:"TCKT-001"`
	IDUser         int    `json:"id_user" example:"1"`
	Judul          string `json:"judul" example:"Tidak bisa login"`
	Deskripsi      string `json:"deskripsi" example:"Saya tidak bisa login ke akun saya."`
	IDCategory     int    `json:"id_category" example:"1"`
	IDPriority     int    `json:"id_priority" example:"1"`
	IDStatus       int    `json:"id_status" example:"1"`
	// Allowed: pelanggan, penjual
	TipePengaduan  string `json:"tipe_pengaduan" example:"pelanggan" enums:"pelanggan,penjual"`
}

type TicketResponse struct {
	IDTicket        int    `json:"id_ticket"`
	KodeTiket       string `json:"kode_tiket"`
	IDUser          int    `json:"id_user"`
	Judul           string `json:"judul"`
	Deskripsi       string `json:"deskripsi"`
	IDCategory      int    `json:"id_category"`
	IDPriority      int    `json:"id_priority"`
	IDStatus        int    `json:"id_status"`
	TipePengaduan   string `json:"tipe_pengaduan"`
	TanggalDibuat   string `json:"tanggal_dibuat"`
	TanggalDiperbarui string `json:"tanggal_diperbarui"`
}
