// Platform Desa — Structs Auth
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Name  string  `json:"name" binding:"omitempty,min=2"`
	Email string  `json:"email" binding:"omitempty,email"`
	Phone *string `json:"phone"`
}

type ChangePasswordRequest struct {
	PasswordLama string `json:"password_lama" binding:"required"`
	PasswordBaru string `json:"password_baru" binding:"required,min=8"`
}
