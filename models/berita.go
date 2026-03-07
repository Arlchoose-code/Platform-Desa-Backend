// Platform Desa — Model Berita
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type KategoriBerita struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Nama      string `gorm:"size:100;not null" json:"nama"`
	Slug      string `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Warna     string `gorm:"size:7;default:#3B82F6" json:"warna"`
	IsDeleted bool   `gorm:"default:false" json:"-"`
}

type Berita struct {
	ID           uint            `gorm:"primarykey" json:"id"`
	KategoriID   *uint           `json:"kategori_id"`
	KategoriNama *string         `gorm:"size:100" json:"kategori_nama"`
	Judul        string          `gorm:"size:255;not null" json:"judul"`
	Slug         string          `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Ringkasan    *string         `gorm:"type:text" json:"ringkasan"`
	Isi          string          `gorm:"type:longtext;not null" json:"isi"`
	ThumbnailID  *uint           `json:"thumbnail_id"`
	PenulisID    *uint           `json:"penulis_id"`
	Status       string          `gorm:"size:20;default:draft" json:"status"`
	Views        uint            `gorm:"default:0" json:"views"`
	PublishedAt  *time.Time      `json:"published_at"`
	IsDeleted    bool            `gorm:"default:false" json:"-"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Kategori     *KategoriBerita `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Thumbnail    *Media          `gorm:"foreignKey:ThumbnailID" json:"thumbnail,omitempty"`
	Penulis      *User           `gorm:"foreignKey:PenulisID" json:"penulis,omitempty"`
}
