// Platform Desa — Model User
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type User struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Email       string     `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	Role        string     `gorm:"size:20;not null;default:operator" json:"role"`
	Avatar      *string    `gorm:"size:255" json:"avatar"`
	Phone       *string    `gorm:"size:20" json:"phone"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ActivityLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      *uint     `json:"user_id"`
	Action      string    `gorm:"size:100;not null" json:"action"`
	Module      string    `gorm:"size:50;not null" json:"module"`
	Description *string   `gorm:"type:text" json:"description"`
	IPAddress   *string   `gorm:"size:45" json:"ip_address"`
	CreatedAt   time.Time `json:"created_at"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
