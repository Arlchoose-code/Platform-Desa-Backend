// Platform Desa — Admin User Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"strconv"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	pg := helpers.GetPagination(c)

	var users []models.User
	var total int64

	q := database.DB.Model(&models.User{})

	if role := c.Query("role"); role != "" {
		q = q.Where("role = ?", role)
	}
	if active := c.Query("is_active"); active != "" {
		q = q.Where("is_active = ?", active == "true")
	}
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR email LIKE ?", like, like)
	}

	q.Count(&total)
	q.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&users)

	helpers.OKPaginated(c, "Berhasil", users, pg.Meta(total))
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", user)
}

func CreateUser(c *gin.Context) {
	var req structs.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Email sudah digunakan", nil)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.InternalError(c, "Gagal memproses password")
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
		Phone:    req.Phone,
		IsActive: true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		helpers.InternalError(c, "Gagal membuat user")
		return
	}

	userID := c.GetUint("user_id")
	desc := "Membuat user baru: " + user.Email
	database.DB.Create(&models.ActivityLog{
		UserID: &userID, Action: "create", Module: "user", Description: &desc,
	})

	helpers.Created(c, "User berhasil dibuat", user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	var req structs.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Email sudah digunakan user lain", nil)
		return
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	user.Phone = req.Phone
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&user).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui user")
		return
	}

	userID := c.GetUint("user_id")
	desc := "Memperbarui user: " + user.Email
	database.DB.Create(&models.ActivityLog{
		UserID: &userID, Action: "update", Module: "user", Description: &desc,
	})

	helpers.OK(c, "User berhasil diperbarui", user)
}

func UpdateAvatar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		helpers.BadRequest(c, "File avatar tidak ditemukan", nil)
		return
	}

	result, err := helpers.UploadFile(file, "avatar")
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	if user.Avatar != nil && *user.Avatar != "" {
		helpers.DeleteFile(*user.Avatar)
	}

	database.DB.Model(&user).Update("avatar", result.Path)

	userID := c.GetUint("user_id")
	desc := "Update avatar user: " + user.Email
	database.DB.Create(&models.ActivityLog{
		UserID: &userID, Action: "update_avatar", Module: "user", Description: &desc,
	})

	helpers.OK(c, "Avatar berhasil diperbarui", gin.H{"avatar": result.Path})
}

func DeleteAvatar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	if user.Avatar == nil || *user.Avatar == "" {
		helpers.BadRequest(c, "User tidak memiliki avatar", nil)
		return
	}

	helpers.DeleteFile(*user.Avatar)
	database.DB.Model(&user).Update("avatar", nil)

	helpers.OK(c, "Avatar berhasil dihapus", nil)
}

func ResetPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.InternalError(c, "Gagal memproses password")
		return
	}

	database.DB.Model(&user).Update("password", string(hashed))

	userID := c.GetUint("user_id")
	desc := "Reset password user: " + user.Email
	database.DB.Create(&models.ActivityLog{
		UserID: &userID, Action: "reset_password", Module: "user", Description: &desc,
	})

	helpers.OK(c, "Password berhasil direset", nil)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	if uint(id) == c.GetUint("user_id") {
		helpers.BadRequest(c, "Tidak dapat menghapus akun sendiri", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	if user.Avatar != nil && *user.Avatar != "" {
		helpers.DeleteFile(*user.Avatar)
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		helpers.InternalError(c, "Gagal menghapus user")
		return
	}

	userID := c.GetUint("user_id")
	desc := "Menghapus user: " + user.Email
	database.DB.Create(&models.ActivityLog{
		UserID: &userID, Action: "delete", Module: "user", Description: &desc,
	})

	helpers.OK(c, "User berhasil dihapus", nil)
}

func GetActivityLogs(c *gin.Context) {
	pg := helpers.GetPagination(c)

	var logs []models.ActivityLog
	var total int64

	q := database.DB.Model(&models.ActivityLog{}).Preload("User")

	if module := c.Query("module"); module != "" {
		q = q.Where("module = ?", module)
	}
	if action := c.Query("action"); action != "" {
		q = q.Where("action = ?", action)
	}
	if userID := c.Query("user_id"); userID != "" {
		q = q.Where("user_id = ?", userID)
	}

	q.Count(&total)
	q.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&logs)

	helpers.OKPaginated(c, "Berhasil", logs, pg.Meta(total))
}
