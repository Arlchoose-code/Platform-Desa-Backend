// Platform Desa — Public Profil Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetProfilDesa(c *gin.Context) {
	var profil models.ProfilDesa
	if err := database.DB.Preload("Logo").Preload("FotoDesa").First(&profil).Error; err != nil {
		helpers.NotFound(c, "Profil desa belum dikonfigurasi")
		return
	}

	helpers.OK(c, "Berhasil", profil)
}

func GetPotensiList(c *gin.Context) {
	kategori := c.Query("kategori")

	query := database.DB.Model(&models.PotensiDesa{}).
		Preload("Foto").
		Where("is_published = true AND is_deleted = false")

	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	var potensi []models.PotensiDesa
	query.Order("urutan ASC, created_at DESC").Find(&potensi)

	helpers.OK(c, "Berhasil", potensi)
}
