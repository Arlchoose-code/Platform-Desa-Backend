// Platform Desa — Public Pengaduan Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
)

func BuatPengaduanPublik(c *gin.Context) {
	var req structs.PengaduanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	nomorTiket := helpers.GenerateNomorTiket()

	namaPelapor := req.NamaPelapor
	if req.IsAnonim {
		namaPelapor = "Anonim"
	}

	pengaduan := models.Pengaduan{
		NomorTiket:    nomorTiket,
		NamaPelapor:   namaPelapor,
		KontakPelapor: req.KontakPelapor,
		IsAnonim:      req.IsAnonim,
		Kategori:      req.Kategori,
		Judul:         req.Judul,
		Isi:           req.Isi,
		Lokasi:        req.Lokasi,
		FotoID:        req.FotoID,
		Status:        "masuk",
	}

	if err := database.DB.Create(&pengaduan).Error; err != nil {
		helpers.InternalError(c, "Gagal mengirim pengaduan")
		return
	}

	helpers.Created(c, "Pengaduan berhasil dikirim", gin.H{
		"nomor_tiket": nomorTiket,
	})
}

func CekStatusPengaduan(c *gin.Context) {
	nomor := c.Param("nomor")

	var pengaduan models.Pengaduan
	if err := database.DB.
		Where("nomor_tiket = ?", nomor).
		First(&pengaduan).Error; err != nil {
		helpers.NotFound(c, "Nomor tiket tidak ditemukan")
		return
	}

	// sembunyikan data sensitif jika anonim
	result := gin.H{
		"nomor_tiket":  pengaduan.NomorTiket,
		"judul":        pengaduan.Judul,
		"kategori":     pengaduan.Kategori,
		"status":       pengaduan.Status,
		"respon_admin": pengaduan.ResponAdmin,
		"created_at":   pengaduan.CreatedAt,
		"updated_at":   pengaduan.UpdatedAt,
	}

	helpers.OK(c, "Berhasil", result)
}
