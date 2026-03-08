// Platform Desa — Admin Surat Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"strconv"
	"time"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
)

func GetJenisSuratList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")
	isActive := c.Query("is_active")

	query := database.DB.Model(&models.JenisSurat{}).
		Preload("TemplateFile").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ? OR kode LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var jenis []models.JenisSurat
	query.Order("nama ASC").Find(&jenis)

	helpers.OK(c, "Berhasil", jenis)
}

func GetJenisSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Preload("TemplateFile").
		Where("id = ? AND is_deleted = false", id).
		First(&jenis).Error; err != nil {
		helpers.NotFound(c, "Jenis surat tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", jenis)
}

func CreateJenisSurat(c *gin.Context) {
	var req structs.JenisSuratRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.JenisSurat
	if err := database.DB.Where("kode = ? AND is_deleted = false", req.Kode).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kode surat sudah ada", nil)
		return
	}

	jenis := models.JenisSurat{
		Nama:           req.Nama,
		Kode:           req.Kode,
		Deskripsi:      req.Deskripsi,
		Syarat:         req.Syarat,
		TemplateFileID: req.TemplateFileID,
		EstimasiHari:   req.EstimasiHari,
		IsActive:       req.IsActive,
	}

	if err := database.DB.Create(&jenis).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah jenis surat")
		return
	}

	database.DB.Preload("TemplateFile").First(&jenis, jenis.ID)
	helpers.Created(c, "Jenis surat berhasil ditambahkan", jenis)
}

func UpdateJenisSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&jenis).Error; err != nil {
		helpers.NotFound(c, "Jenis surat tidak ditemukan")
		return
	}

	var req structs.JenisSuratRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.JenisSurat
	if err := database.DB.Where("kode = ? AND id != ? AND is_deleted = false", req.Kode, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kode surat sudah ada", nil)
		return
	}

	jenis.Nama = req.Nama
	jenis.Kode = req.Kode
	jenis.Deskripsi = req.Deskripsi
	jenis.Syarat = req.Syarat
	jenis.TemplateFileID = req.TemplateFileID
	jenis.EstimasiHari = req.EstimasiHari
	jenis.IsActive = req.IsActive

	if err := database.DB.Save(&jenis).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui jenis surat")
		return
	}

	database.DB.Preload("TemplateFile").First(&jenis, jenis.ID)
	helpers.OK(c, "Jenis surat berhasil diperbarui", jenis)
}

func DeleteJenisSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&jenis).Error; err != nil {
		helpers.NotFound(c, "Jenis surat tidak ditemukan")
		return
	}

	database.DB.Model(&jenis).Update("is_deleted", true)
	helpers.OK(c, "Jenis surat berhasil dihapus", nil)
}

func ForceDeleteJenisSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&jenis).Error; err != nil {
		helpers.NotFound(c, "Jenis surat tidak ditemukan di trash")
		return
	}

	if jenis.TemplateFileID != nil {
		var media models.Media
		if err := database.DB.First(&media, jenis.TemplateFileID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&jenis)
	helpers.OK(c, "Jenis surat berhasil dihapus permanen", nil)
}

func RestoreJenisSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&jenis).Error; err != nil {
		helpers.NotFound(c, "Jenis surat tidak ditemukan di trash")
		return
	}

	database.DB.Model(&jenis).Update("is_deleted", false)
	helpers.OK(c, "Jenis surat berhasil dipulihkan", nil)
}

func GetPengajuanSuratList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	status := c.Query("status")
	jenisSuratID := c.Query("jenis_surat_id")

	query := database.DB.Model(&models.PengajuanSurat{}).
		Preload("JenisSurat").
		Preload("FileHasil")

	if search != "" {
		query = query.Where("nama_pemohon LIKE ? OR nik LIKE ? OR nomor_pengajuan LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if jenisSuratID != "" {
		query = query.Where("jenis_surat_id = ?", jenisSuratID)
	}

	var total int64
	query.Count(&total)

	var pengajuan []models.PengajuanSurat
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&pengajuan)

	helpers.OKPaginated(c, "Berhasil", pengajuan, pg.Meta(total))
}

func GetPengajuanSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengajuan models.PengajuanSurat
	if err := database.DB.Preload("JenisSurat").Preload("FileHasil").Preload("Riwayat.OlehUser").
		First(&pengajuan, id).Error; err != nil {
		helpers.NotFound(c, "Pengajuan tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", pengajuan)
}

func AjukanSurat(c *gin.Context) {
	var req structs.AjukanSuratRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var jenis models.JenisSurat
	if err := database.DB.Where("id = ? AND is_active = true AND is_deleted = false", req.JenisSuratID).First(&jenis).Error; err != nil {
		helpers.BadRequest(c, "Jenis surat tidak tersedia", nil)
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

	database.DB.Create(&models.RiwayatSurat{
		PengajuanID: pengajuan.ID,
		Status:      "pending",
	})

	helpers.Created(c, "Pengajuan berhasil dibuat", gin.H{
		"nomor_pengajuan": nomorPengajuan,
		"id":              pengajuan.ID,
	})
}

func ProsesSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengajuan models.PengajuanSurat
	if err := database.DB.First(&pengajuan, id).Error; err != nil {
		helpers.NotFound(c, "Pengajuan tidak ditemukan")
		return
	}

	if pengajuan.Status != "pending" {
		helpers.BadRequest(c, "Pengajuan tidak dalam status pending", nil)
		return
	}

	var req structs.ProsesSuratRequest
	c.ShouldBindJSON(&req)

	userID := c.GetUint("user_id")
	now := time.Now()

	database.DB.Model(&pengajuan).Updates(map[string]interface{}{
		"status":        "diproses",
		"catatan_admin": req.CatatanAdmin,
		"diproses_oleh": userID,
		"diproses_pada": now,
	})

	database.DB.Create(&models.RiwayatSurat{
		PengajuanID: pengajuan.ID,
		Status:      "diproses",
		Catatan:     req.CatatanAdmin,
		Oleh:        &userID,
	})

	helpers.OK(c, "Pengajuan sedang diproses", nil)
}

func SelesaikanSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengajuan models.PengajuanSurat
	if err := database.DB.First(&pengajuan, id).Error; err != nil {
		helpers.NotFound(c, "Pengajuan tidak ditemukan")
		return
	}

	if pengajuan.Status != "diproses" {
		helpers.BadRequest(c, "Pengajuan tidak dalam status diproses", nil)
		return
	}

	var req structs.SelesaikanSuratRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	userID := c.GetUint("user_id")
	now := time.Now()

	database.DB.Model(&pengajuan).Updates(map[string]interface{}{
		"status":        "selesai",
		"file_hasil_id": req.FileHasilID,
		"catatan_admin": req.CatatanAdmin,
		"selesai_pada":  now,
	})

	database.DB.Create(&models.RiwayatSurat{
		PengajuanID: pengajuan.ID,
		Status:      "selesai",
		Catatan:     req.CatatanAdmin,
		Oleh:        &userID,
	})

	helpers.OK(c, "Pengajuan berhasil diselesaikan", nil)
}

func TolakSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengajuan models.PengajuanSurat
	if err := database.DB.First(&pengajuan, id).Error; err != nil {
		helpers.NotFound(c, "Pengajuan tidak ditemukan")
		return
	}

	if pengajuan.Status == "selesai" || pengajuan.Status == "ditolak" {
		helpers.BadRequest(c, "Pengajuan tidak dapat ditolak", nil)
		return
	}

	var req structs.TolakSuratRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	userID := c.GetUint("user_id")

	database.DB.Model(&pengajuan).Updates(map[string]interface{}{
		"status":        "ditolak",
		"catatan_admin": req.Alasan,
	})

	database.DB.Create(&models.RiwayatSurat{
		PengajuanID: pengajuan.ID,
		Status:      "ditolak",
		Catatan:     &req.Alasan,
		Oleh:        &userID,
	})

	helpers.OK(c, "Pengajuan berhasil ditolak", nil)
}
