// Platform Desa — Admin APBDes Controller
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

func GetAPBDesList(c *gin.Context) {
	isPublished := c.Query("is_published")

	query := database.DB.Model(&models.APBDes{}).Preload("Dokumen")

	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var data []models.APBDes
	query.Order("tahun DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetAPBDes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.APBDes
	if err := database.DB.Preload("Dokumen").Preload("Detail").
		First(&data, id).Error; err != nil {
		helpers.NotFound(c, "APBDes tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func CreateAPBDes(c *gin.Context) {
	var req structs.APBDesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.APBDes
	if err := database.DB.Where("tahun = ?", req.Tahun).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "APBDes tahun ini sudah ada", nil)
		return
	}

	data := models.APBDes{
		Tahun:       req.Tahun,
		DokumenID:   req.DokumenID,
		IsPublished: req.IsPublished,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah APBDes")
		return
	}

	helpers.Log(c, "create", "apbdes", "Menambah APBDes tahun "+strconv.Itoa(req.Tahun))

	if req.IsPublished {
		go helpers.RevalidatePath("/apbdes")
	}

	database.DB.Preload("Dokumen").First(&data, data.ID)
	helpers.Created(c, "APBDes berhasil ditambahkan", data)
}

func UpdateAPBDes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.APBDes
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "APBDes tidak ditemukan")
		return
	}

	var req structs.APBDesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.APBDes
	if err := database.DB.Where("tahun = ? AND id != ?", req.Tahun, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "APBDes tahun ini sudah ada", nil)
		return
	}

	data.Tahun = req.Tahun
	data.DokumenID = req.DokumenID
	data.IsPublished = req.IsPublished

	if err := database.DB.Save(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui APBDes")
		return
	}

	helpers.Log(c, "update", "apbdes", "Memperbarui APBDes tahun "+strconv.Itoa(req.Tahun))

	go helpers.RevalidatePath("/apbdes")

	database.DB.Preload("Dokumen").First(&data, data.ID)
	helpers.OK(c, "APBDes berhasil diperbarui", data)
}

func DeleteAPBDes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.APBDes
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "APBDes tidak ditemukan")
		return
	}

	database.DB.Where("apbdes_id = ?", data.ID).Delete(&models.APBDesDetail{})
	database.DB.Delete(&data)

	helpers.Log(c, "delete", "apbdes", "Menghapus APBDes tahun "+strconv.Itoa(data.Tahun))

	go helpers.RevalidatePath("/apbdes")

	helpers.OK(c, "APBDes berhasil dihapus", nil)
}

func PublishAPBDes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.APBDes
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "APBDes tidak ditemukan")
		return
	}

	newStatus := !data.IsPublished
	database.DB.Model(&data).Update("is_published", newStatus)

	msg := "APBDes berhasil dipublikasikan"
	if !newStatus {
		msg = "APBDes berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "apbdes", msg+" tahun "+strconv.Itoa(data.Tahun))

	go helpers.RevalidatePath("/apbdes")

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}

func GetAPBDesDetailList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	jenis := c.Query("jenis")

	query := database.DB.Where("apbdes_id = ?", id)
	if jenis != "" {
		query = query.Where("jenis = ?", jenis)
	}

	var detail []models.APBDesDetail
	query.Order("jenis ASC, urutan ASC").Find(&detail)

	helpers.OK(c, "Berhasil", detail)
}

func CreateAPBDesDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var apbdes models.APBDes
	if err := database.DB.First(&apbdes, id).Error; err != nil {
		helpers.NotFound(c, "APBDes tidak ditemukan")
		return
	}

	var req structs.APBDesDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	detail := models.APBDesDetail{
		APBDesID:  apbdes.ID,
		Jenis:     req.Jenis,
		Kategori:  req.Kategori,
		Uraian:    req.Uraian,
		Anggaran:  req.Anggaran,
		Realisasi: req.Realisasi,
		Urutan:    req.Urutan,
	}

	if err := database.DB.Create(&detail).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah detail")
		return
	}

	helpers.Log(c, "create", "apbdes", "Menambah detail "+req.Jenis+" APBDes tahun "+strconv.Itoa(apbdes.Tahun))

	recalcAPBDes(apbdes.ID)
	go helpers.RevalidatePath("/apbdes")

	helpers.Created(c, "Detail berhasil ditambahkan", detail)
}

func UpdateAPBDesDetail(c *gin.Context) {
	apbdesID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	detailID, err := strconv.Atoi(c.Param("detail_id"))
	if err != nil {
		helpers.BadRequest(c, "ID detail tidak valid", nil)
		return
	}

	var detail models.APBDesDetail
	if err := database.DB.Where("id = ? AND apbdes_id = ?", detailID, apbdesID).First(&detail).Error; err != nil {
		helpers.NotFound(c, "Detail tidak ditemukan")
		return
	}

	var req structs.APBDesDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	detail.Jenis = req.Jenis
	detail.Kategori = req.Kategori
	detail.Uraian = req.Uraian
	detail.Anggaran = req.Anggaran
	detail.Realisasi = req.Realisasi
	detail.Urutan = req.Urutan

	if err := database.DB.Save(&detail).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui detail")
		return
	}

	helpers.Log(c, "update", "apbdes", "Memperbarui detail "+req.Jenis)

	recalcAPBDes(uint(apbdesID))
	go helpers.RevalidatePath("/apbdes")

	helpers.OK(c, "Detail berhasil diperbarui", detail)
}

func DeleteAPBDesDetail(c *gin.Context) {
	apbdesID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	detailID, err := strconv.Atoi(c.Param("detail_id"))
	if err != nil {
		helpers.BadRequest(c, "ID detail tidak valid", nil)
		return
	}

	var detail models.APBDesDetail
	if err := database.DB.Where("id = ? AND apbdes_id = ?", detailID, apbdesID).First(&detail).Error; err != nil {
		helpers.NotFound(c, "Detail tidak ditemukan")
		return
	}

	database.DB.Delete(&detail)

	helpers.Log(c, "delete", "apbdes", "Menghapus detail "+detail.Jenis+": "+detail.Uraian)

	recalcAPBDes(uint(apbdesID))
	go helpers.RevalidatePath("/apbdes")

	helpers.OK(c, "Detail berhasil dihapus", nil)
}

func recalcAPBDes(apbdesID uint) {
	type Result struct {
		Jenis string `gorm:"column:jenis"`
		Total int64  `gorm:"column:total"`
	}

	var results []Result
	database.DB.Model(&models.APBDesDetail{}).
		Select("jenis, SUM(anggaran) as total").
		Where("apbdes_id = ?", apbdesID).
		Group("jenis").
		Scan(&results)

	var pendapatan, belanja, pembiayaan int64
	for _, r := range results {
		switch r.Jenis {
		case "pendapatan":
			pendapatan = r.Total
		case "belanja":
			belanja = r.Total
		case "pembiayaan":
			pembiayaan = r.Total
		}
	}

	database.DB.Model(&models.APBDes{}).Where("id = ?", apbdesID).
		Select("total_pendapatan", "total_belanja", "total_pembiayaan").
		Updates(models.APBDes{
			TotalPendapatan: pendapatan,
			TotalBelanja:    belanja,
			TotalPembiayaan: pembiayaan,
		})
}
