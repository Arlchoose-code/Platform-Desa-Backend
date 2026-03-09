// Platform Desa — Public UMKM Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetKategoriUMKMList(c *gin.Context) {
	var kategori []models.KategoriUMKM
	database.DB.Where("is_deleted = false").Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetUMKMList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriSlug := c.Query("kategori")

	query := database.DB.Model(&models.UMKM{}).
		Preload("Kategori").
		Preload("Foto").
		Where("is_published = true AND is_deleted = false")

	if search != "" {
		query = query.Where("nama_usaha LIKE ? OR nama_pemilik LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriSlug != "" {
		var kategori models.KategoriUMKM
		if err := database.DB.Where("slug = ? AND is_deleted = false", kategoriSlug).First(&kategori).Error; err == nil {
			query = query.Where("kategori_id = ?", kategori.ID)
		}
	}

	var total int64
	query.Count(&total)

	var umkm []models.UMKM
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&umkm)

	helpers.OKPaginated(c, "Berhasil", umkm, pg.Meta(total))
}

func GetUMKM(c *gin.Context) {
	slug := c.Param("slug")

	var umkm models.UMKM
	if err := database.DB.
		Preload("Kategori").
		Preload("Foto").
		Preload("Produk.Foto").
		Where("slug = ? AND is_published = true AND is_deleted = false", slug).
		First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", umkm)
}
