// Platform Desa — Model Wisata
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type KategoriWisata struct {
	ID        uint    `gorm:"primarykey" json:"id"`
	Nama      string  `gorm:"size:100;not null" json:"nama"`
	Slug      string  `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Icon      *string `gorm:"size:50" json:"icon"`
	IsDeleted bool    `gorm:"default:false" json:"-"`
}

type Wisata struct {
	ID           uint            `gorm:"primarykey" json:"id"`
	KategoriID   *uint           `json:"kategori_id"`
	KategoriNama *string         `gorm:"size:100" json:"kategori_nama"`
	Nama         string          `gorm:"size:150;not null" json:"nama"`
	Slug         string          `gorm:"size:150;not null;uniqueIndex" json:"slug"`
	Deskripsi    *string         `gorm:"type:longtext" json:"deskripsi"`
	Alamat       *string         `gorm:"type:text" json:"alamat"`
	Latitude     *float64        `json:"latitude"`
	Longitude    *float64        `json:"longitude"`
	JamBuka      *string         `gorm:"size:8" json:"jam_buka"`
	JamTutup     *string         `gorm:"size:8" json:"jam_tutup"`
	HargaTiket   int64           `gorm:"default:0" json:"harga_tiket"`
	Kontak       *string         `gorm:"size:20" json:"kontak"`
	ThumbnailID  *uint           `json:"thumbnail_id"`
	Fasilitas    *string         `gorm:"type:json" json:"fasilitas"`
	IsPublished  bool            `gorm:"default:true" json:"is_published"`
	IsDeleted    bool            `gorm:"default:false" json:"-"`
	Views        uint            `gorm:"default:0" json:"views"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Kategori     *KategoriWisata `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Thumbnail    *Media          `gorm:"foreignKey:ThumbnailID" json:"thumbnail,omitempty"`
	Galeri       []WisataGaleri  `gorm:"foreignKey:WisataID" json:"galeri,omitempty"`
}

type WisataGaleri struct {
	ID       uint    `gorm:"primarykey" json:"id"`
	WisataID uint    `gorm:"not null" json:"wisata_id"`
	MediaID  uint    `gorm:"not null" json:"media_id"`
	Caption  *string `gorm:"size:255" json:"caption"`
	Urutan   uint    `gorm:"default:0" json:"urutan"`
	Media    *Media  `gorm:"foreignKey:MediaID" json:"media,omitempty"`
}
