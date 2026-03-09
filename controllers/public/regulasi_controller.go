// Platform Desa — Public Regulasi Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetKategoriRegulasiList(c *gin.Context) {
	var kategori []models.KategoriRegulasi
	database.DB.Where("is_deleted = false").Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetRegulasiList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriSlug := c.Query("kategori")
	tahun := c.Query("tahun")

	query := database.DB.Model(&models.Regulasi{}).
		Preload("Kategori").
		Preload("File").
		Where("is_published = true AND is_deleted = false")

	if search != "" {
		query = query.Where("judul LIKE ? OR nomor LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriSlug != "" {
		var kategori models.KategoriRegulasi
		if err := database.DB.Where("slug = ? AND is_deleted = false", kategoriSlug).First(&kategori).Error; err == nil {
			query = query.Where("kategori_id = ?", kategori.ID)
		}
	}
	if tahun != "" {
		query = query.Where("tahun = ?", tahun)
	}

	var total int64
	query.Count(&total)

	var regulasi []models.Regulasi
	query.Scopes(helpers.Paginate(pg)).Order("tanggal_terbit DESC, created_at DESC").Find(&regulasi)

	helpers.OKPaginated(c, "Berhasil", regulasi, pg.Meta(total))
}

func GetRegulasi(c *gin.Context) {
	slug := c.Param("slug")

	var regulasi models.Regulasi
	if err := database.DB.
		Preload("Kategori").
		Preload("File").
		Where("slug = ? AND is_published = true AND is_deleted = false", slug).
		First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", regulasi)
}
