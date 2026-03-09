// Platform Desa — Admin Dashboard Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	var totalBerita, totalBeritaPublished int64
	var totalPengumuman, totalAgenda int64
	var totalWisata, totalUMKM int64
	var totalPengajuanSurat, totalSuratPending int64
	var totalPengaduan, totalPengaduanBaru int64
	var totalMedia int64
	var totalUser int64

	database.DB.Model(&models.Berita{}).Where("is_deleted = false").Count(&totalBerita)
	database.DB.Model(&models.Berita{}).Where("is_deleted = false AND status = ?", "published").Count(&totalBeritaPublished)

	database.DB.Model(&models.Pengumuman{}).Where("is_deleted = false").Count(&totalPengumuman)
	database.DB.Model(&models.Agenda{}).Where("is_deleted = false").Count(&totalAgenda)

	database.DB.Model(&models.Wisata{}).Where("is_deleted = false").Count(&totalWisata)
	database.DB.Model(&models.UMKM{}).Where("is_deleted = false").Count(&totalUMKM)

	database.DB.Model(&models.PengajuanSurat{}).Count(&totalPengajuanSurat)
	database.DB.Model(&models.PengajuanSurat{}).Where("status = ?", "pending").Count(&totalSuratPending)

	database.DB.Model(&models.Pengaduan{}).Count(&totalPengaduan)
	database.DB.Model(&models.Pengaduan{}).Where("status = ?", "masuk").Count(&totalPengaduanBaru)

	database.DB.Model(&models.Media{}).Count(&totalMedia)
	database.DB.Model(&models.User{}).Where("is_active = true").Count(&totalUser)

	var beritaTerbaru []models.Berita
	database.DB.Where("is_deleted = false").
		Order("created_at DESC").
		Limit(5).
		Find(&beritaTerbaru)

	var suratTerbaru []models.PengajuanSurat
	database.DB.Preload("JenisSurat").
		Order("created_at DESC").
		Limit(5).
		Find(&suratTerbaru)

	var pengaduanTerbaru []models.Pengaduan
	database.DB.Order("created_at DESC").
		Limit(5).
		Find(&pengaduanTerbaru)

	var statistikPenduduk models.StatistikPenduduk
	database.DB.Order("tahun DESC, bulan DESC").First(&statistikPenduduk)

	helpers.OK(c, "Berhasil", gin.H{
		"statistik": gin.H{
			"berita": gin.H{
				"total":     totalBerita,
				"published": totalBeritaPublished,
				"draft":     totalBerita - totalBeritaPublished,
			},
			"pengumuman": totalPengumuman,
			"agenda":     totalAgenda,
			"wisata":     totalWisata,
			"umkm":       totalUMKM,
			"surat": gin.H{
				"total":   totalPengajuanSurat,
				"pending": totalSuratPending,
			},
			"pengaduan": gin.H{
				"total": totalPengaduan,
				"baru":  totalPengaduanBaru,
			},
			"media": totalMedia,
			"user":  totalUser,
		},
		"penduduk":          statistikPenduduk,
		"berita_terbaru":    beritaTerbaru,
		"surat_terbaru":     suratTerbaru,
		"pengaduan_terbaru": pengaduanTerbaru,
	})
}
