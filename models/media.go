// Platform Desa — Model Media
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type Media struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Folder       string    `gorm:"size:50;not null" json:"folder"`
	Filename     string    `gorm:"size:255;not null" json:"filename"`
	OriginalName string    `gorm:"size:255;not null" json:"original_name"`
	MimeType     string    `gorm:"size:100;not null" json:"mime_type"`
	Size         int64     `gorm:"not null" json:"size"`
	Path         string    `gorm:"size:500;not null" json:"path"`
	Type         string    `gorm:"size:20;not null" json:"type"`
	UploadedBy   *uint     `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}
