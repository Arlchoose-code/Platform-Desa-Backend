// Platform Desa — Structs Galeri
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

import "time"

type GaleriRequest struct {
	KategoriID  *uint      `json:"kategori_id"`
	Judul       string     `json:"judul" binding:"required,min=3"`
	Deskripsi   *string    `json:"deskripsi"`
	MediaID     uint       `json:"media_id" binding:"required"`
	ThumbnailID *uint      `json:"thumbnail_id"`
	Tanggal     *time.Time `json:"tanggal"`
	Fotografer  *string    `json:"fotografer"`
	IsPublished bool       `json:"is_published"`
	Urutan      uint       `json:"urutan"`
}

type KategoriGaleriRequest struct {
	Nama string `json:"nama" binding:"required,min=2"`
	Tipe string `json:"tipe" binding:"required,oneof=foto video semua"`
}
