// Platform Desa — Public Pemerintahan Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetJabatanList(c *gin.Context) {
	var jabatan []models.Jabatan
	database.DB.Where("is_deleted = false AND parent_id IS NULL").
		Preload("Children").
		Order("urutan ASC, nama ASC").
		Find(&jabatan)

	helpers.OK(c, "Berhasil", jabatan)
}

func GetPejabatList(c *gin.Context) {
	jabatanID := c.Query("jabatan_id")

	query := database.DB.Model(&models.Pejabat{}).
		Preload("Jabatan").
		Preload("Foto").
		Preload("Pendidikan").
		Where("is_deleted = false AND is_active = true")

	if jabatanID != "" {
		query = query.Where("jabatan_id = ?", jabatanID)
	}

	var pejabat []models.Pejabat
	query.Order("urutan ASC, nama ASC").Find(&pejabat)

	helpers.OK(c, "Berhasil", pejabat)
}

func GetPejabat(c *gin.Context) {
	slug := c.Param("slug")

	var pejabat models.Pejabat
	if err := database.DB.
		Preload("Jabatan").
		Preload("Foto").
		Preload("Pendidikan").
		Where("slug = ? AND is_deleted = false AND is_active = true", slug).
		First(&pejabat).Error; err != nil {
		helpers.NotFound(c, "Pejabat tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", pejabat)
}

func GetLembagaList(c *gin.Context) {
	var lembaga []models.LembagaDesa
	database.DB.
		Preload("Logo").
		Where("is_deleted = false AND is_active = true").
		Order("urutan ASC, nama ASC").
		Find(&lembaga)

	helpers.OK(c, "Berhasil", lembaga)
}
