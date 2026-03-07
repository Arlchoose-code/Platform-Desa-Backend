// Platform Desa — Model Setting
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type Setting struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"size:100;not null;uniqueIndex" json:"key"`
	Value     *string   `gorm:"type:text" json:"value"`
	Group     string    `gorm:"size:50;not null;default:general" json:"group"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
