// Platform Desa — Model UMKM
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type KategoriUMKM struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Nama      string `gorm:"size:100;not null" json:"nama"`
	Slug      string `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	IsDeleted bool   `gorm:"default:false" json:"-"`
}

type UMKM struct {
	ID           uint          `gorm:"primarykey" json:"id"`
	KategoriID   *uint         `json:"kategori_id"`
	KategoriNama *string       `gorm:"size:100" json:"kategori_nama"`
	NamaUsaha    string        `gorm:"size:150;not null" json:"nama_usaha"`
	Slug         string        `gorm:"size:150;not null;uniqueIndex" json:"slug"`
	NamaPemilik  string        `gorm:"size:100;not null" json:"nama_pemilik"`
	Deskripsi    *string       `gorm:"type:text" json:"deskripsi"`
	Alamat       *string       `gorm:"type:text" json:"alamat"`
	Telepon      *string       `gorm:"size:20" json:"telepon"`
	WhatsApp     *string       `gorm:"size:20" json:"whatsapp"`
	Email        *string       `gorm:"size:100" json:"email"`
	Instagram    *string       `gorm:"size:100" json:"instagram"`
	FotoID       *uint         `json:"foto_id"`
	IsPublished  bool          `gorm:"default:true" json:"is_published"`
	IsDeleted    bool          `gorm:"default:false" json:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Kategori     *KategoriUMKM `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Foto         *Media        `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
	Produk       []ProdukUMKM  `gorm:"foreignKey:UMKMID" json:"produk,omitempty"`
}

type ProdukUMKM struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UMKMID      uint      `gorm:"not null" json:"umkm_id"`
	Nama        string    `gorm:"size:150;not null" json:"nama"`
	Deskripsi   *string   `gorm:"type:text" json:"deskripsi"`
	Harga       *int64    `json:"harga"`
	FotoID      *uint     `json:"foto_id"`
	IsAvailable bool      `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	Foto        *Media    `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
}
