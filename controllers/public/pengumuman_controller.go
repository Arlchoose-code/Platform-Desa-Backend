// Platform Desa — Public Pengumuman & Agenda Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetPengumumanList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	penting := c.Query("penting")

	query := database.DB.Model(&models.Pengumuman{}).
		Preload("File").
		Preload("Penulis").
		Where("status = 'published' AND is_deleted = false")

	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}
	if penting == "true" {
		query = query.Where("is_penting = true")
	}

	var total int64
	query.Count(&total)

	var pengumuman []models.Pengumuman
	query.Scopes(helpers.Paginate(pg)).Order("is_penting DESC, published_at DESC").Find(&pengumuman)

	helpers.OKPaginated(c, "Berhasil", pengumuman, pg.Meta(total))
}

func GetPengumuman(c *gin.Context) {
	slug := c.Param("slug")

	var pengumuman models.Pengumuman
	if err := database.DB.
		Preload("File").
		Preload("Penulis").
		Where("slug = ? AND status = 'published' AND is_deleted = false", slug).
		First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", pengumuman)
}

func GetAgendaList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	status := c.Query("status")

	query := database.DB.Model(&models.Agenda{}).
		Where("is_published = true AND is_deleted = false")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var agenda []models.Agenda
	query.Scopes(helpers.Paginate(pg)).Order("tanggal_mulai ASC").Find(&agenda)

	helpers.OKPaginated(c, "Berhasil", agenda, pg.Meta(total))
}

func GetAgenda(c *gin.Context) {
	slug := c.Param("slug")

	var agenda models.Agenda
	if err := database.DB.
		Where("slug = ? AND is_published = true AND is_deleted = false", slug).
		First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", agenda)
}
