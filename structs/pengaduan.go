// Platform Desa — Structs Pengaduan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type PengaduanRequest struct {
	NamaPelapor   string  `json:"nama_pelapor" binding:"required,min=3"`
	KontakPelapor *string `json:"kontak_pelapor"`
	IsAnonim      bool    `json:"is_anonim"`
	Kategori      string  `json:"kategori" binding:"required,oneof=infrastruktur pelayanan keamanan lingkungan lainnya"`
	Judul         string  `json:"judul" binding:"required,min=5"`
	Isi           string  `json:"isi" binding:"required,min=10"`
	Lokasi        *string `json:"lokasi"`
	FotoID        *uint   `json:"foto_id"`
}

type ResponPengaduanRequest struct {
	Respon string `json:"respon" binding:"required"`
}

type TolakPengaduanRequest struct {
	Alasan string `json:"alasan" binding:"required"`
}
