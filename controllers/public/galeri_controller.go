// Platform Desa — Public Galeri Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetKategoriGaleriList(c *gin.Context) {
	var kategori []models.KategoriGaleri
	database.DB.Where("is_deleted = false").Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetGaleriList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	kategoriSlug := c.Query("kategori")

	query := database.DB.Model(&models.Galeri{}).
		Preload("Kategori").
		Preload("Media").
		Preload("Thumbnail").
		Where("is_published = true AND is_deleted = false")

	if kategoriSlug != "" {
		var kategori models.KategoriGaleri
		if err := database.DB.Where("slug = ? AND is_deleted = false", kategoriSlug).First(&kategori).Error; err == nil {
			query = query.Where("kategori_id = ?", kategori.ID)
		}
	}

	var total int64
	query.Count(&total)

	var galeri []models.Galeri
	query.Scopes(helpers.Paginate(pg)).Order("urutan ASC, tanggal DESC, created_at DESC").Find(&galeri)

	helpers.OKPaginated(c, "Berhasil", galeri, pg.Meta(total))
}
