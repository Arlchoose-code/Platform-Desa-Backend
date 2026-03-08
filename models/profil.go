// Platform Desa — Model Profil Desa
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type ProfilDesa struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	NamaDesa      string    `gorm:"size:100;not null" json:"nama_desa"`
	NamaKecamatan *string   `gorm:"size:100" json:"nama_kecamatan"`
	NamaKabupaten *string   `gorm:"size:100" json:"nama_kabupaten"`
	NamaProvinsi  *string   `gorm:"size:100" json:"nama_provinsi"`
	KodePos       *string   `gorm:"size:10" json:"kode_pos"`
	Alamat        *string   `gorm:"type:text" json:"alamat"`
	Telepon       *string   `gorm:"size:20" json:"telepon"`
	Email         *string   `gorm:"size:100" json:"email"`
	Website       *string   `gorm:"size:255" json:"website"`
	Facebook      *string   `gorm:"size:255" json:"facebook"`
	Instagram     *string   `gorm:"size:255" json:"instagram"`
	Twitter       *string   `gorm:"size:255" json:"twitter"`
	Youtube       *string   `gorm:"size:255" json:"youtube"`
	Tiktok        *string   `gorm:"size:255" json:"tiktok"`
	LuasWilayah   *float64  `json:"luas_wilayah"`
	JumlahDusun   uint      `gorm:"default:0" json:"jumlah_dusun"`
	JumlahRW      uint      `gorm:"default:0" json:"jumlah_rw"`
	JumlahRT      uint      `gorm:"default:0" json:"jumlah_rt"`
	BatasUtara    *string   `gorm:"size:255" json:"batas_utara"`
	BatasSelatan  *string   `gorm:"size:255" json:"batas_selatan"`
	BatasTimur    *string   `gorm:"size:255" json:"batas_timur"`
	BatasBarat    *string   `gorm:"size:255" json:"batas_barat"`
	Latitude      *float64  `json:"latitude"`
	Longitude     *float64  `json:"longitude"`
	LogoID        *uint     `json:"logo_id"`
	FotoDesaID    *uint     `json:"foto_desa_id"`
	Sejarah       *string   `gorm:"type:longtext" json:"sejarah"`
	Visi          *string   `gorm:"type:text" json:"visi"`
	Misi          *string   `gorm:"type:longtext" json:"misi"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Logo          *Media    `gorm:"foreignKey:LogoID" json:"logo,omitempty"`
	FotoDesa      *Media    `gorm:"foreignKey:FotoDesaID" json:"foto_desa,omitempty"`
}

type PotensiDesa struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Kategori    string    `gorm:"size:20;not null" json:"kategori"`
	Judul       string    `gorm:"size:255;not null" json:"judul"`
	Deskripsi   *string   `gorm:"type:longtext" json:"deskripsi"`
	FotoID      *uint     `json:"foto_id"`
	Urutan      uint      `gorm:"default:0" json:"urutan"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	IsDeleted   bool      `gorm:"default:false" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Foto        *Media    `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
}
