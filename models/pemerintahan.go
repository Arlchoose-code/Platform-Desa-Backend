// Platform Desa — Model Pemerintahan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type Jabatan struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Nama      string    `gorm:"size:100;not null" json:"nama"`
	Level     uint      `gorm:"default:0" json:"level"`
	ParentID  *uint     `json:"parent_id"`
	Urutan    uint      `gorm:"default:0" json:"urutan"`
	IsDeleted bool      `gorm:"default:false" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Parent    *Jabatan  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Jabatan `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type Pejabat struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	JabatanID      uint      `gorm:"not null" json:"jabatan_id"`
	JabatanNama    string    `gorm:"size:100;not null" json:"jabatan_nama"`
	Nama           string    `gorm:"size:100;not null" json:"nama"`
	NIP            *string   `gorm:"size:50" json:"nip"`
	FotoID         *uint     `json:"foto_id"`
	Pendidikan     *string   `gorm:"size:100" json:"pendidikan"`
	PeriodeMulai   *int      `json:"periode_mulai"`
	PeriodeSelesai *int      `json:"periode_selesai"`
	Biodata        *string   `gorm:"type:text" json:"biodata"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	Urutan         uint      `gorm:"default:0" json:"urutan"`
	IsDeleted      bool      `gorm:"default:false" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Jabatan        *Jabatan  `gorm:"foreignKey:JabatanID" json:"jabatan,omitempty"`
	Foto           *Media    `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
}

type LembagaDesa struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Nama      string    `gorm:"size:150;not null" json:"nama"`
	Singkatan *string   `gorm:"size:20" json:"singkatan"`
	Deskripsi *string   `gorm:"type:text" json:"deskripsi"`
	Ketua     *string   `gorm:"size:100" json:"ketua"`
	LogoID    *uint     `json:"logo_id"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	Urutan    uint      `gorm:"default:0" json:"urutan"`
	IsDeleted bool      `gorm:"default:false" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Logo      *Media    `gorm:"foreignKey:LogoID" json:"logo,omitempty"`
}
