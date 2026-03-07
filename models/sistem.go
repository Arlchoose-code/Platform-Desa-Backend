// Platform Desa — Model Peta, FAQ, Tema, System
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type PetaFasilitas struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Nama        string    `gorm:"size:150;not null" json:"nama"`
	Kategori    string    `gorm:"size:20;not null" json:"kategori"`
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

type Tema struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Nama        string    `gorm:"size:100;not null" json:"nama"`
	Slug        string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Deskripsi   *string   `gorm:"type:text" json:"deskripsi"`
	Versi       *string   `gorm:"size:20" json:"versi"`
	Author      *string   `gorm:"size:100" json:"author"`
	RepoURL     *string   `gorm:"size:500" json:"repo_url"`
	PreviewURL  *string   `gorm:"size:500" json:"preview_url"`
	ThumbnailID *uint     `json:"thumbnail_id"`
	IsActive    bool      `gorm:"default:false" json:"is_active"`
	InstalledAt time.Time `json:"installed_at"`
	Thumbnail   *Media    `gorm:"foreignKey:ThumbnailID" json:"thumbnail,omitempty"`
}

type UpdateLog struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	VersiLama     string     `gorm:"size:20;not null" json:"versi_lama"`
	VersiBaru     string     `gorm:"size:20;not null" json:"versi_baru"`
	Status        string     `gorm:"size:20;default:pending" json:"status"`
	Changelog     *string    `gorm:"type:text" json:"changelog"`
	BackupPath    *string    `gorm:"size:500" json:"backup_path"`
	ErrorLog      *string    `gorm:"type:text" json:"error_log"`
	DilakukanOleh *uint      `json:"dilakukan_oleh"`
	StartedAt     time.Time  `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	OlehUser      *User      `gorm:"foreignKey:DilakukanOleh" json:"oleh_user,omitempty"`
}
