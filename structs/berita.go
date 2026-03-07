// Platform Desa — Structs Berita
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type BeritaRequest struct {
	KategoriID  *uint   `json:"kategori_id"`
	Judul       string  `json:"judul" binding:"required,min=3"`
	Ringkasan   *string `json:"ringkasan"`
	Isi         string  `json:"isi" binding:"required"`
	ThumbnailID *uint   `json:"thumbnail_id"`
	Status      string  `json:"status" binding:"required,oneof=draft published archived"`
}

type KategoriBeritaRequest struct {
	Nama  string `json:"nama" binding:"required,min=2"`
	Warna string `json:"warna" binding:"required"`
}
