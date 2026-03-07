// Platform Desa — Model Persuratan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type JenisSurat struct {
	ID             uint    `gorm:"primarykey" json:"id"`
	Nama           string  `gorm:"size:150;not null" json:"nama"`
	Kode           string  `gorm:"size:20;not null;uniqueIndex" json:"kode"`
	Deskripsi      *string `gorm:"type:text" json:"deskripsi"`
	Syarat         *string `gorm:"type:json" json:"syarat"`
	TemplateFileID *uint   `json:"template_file_id"`
	EstimasiHari   uint    `gorm:"default:3" json:"estimasi_hari"`
	IsActive       bool    `gorm:"default:true" json:"is_active"`
	IsDeleted      bool    `gorm:"default:false" json:"-"`
	TemplateFile   *Media  `gorm:"foreignKey:TemplateFileID" json:"template_file,omitempty"`
}

type PengajuanSurat struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	NomorPengajuan string         `gorm:"size:50;not null;uniqueIndex" json:"nomor_pengajuan"`
	JenisSuratID   uint           `gorm:"not null" json:"jenis_surat_id"`
	JenisSuratNama string         `gorm:"size:150;not null" json:"jenis_surat_nama"`
	PemohonID      *uint          `json:"pemohon_id"`
	NamaPemohon    string         `gorm:"size:100;not null" json:"nama_pemohon"`
	NIK            string         `gorm:"size:16;not null" json:"nik"`
	Keperluan      *string        `gorm:"type:text" json:"keperluan"`
	DataTambahan   *string        `gorm:"type:json" json:"data_tambahan"`
	Status         string         `gorm:"size:20;default:pending" json:"status"`
	CatatanAdmin   *string        `gorm:"type:text" json:"catatan_admin"`
	FileHasilID    *uint          `json:"file_hasil_id"`
	DipresesOleh   *uint          `json:"diproses_oleh"`
	DiprosesPada   *time.Time     `json:"diproses_pada"`
	SelesaiPada    *time.Time     `json:"selesai_pada"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	JenisSurat     *JenisSurat    `gorm:"foreignKey:JenisSuratID" json:"jenis_surat,omitempty"`
	FileHasil      *Media         `gorm:"foreignKey:FileHasilID" json:"file_hasil,omitempty"`
	Riwayat        []RiwayatSurat `gorm:"foreignKey:PengajuanID" json:"riwayat,omitempty"`
}

type RiwayatSurat struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	PengajuanID uint      `gorm:"not null" json:"pengajuan_id"`
	Status      string    `gorm:"size:20;not null" json:"status"`
	Catatan     *string   `gorm:"type:text" json:"catatan"`
	Oleh        *uint     `json:"oleh"`
	CreatedAt   time.Time `json:"created_at"`
	OlehUser    *User     `gorm:"foreignKey:Oleh" json:"oleh_user,omitempty"`
}
