// Platform Desa — Model Galeri
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type KategoriGaleri struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Nama      string `gorm:"size:100;not null" json:"nama"`
	Slug      string `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Tipe      string `gorm:"size:10;default:semua" json:"tipe"`
	IsDeleted bool   `gorm:"default:false" json:"-"`
}

type Galeri struct {
	ID           uint            `gorm:"primarykey" json:"id"`
	KategoriID   *uint           `json:"kategori_id"`
	KategoriNama *string         `gorm:"size:100" json:"kategori_nama"`
	Judul        string          `gorm:"size:255;not null" json:"judul"`
	Deskripsi    *string         `gorm:"type:text" json:"deskripsi"`
	MediaID      uint            `gorm:"not null" json:"media_id"`
	ThumbnailID  *uint           `json:"thumbnail_id"`
	Tanggal      *time.Time      `json:"tanggal"`
	Fotografer   *string         `gorm:"size:100" json:"fotografer"`
	IsPublished  bool            `gorm:"default:true" json:"is_published"`
	IsDeleted    bool            `gorm:"default:false" json:"-"`
	Urutan       uint            `gorm:"default:0" json:"urutan"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Kategori     *KategoriGaleri `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Media        *Media          `gorm:"foreignKey:MediaID" json:"media,omitempty"`
	Thumbnail    *Media          `gorm:"foreignKey:ThumbnailID" json:"thumbnail,omitempty"`
}
