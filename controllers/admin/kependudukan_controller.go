// Platform Desa — Admin Kependudukan Controller
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

func GetStatistikPendudukList(c *gin.Context) {
	tahun := c.Query("tahun")

	query := database.DB.Model(&models.StatistikPenduduk{})

	if tahun != "" {
		query = query.Where("tahun = ?", tahun)
	}

	var data []models.StatistikPenduduk
	query.Order("tahun DESC, bulan DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetStatistikPenduduk(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPenduduk
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func CreateStatistikPenduduk(c *gin.Context) {
	var req structs.StatistikPendudukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.StatistikPenduduk
	if err := database.DB.Where("tahun = ? AND bulan = ?", req.Tahun, req.Bulan).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Data bulan ini sudah ada", nil)
		return
	}

	data := models.StatistikPenduduk{
		Tahun:         req.Tahun,
		Bulan:         req.Bulan,
		TotalPenduduk: req.TotalPenduduk,
		LakiLaki:      req.LakiLaki,
		Perempuan:     req.Perempuan,
		TotalKK:       req.TotalKK,
		Kelahiran:     req.Kelahiran,
		Kematian:      req.Kematian,
		PindahMasuk:   req.PindahMasuk,
		PindahKeluar:  req.PindahKeluar,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah data")
		return
	}

	go helpers.RevalidatePath("/kependudukan")

	helpers.Created(c, "Data berhasil ditambahkan", data)
}

func UpdateStatistikPenduduk(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPenduduk
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	var req structs.StatistikPendudukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.StatistikPenduduk
	if err := database.DB.Where("tahun = ? AND bulan = ? AND id != ?", req.Tahun, req.Bulan, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Data bulan ini sudah ada", nil)
		return
	}

	data.Tahun = req.Tahun
	data.Bulan = req.Bulan
	data.TotalPenduduk = req.TotalPenduduk
	data.LakiLaki = req.LakiLaki
	data.Perempuan = req.Perempuan
	data.TotalKK = req.TotalKK
	data.Kelahiran = req.Kelahiran
	data.Kematian = req.Kematian
	data.PindahMasuk = req.PindahMasuk
	data.PindahKeluar = req.PindahKeluar

	if err := database.DB.Save(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui data")
		return
	}

	go helpers.RevalidatePath("/kependudukan")

	helpers.OK(c, "Data berhasil diperbarui", data)
}

func DeleteStatistikPenduduk(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPenduduk
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	database.DB.Delete(&data)

	go helpers.RevalidatePath("/kependudukan")

	helpers.OK(c, "Data berhasil dihapus", nil)
}

func GetStatistikPendidikanList(c *gin.Context) {
	var data []models.StatistikPendidikan
	database.DB.Order("tahun DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetStatistikPendidikan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPendidikan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func CreateStatistikPendidikan(c *gin.Context) {
	var req structs.StatistikPendidikanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.StatistikPendidikan
	if err := database.DB.Where("tahun = ?", req.Tahun).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Data tahun ini sudah ada", nil)
		return
	}

	data := models.StatistikPendidikan{
		Tahun:        req.Tahun,
		TidakSekolah: req.TidakSekolah,
		SD:           req.SD,
		SMP:          req.SMP,
		SMA:          req.SMA,
		Diploma:      req.Diploma,
		Sarjana:      req.Sarjana,
		Pascasarjana: req.Pascasarjana,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah data")
		return
	}

	go helpers.RevalidatePath("/kependudukan")

	helpers.Created(c, "Data berhasil ditambahkan", data)
}

func UpdateStatistikPendidikan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPendidikan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	var req structs.StatistikPendidikanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.StatistikPendidikan
	if err := database.DB.Where("tahun = ? AND id != ?", req.Tahun, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Data tahun ini sudah ada", nil)
		return
	}

	data.Tahun = req.Tahun
	data.TidakSekolah = req.TidakSekolah
	data.SD = req.SD
	data.SMP = req.SMP
	data.SMA = req.SMA
	data.Diploma = req.Diploma
	data.Sarjana = req.Sarjana
	data.Pascasarjana = req.Pascasarjana

	if err := database.DB.Save(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui data")
		return
	}

	go helpers.RevalidatePath("/kependudukan")

	helpers.OK(c, "Data berhasil diperbarui", data)
}

func DeleteStatistikPendidikan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPendidikan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	database.DB.Delete(&data)

	go helpers.RevalidatePath("/kependudukan")

	helpers.OK(c, "Data berhasil dihapus", nil)
}

func GetStatistikPekerjaanList(c *gin.Context) {
	var data []models.StatistikPekerjaan
	database.DB.Order("tahun DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetStatistikPekerjaan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPekerjaan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func CreateStatistikPekerjaan(c *gin.Context) {
	var req structs.StatistikPekerjaanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.StatistikPekerjaan
	if err := database.DB.Where("tahun = ?", req.Tahun).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Data tahun ini sudah ada", nil)
		return
	}

	data := models.StatistikPekerjaan{
		Tahun:      req.Tahun,
		Petani:     req.Petani,
		Pedagang:   req.Pedagang,
		PNS:        req.PNS,
		Swasta:     req.Swasta,
		Wiraswasta: req.Wiraswasta,
		Buruh:      req.Buruh,
		Lainnya:    req.Lainnya,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah data")
		return
	}

	go helpers.RevalidatePath("/kependudukan")

	helpers.Created(c, "Data berhasil ditambahkan", data)
}

func UpdateStatistikPekerjaan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPekerjaan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	var req structs.StatistikPekerjaanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.StatistikPekerjaan
	if err := database.DB.Where("tahun = ? AND id != ?", req.Tahun, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Data tahun ini sudah ada", nil)
		return
	}

	data.Tahun = req.Tahun
	data.Petani = req.Petani
	data.Pedagang = req.Pedagang
	data.PNS = req.PNS
	data.Swasta = req.Swasta
	data.Wiraswasta = req.Wiraswasta
	data.Buruh = req.Buruh
	data.Lainnya = req.Lainnya

	if err := database.DB.Save(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui data")
		return
	}

	go helpers.RevalidatePath("/kependudukan")

	helpers.OK(c, "Data berhasil diperbarui", data)
}

func DeleteStatistikPekerjaan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.StatistikPekerjaan
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Data tidak ditemukan")
		return
	}

	database.DB.Delete(&data)

	go helpers.RevalidatePath("/kependudukan")

	helpers.OK(c, "Data berhasil dihapus", nil)
}
