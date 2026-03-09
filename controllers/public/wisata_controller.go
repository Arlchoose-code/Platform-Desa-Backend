// Platform Desa — Public Wisata Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetKategoriWisataList(c *gin.Context) {
	var kategori []models.KategoriWisata
	database.DB.Where("is_deleted = false").Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetWisataList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriSlug := c.Query("kategori")

	query := database.DB.Model(&models.Wisata{}).
		Preload("Kategori").
		Preload("Thumbnail").
		Where("is_published = true AND is_deleted = false")

	if search != "" {
		query = query.Where("nama LIKE ? OR deskripsi LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriSlug != "" {
		var kategori models.KategoriWisata
		if err := database.DB.Where("slug = ? AND is_deleted = false", kategoriSlug).First(&kategori).Error; err == nil {
			query = query.Where("kategori_id = ?", kategori.ID)
		}
	}

	var total int64
	query.Count(&total)

	var wisata []models.Wisata
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&wisata)

	helpers.OKPaginated(c, "Berhasil", wisata, pg.Meta(total))
}

func GetWisata(c *gin.Context) {
	slug := c.Param("slug")

	var wisata models.Wisata
	if err := database.DB.
		Preload("Kategori").
		Preload("Thumbnail").
		Preload("Galeri.Media").
		Where("slug = ? AND is_published = true AND is_deleted = false", slug).
		First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan")
		return
	}

	// increment views
	database.DB.Model(&wisata).UpdateColumn("views", wisata.Views+1)

	helpers.OK(c, "Berhasil", wisata)
}
