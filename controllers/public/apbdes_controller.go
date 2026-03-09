// Platform Desa — Public APBDes Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetAPBDesList(c *gin.Context) {
	var apbdes []models.APBDes
	database.DB.
		Preload("Dokumen").
		Where("is_published = true").
		Order("tahun DESC").
		Find(&apbdes)

	helpers.OK(c, "Berhasil", apbdes)
}

func GetAPBDes(c *gin.Context) {
	tahun := c.Param("tahun")

	var apbdes models.APBDes
	if err := database.DB.
		Preload("Dokumen").
		Preload("Detail").
		Where("tahun = ? AND is_published = true", tahun).
		First(&apbdes).Error; err != nil {
		helpers.NotFound(c, "Data APBDes tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", apbdes)
}
