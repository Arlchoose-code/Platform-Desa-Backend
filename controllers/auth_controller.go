// Platform Desa — Auth Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package controllers

import (
	"time"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req structs.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("email = ? AND is_active = true", req.Email).First(&user).Error; err != nil {
		helpers.Unauthorized(c, "Email atau password salah")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		helpers.Unauthorized(c, "Email atau password salah")
		return
	}

	accessToken, err := helpers.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		helpers.InternalError(c, "Gagal membuat token")
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		helpers.InternalError(c, "Gagal membuat token")
		return
	}

	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", now)

	ip := c.ClientIP()
	desc := "Login berhasil"
	database.DB.Create(&models.ActivityLog{
		UserID:      &user.ID,
		Action:      "login",
		Module:      "auth",
		Description: &desc,
		IPAddress:   &ip,
	})

	helpers.OK(c, "Login berhasil", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	database.DB.Create(&models.ActivityLog{
		UserID: &userID,
		Action: "logout",
		Module: "auth",
	})
	helpers.OK(c, "Logout berhasil", nil)
}

func RefreshToken(c *gin.Context) {
	var req structs.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	claims, err := helpers.ParseToken(req.RefreshToken)
	if err != nil || claims.Type != "refresh" {
		helpers.Unauthorized(c, "Refresh token tidak valid")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ? AND is_active = true", claims.UserID).First(&user).Error; err != nil {
		helpers.Unauthorized(c, "User tidak ditemukan atau tidak aktif")
		return
	}

	accessToken, err := helpers.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		helpers.InternalError(c, "Gagal membuat token")
		return
	}

	helpers.OK(c, "Token diperbarui", gin.H{
		"access_token": accessToken,
	})
}

func GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}
	helpers.OK(c, "Berhasil", user)
}

func UpdateMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req structs.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" && req.Email != user.Email {
		var existing models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existing).Error; err == nil {
			helpers.BadRequest(c, "Email sudah digunakan", nil)
			return
		}
		user.Email = req.Email
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if err := database.DB.Save(&user).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui profil")
		return
	}

	helpers.OK(c, "Profil berhasil diperbarui", user)
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req structs.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.PasswordLama)); err != nil {
		helpers.BadRequest(c, "Password lama tidak sesuai", nil)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		helpers.InternalError(c, "Gagal memproses password")
		return
	}

	database.DB.Model(&user).Update("password", string(hashed))
	helpers.OK(c, "Password berhasil diubah", nil)
}
