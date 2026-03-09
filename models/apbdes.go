// Platform Desa — Model APBDes & Pengaduan & Regulasi
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type APBDes struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Tahun           int            `gorm:"not null;uniqueIndex" json:"tahun"`
	TotalPendapatan int64          `gorm:"default:0" json:"total_pendapatan"`
	TotalBelanja    int64          `gorm:"default:0" json:"total_belanja"`
	TotalPembiayaan int64          `gorm:"default:0" json:"total_pembiayaan"`
	DokumenID       *uint          `json:"dokumen_id"`
	IsPublished     bool           `gorm:"default:false" json:"is_published"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Dokumen         *Media         `gorm:"foreignKey:DokumenID" json:"dokumen,omitempty"`
	Detail          []APBDesDetail `gorm:"foreignKey:APBDesID" json:"detail,omitempty"`
}

func (APBDes) TableName() string {
	return "apbdes"
}

type APBDesDetail struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	APBDesID  uint   `gorm:"column:apbdes_id;not null" json:"apbdes_id"`
	Jenis     string `gorm:"size:20;not null" json:"jenis"`
	Kategori  string `gorm:"size:100;not null" json:"kategori"`
	Uraian    string `gorm:"size:255;not null" json:"uraian"`
	Anggaran  int64  `gorm:"default:0" json:"anggaran"`
	Realisasi int64  `gorm:"default:0" json:"realisasi"`
	Urutan    uint   `gorm:"default:0" json:"urutan"`
}

func (APBDesDetail) TableName() string {
	return "apbdes_details"
}

type Pengaduan struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	NomorTiket    string    `gorm:"size:50;not null;uniqueIndex" json:"nomor_tiket"`
	NamaPelapor   string    `gorm:"size:100;not null" json:"nama_pelapor"`
	KontakPelapor *string   `gorm:"size:100" json:"kontak_pelapor"`
	IsAnonim      bool      `gorm:"default:false" json:"is_anonim"`
	Kategori      string    `gorm:"size:20;not null" json:"kategori"`
	Judul         string    `gorm:"size:255;not null" json:"judul"`
	Isi           string    `gorm:"type:longtext;not null" json:"isi"`
	Lokasi        *string   `gorm:"size:255" json:"lokasi"`
	FotoID        *uint     `json:"foto_id"`
	Status        string    `gorm:"size:20;default:masuk" json:"status"`
	DitanganiOleh *uint     `json:"ditangani_oleh"`
	ResponAdmin   *string   `gorm:"type:text" json:"respon_admin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Foto          *Media    `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
}

type KategoriRegulasi struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Nama      string `gorm:"size:100;not null" json:"nama"`
	Slug      string `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	IsDeleted bool   `gorm:"default:false" json:"-"`
}

type Regulasi struct {
	ID            uint              `gorm:"primarykey" json:"id"`
	KategoriID    *uint             `json:"kategori_id"`
	KategoriNama  *string           `gorm:"size:100" json:"kategori_nama"`
	Nomor         *string           `gorm:"size:100" json:"nomor"`
	Judul         string            `gorm:"size:255;not null" json:"judul"`
	Slug          string            `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Tentang       *string           `gorm:"type:text" json:"tentang"`
	Tahun         *int              `json:"tahun"`
	TanggalTerbit *time.Time        `json:"tanggal_terbit"`
	FileID        *uint             `json:"file_id"`
	IsPublished   bool              `gorm:"default:true" json:"is_published"`
	IsDeleted     bool              `gorm:"default:false" json:"-"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Kategori      *KategoriRegulasi `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	File          *Media            `gorm:"foreignKey:FileID" json:"file,omitempty"`
}
