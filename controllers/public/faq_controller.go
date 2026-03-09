// Platform Desa — Public FAQ Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetFAQList(c *gin.Context) {
	kategori := c.Query("kategori")
	search := c.Query("search")

	query := database.DB.Model(&models.FAQ{}).
		Where("is_published = true AND is_deleted = false")

	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if search != "" {
		query = query.Where("pertanyaan LIKE ? OR jawaban LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var faq []models.FAQ
	query.Order("urutan ASC, created_at ASC").Find(&faq)

	helpers.OK(c, "Berhasil", faq)
}
