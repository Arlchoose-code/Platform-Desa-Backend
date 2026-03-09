// Platform Desa — Public Peta Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetPetaFasilitasList(c *gin.Context) {
	kategori := c.Query("kategori")

	query := database.DB.Model(&models.PetaFasilitas{}).
		Preload("Foto").
		Where("is_published = true")

	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	var peta []models.PetaFasilitas
	query.Order("nama ASC").Find(&peta)

	helpers.OK(c, "Berhasil", peta)
}
