// Platform Desa — Public Surat Controller
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

func GetJenisSuratPublikList(c *gin.Context) {
	var jenis []models.JenisSurat
	database.DB.Where("is_active = true AND is_deleted = false").
		Order("nama ASC").
		Find(&jenis)

	helpers.OK(c, "Berhasil", jenis)
}

func AjukanSuratPublik(c *gin.Context) {
	var req structs.AjukanSuratRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Where("id = ? AND is_active = true AND is_deleted = false", req.JenisSuratID).First(&jenis).Error; err != nil {
		helpers.NotFound(c, "Jenis surat tidak tersedia")
		return
	}

	nomorPengajuan := helpers.GenerateNomorPengajuan()

	pengajuan := models.PengajuanSurat{
		NomorPengajuan: nomorPengajuan,
		JenisSuratID:   req.JenisSuratID,
		JenisSuratNama: jenis.Nama,
		NamaPemohon:    req.NamaPemohon,
		NIK:            req.NIK,
		Keperluan:      req.Keperluan,
		DataTambahan:   req.DataTambahan,
		Status:         "pending",
	}

	if err := database.DB.Create(&pengajuan).Error; err != nil {
		helpers.InternalError(c, "Gagal mengajukan surat")
		return
	}

	helpers.Created(c, "Pengajuan surat berhasil", gin.H{
		"nomor_pengajuan": nomorPengajuan,
	})
}

func CekStatusSurat(c *gin.Context) {
	nomor := c.Param("nomor")

	var pengajuan models.PengajuanSurat
	if err := database.DB.
		Preload("JenisSurat").
		Preload("FileHasil").
		Preload("Riwayat").
		Where("nomor_pengajuan = ?", nomor).
		First(&pengajuan).Error; err != nil {
		helpers.NotFound(c, "Nomor pengajuan tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", pengajuan)
}
