// Platform Desa — Admin Pengaduan Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"strconv"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
)

func GetPengaduanList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	status := c.Query("status")
	kategori := c.Query("kategori")

	query := database.DB.Model(&models.Pengaduan{}).Preload("Foto")

	if search != "" {
		query = query.Where("judul LIKE ? OR nama_pelapor LIKE ? OR nomor_tiket LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	var total int64
	query.Count(&total)

	var data []models.Pengaduan
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&data)

	helpers.OKPaginated(c, "Berhasil", data, pg.Meta(total))
}

func GetPengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.Pengaduan
	if err := database.DB.Preload("Foto").First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Pengaduan tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func BuatPengaduan(c *gin.Context) {
	var req structs.PengaduanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	namaPelapor := req.NamaPelapor
	if req.IsAnonim {
		namaPelapor = "Anonim"
	}

	data := models.Pengaduan{
		NomorTiket:    helpers.GenerateNomorTiket(),
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

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal membuat pengaduan")
		return
	}

	helpers.Created(c, "Pengaduan berhasil dibuat", gin.H{
		"nomor_tiket": data.NomorTiket,
		"id":          data.ID,
	})
}

func VerifikasiPengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.Pengaduan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Pengaduan tidak ditemukan")
		return
	}

	if data.Status != "masuk" {
		helpers.BadRequest(c, "Pengaduan tidak dalam status masuk", nil)
		return
	}

	userID := c.GetUint("user_id")
	database.DB.Model(&data).Updates(map[string]interface{}{
		"status":         "diverifikasi",
		"ditangani_oleh": userID,
	})

	helpers.OK(c, "Pengaduan berhasil diverifikasi", nil)
}

func ProsesPengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.Pengaduan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Pengaduan tidak ditemukan")
		return
	}

	if data.Status != "diverifikasi" {
		helpers.BadRequest(c, "Pengaduan tidak dalam status diverifikasi", nil)
		return
	}

	userID := c.GetUint("user_id")
	database.DB.Model(&data).Updates(map[string]interface{}{
		"status":         "diproses",
		"ditangani_oleh": userID,
	})

	helpers.OK(c, "Pengaduan sedang diproses", nil)
}

func SelesaikanPengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.Pengaduan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Pengaduan tidak ditemukan")
		return
	}

	if data.Status != "diproses" {
		helpers.BadRequest(c, "Pengaduan tidak dalam status diproses", nil)
		return
	}

	var req structs.ResponPengaduanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	database.DB.Model(&data).Updates(map[string]interface{}{
		"status":       "selesai",
		"respon_admin": req.Respon,
	})

	helpers.OK(c, "Pengaduan berhasil diselesaikan", nil)
}

func TolakPengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.Pengaduan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Pengaduan tidak ditemukan")
		return
	}

	if data.Status == "selesai" || data.Status == "ditolak" {
		helpers.BadRequest(c, "Pengaduan tidak dapat ditolak", nil)
		return
	}

	var req structs.TolakPengaduanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	database.DB.Model(&data).Updates(map[string]interface{}{
		"status":       "ditolak",
		"respon_admin": req.Alasan,
	})

	helpers.OK(c, "Pengaduan berhasil ditolak", nil)
}

func DeletePengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.Pengaduan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Pengaduan tidak ditemukan")
		return
	}

	if data.FotoID != nil {
		var foto models.Media
		if err := database.DB.First(&foto, data.FotoID).Error; err == nil {
			go helpers.DeleteFile(foto.Path)
			database.DB.Delete(&foto)
		}
	}

	database.DB.Delete(&data)
	helpers.OK(c, "Pengaduan berhasil dihapus", nil)
}
