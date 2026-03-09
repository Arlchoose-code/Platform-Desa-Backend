// Platform Desa — Public Berita Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetKategoriBeritaList(c *gin.Context) {
	var kategori []models.KategoriBerita
	database.DB.Where("is_deleted = false").Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetBeritaList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriSlug := c.Query("kategori")

	query := database.DB.Model(&models.Berita{}).
		Preload("Kategori").
		Preload("Thumbnail").
		Preload("Penulis").
		Where("status = 'published' AND is_deleted = false")

	if search != "" {
		query = query.Where("judul LIKE ? OR ringkasan LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriSlug != "" {
		var kategori models.KategoriBerita
		if err := database.DB.Where("slug = ? AND is_deleted = false", kategoriSlug).First(&kategori).Error; err == nil {
			query = query.Where("kategori_id = ?", kategori.ID)
		}
	}

	var total int64
	query.Count(&total)

	var berita []models.Berita
	query.Scopes(helpers.Paginate(pg)).Order("published_at DESC").Find(&berita)

	helpers.OKPaginated(c, "Berhasil", berita, pg.Meta(total))
}

func GetBerita(c *gin.Context) {
	slug := c.Param("slug")

	var berita models.Berita
	if err := database.DB.
		Preload("Kategori").
		Preload("Thumbnail").
		Preload("Penulis").
		Where("slug = ? AND status = 'published' AND is_deleted = false", slug).
		First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan")
		return
	}

	// increment views
	database.DB.Model(&berita).UpdateColumn("views", berita.Views+1)

	helpers.OK(c, "Berhasil", berita)
}
