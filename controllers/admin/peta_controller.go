// Platform Desa — Admin Peta Controller
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

func GetPetaFasilitasList(c *gin.Context) {
	search := c.Query("search")
	kategori := c.Query("kategori")
	isPublished := c.Query("is_published")

	query := database.DB.Model(&models.PetaFasilitas{}).Preload("Foto")

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var data []models.PetaFasilitas
	query.Order("kategori ASC, nama ASC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetPetaFasilitas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.PetaFasilitas
	if err := database.DB.Preload("Foto").First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Fasilitas tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func CreatePetaFasilitas(c *gin.Context) {
	var req structs.PetaFasilitasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	data := models.PetaFasilitas{
		Nama:        req.Nama,
		Kategori:    req.Kategori,
		Alamat:      req.Alamat,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Deskripsi:   req.Deskripsi,
		FotoID:      req.FotoID,
		IsPublished: req.IsPublished,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah fasilitas")
		return
	}

	if req.IsPublished {
		go helpers.RevalidatePath("/peta")
	}

	database.DB.Preload("Foto").First(&data, data.ID)
	helpers.Created(c, "Fasilitas berhasil ditambahkan", data)
}

func UpdatePetaFasilitas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.PetaFasilitas
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Fasilitas tidak ditemukan")
		return
	}

	var req structs.PetaFasilitasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	data.Nama = req.Nama
	data.Kategori = req.Kategori
	data.Alamat = req.Alamat
	data.Latitude = req.Latitude
	data.Longitude = req.Longitude
	data.Deskripsi = req.Deskripsi
	data.FotoID = req.FotoID
	data.IsPublished = req.IsPublished

	if err := database.DB.Save(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui fasilitas")
		return
	}

	go helpers.RevalidatePath("/peta")

	database.DB.Preload("Foto").First(&data, data.ID)
	helpers.OK(c, "Fasilitas berhasil diperbarui", data)
}

func DeletePetaFasilitas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.PetaFasilitas
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Fasilitas tidak ditemukan")
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

	go helpers.RevalidatePath("/peta")

	helpers.OK(c, "Fasilitas berhasil dihapus", nil)
}

func PublishPetaFasilitas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.PetaFasilitas
	if err := database.DB.First(&data, id).Error; err != nil {
		helpers.NotFound(c, "Fasilitas tidak ditemukan")
		return
	}

	newStatus := !data.IsPublished
	database.DB.Model(&data).Update("is_published", newStatus)

	go helpers.RevalidatePath("/peta")

	msg := "Fasilitas berhasil dipublikasikan"
	if !newStatus {
		msg = "Fasilitas berhasil disembunyikan"
	}
	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}
