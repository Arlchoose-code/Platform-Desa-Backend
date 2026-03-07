// Platform Desa — Model Pengumuman & Agenda
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type Pengumuman struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	Judul          string     `gorm:"size:255;not null" json:"judul"`
	Isi            string     `gorm:"type:longtext;not null" json:"isi"`
	FileID         *uint      `json:"file_id"`
	PenulisID      *uint      `json:"penulis_id"`
	IsPenting      bool       `gorm:"default:false" json:"is_penting"`
	Status         string     `gorm:"size:20;default:draft" json:"status"`
	TanggalMulai   *time.Time `json:"tanggal_mulai"`
	TanggalSelesai *time.Time `json:"tanggal_selesai"`
	PublishedAt    *time.Time `json:"published_at"`
	IsDeleted      bool       `gorm:"default:false" json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	File           *Media     `gorm:"foreignKey:FileID" json:"file,omitempty"`
	Penulis        *User      `gorm:"foreignKey:PenulisID" json:"penulis,omitempty"`
}

type Agenda struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	Judul          string     `gorm:"size:255;not null" json:"judul"`
	Deskripsi      *string    `gorm:"type:text" json:"deskripsi"`
	Lokasi         *string    `gorm:"size:255" json:"lokasi"`
	TanggalMulai   time.Time  `gorm:"not null" json:"tanggal_mulai"`
	TanggalSelesai *time.Time `json:"tanggal_selesai"`
	Penyelenggara  *string    `gorm:"size:100" json:"penyelenggara"`
	Status         string     `gorm:"size:20;default:upcoming" json:"status"`
	IsPublished    bool       `gorm:"default:true" json:"is_published"`
	IsDeleted      bool       `gorm:"default:false" json:"-"`
	CreatedBy      *uint      `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
