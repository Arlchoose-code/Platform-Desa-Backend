// Platform Desa — Structs Regulasi
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

import "time"

type RegulasiRequest struct {
	KategoriID    *uint      `json:"kategori_id"`
	Nomor         *string    `json:"nomor"`
	Judul         string     `json:"judul" binding:"required,min=3"`
	Tentang       *string    `json:"tentang"`
	Tahun         *int       `json:"tahun"`
	TanggalTerbit *time.Time `json:"tanggal_terbit"`
	FileID        *uint      `json:"file_id"`
	IsPublished   bool       `json:"is_published"`
}

type KategoriRegulasiRequest struct {
	Nama string `json:"nama" binding:"required,min=2"`
}
