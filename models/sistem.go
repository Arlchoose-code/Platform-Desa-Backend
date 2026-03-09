// Platform Desa — Model Peta & FAQ
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type PetaFasilitas struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Nama        string    `gorm:"size:150;not null" json:"nama"`
	Kategori    string    `gorm:"size:20;not null" json:"kategori"` // pendidikan, kesehatan, ibadah, pemerintahan, olahraga, lainnya
	Alamat      *string   `gorm:"type:text" json:"alamat"`
	Latitude    float64   `gorm:"not null" json:"latitude"`
	Longitude   float64   `gorm:"not null" json:"longitude"`
	Deskripsi   *string   `gorm:"type:text" json:"deskripsi"`
	FotoID      *uint     `json:"foto_id"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Foto        *Media    `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
}

type FAQ struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Pertanyaan  string    `gorm:"type:text;not null" json:"pertanyaan"`
	Jawaban     string    `gorm:"type:longtext;not null" json:"jawaban"`
	Kategori    *string   `gorm:"size:100" json:"kategori"`
	Urutan      uint      `gorm:"default:0" json:"urutan"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	IsDeleted   bool      `gorm:"default:false" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
