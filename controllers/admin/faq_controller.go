// Platform Desa — Admin FAQ Controller
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

func GetFAQList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")
	kategori := c.Query("kategori")
	isPublished := c.Query("is_published")

	query := database.DB.Model(&models.FAQ{}).Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("pertanyaan LIKE ? OR jawaban LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var data []models.FAQ
	query.Order("urutan ASC, created_at ASC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.FAQ
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&data).Error; err != nil {
		helpers.NotFound(c, "FAQ tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", data)
}

func CreateFAQ(c *gin.Context) {
	var req structs.FAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	data := models.FAQ{
		Pertanyaan:  req.Pertanyaan,
		Jawaban:     req.Jawaban,
		Kategori:    req.Kategori,
		Urutan:      req.Urutan,
		IsPublished: req.IsPublished,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah FAQ")
		return
	}

	helpers.Log(c, "create", "faq", "Menambah FAQ: "+req.Pertanyaan)

	if req.IsPublished {
		go helpers.RevalidatePath("/faq")
	}

	helpers.Created(c, "FAQ berhasil ditambahkan", data)
}

func UpdateFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.FAQ
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&data).Error; err != nil {
		helpers.NotFound(c, "FAQ tidak ditemukan")
		return
	}

	var req structs.FAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	data.Pertanyaan = req.Pertanyaan
	data.Jawaban = req.Jawaban
	data.Kategori = req.Kategori
	data.Urutan = req.Urutan
	data.IsPublished = req.IsPublished

	if err := database.DB.Save(&data).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui FAQ")
		return
	}

	helpers.Log(c, "update", "faq", "Memperbarui FAQ: "+req.Pertanyaan)

	go helpers.RevalidatePath("/faq")

	helpers.OK(c, "FAQ berhasil diperbarui", data)
}

func DeleteFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.FAQ
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&data).Error; err != nil {
		helpers.NotFound(c, "FAQ tidak ditemukan")
		return
	}

	database.DB.Model(&data).Update("is_deleted", true)

	helpers.Log(c, "delete", "faq", "Menghapus FAQ: "+data.Pertanyaan)

	go helpers.RevalidatePath("/faq")

	helpers.OK(c, "FAQ berhasil dihapus", nil)
}

func ForceDeleteFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.FAQ
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&data).Error; err != nil {
		helpers.NotFound(c, "FAQ tidak ditemukan di trash")
		return
	}

	database.DB.Delete(&data)

	helpers.Log(c, "force_delete", "faq", "Menghapus permanen FAQ: "+data.Pertanyaan)

	helpers.OK(c, "FAQ berhasil dihapus permanen", nil)
}

func RestoreFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.FAQ
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&data).Error; err != nil {
		helpers.NotFound(c, "FAQ tidak ditemukan di trash")
		return
	}

	database.DB.Model(&data).Update("is_deleted", false)

	helpers.Log(c, "restore", "faq", "Memulihkan FAQ: "+data.Pertanyaan)

	helpers.OK(c, "FAQ berhasil dipulihkan", nil)
}

func PublishFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var data models.FAQ
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&data).Error; err != nil {
		helpers.NotFound(c, "FAQ tidak ditemukan")
		return
	}

	newStatus := !data.IsPublished
	database.DB.Model(&data).Update("is_published", newStatus)

	go helpers.RevalidatePath("/faq")

	msg := "FAQ berhasil dipublikasikan"
	if !newStatus {
		msg = "FAQ berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "faq", msg+": "+data.Pertanyaan)

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}

func UrutanFAQ(c *gin.Context) {
	type UrutanItem struct {
		ID     uint `json:"id" binding:"required"`
		Urutan uint `json:"urutan"`
	}
	type UrutanRequest struct {
		Items []UrutanItem `json:"items" binding:"required"`
	}

	var req UrutanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	for _, item := range req.Items {
		database.DB.Model(&models.FAQ{}).Where("id = ?", item.ID).Update("urutan", item.Urutan)
	}

	go helpers.RevalidatePath("/faq")

	helpers.OK(c, "Urutan berhasil diperbarui", nil)
}
