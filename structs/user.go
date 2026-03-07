// Platform Desa — Structs User
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type CreateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=2"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	Role     string  `json:"role" binding:"required,oneof=superadmin admin operator"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=2"`
	Email    string  `json:"email" binding:"required,email"`
	Role     string  `json:"role" binding:"required,oneof=superadmin admin operator"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}
