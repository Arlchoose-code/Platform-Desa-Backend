// Platform Desa — Admin Pengumuman & Agenda Controller
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

func GetPengumumanList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	status := c.Query("status")
	isPenting := c.Query("is_penting")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Pengumuman{}).
		Preload("File").
		Preload("Penulis").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if isPenting != "" {
		query = query.Where("is_penting = ?", isPenting == "true")
	}

	var total int64
	query.Count(&total)

	var pengumuman []models.Pengumuman
	query.Scopes(helpers.Paginate(pg)).Order("is_penting DESC, created_at DESC").Find(&pengumuman)

	helpers.OKPaginated(c, "Berhasil", pengumuman, pg.Meta(total))
}

func GetPengumuman(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengumuman models.Pengumuman
	if err := database.DB.Preload("File").Preload("Penulis").
		Where("id = ? AND is_deleted = false", id).
		First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", pengumuman)
}

func CreatePengumuman(c *gin.Context) {
	var req structs.PengumumanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	userID := c.GetUint("user_id")

	pengumuman := models.Pengumuman{
		Judul:          req.Judul,
		Slug:           helpers.UniqueSlug(helpers.GenerateSlug(req.Judul)),
		Isi:            req.Isi,
		FileID:         req.FileID,
		PenulisID:      &userID,
		IsPenting:      req.IsPenting,
		Status:         req.Status,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
	}

	if req.Status == "published" {
		now := time.Now()
		pengumuman.PublishedAt = &now
	}

	if err := database.DB.Create(&pengumuman).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah pengumuman")
		return
	}

	helpers.Log(c, "create", "pengumuman", "Menambah pengumuman: "+req.Judul)

	if req.Status == "published" {
		go helpers.RevalidatePaths("/pengumuman", "/pengumuman/"+pengumuman.Slug)
	}

	database.DB.Preload("File").Preload("Penulis").First(&pengumuman, pengumuman.ID)
	helpers.Created(c, "Pengumuman berhasil ditambahkan", pengumuman)
}

func UpdatePengumuman(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengumuman models.Pengumuman
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan")
		return
	}

	var req structs.PengumumanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Judul != pengumuman.Judul {
		pengumuman.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Judul))
	}

	pengumuman.Judul = req.Judul
	pengumuman.Isi = req.Isi
	pengumuman.FileID = req.FileID
	pengumuman.IsPenting = req.IsPenting
	pengumuman.Status = req.Status
	pengumuman.TanggalMulai = req.TanggalMulai
	pengumuman.TanggalSelesai = req.TanggalSelesai

	if req.Status == "published" && pengumuman.PublishedAt == nil {
		now := time.Now()
		pengumuman.PublishedAt = &now
	}

	if err := database.DB.Save(&pengumuman).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui pengumuman")
		return
	}

	helpers.Log(c, "update", "pengumuman", "Memperbarui pengumuman: "+req.Judul)

	go helpers.RevalidatePaths("/pengumuman", "/pengumuman/"+pengumuman.Slug)

	database.DB.Preload("File").Preload("Penulis").First(&pengumuman, pengumuman.ID)
	helpers.OK(c, "Pengumuman berhasil diperbarui", pengumuman)
}

func DeletePengumuman(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengumuman models.Pengumuman
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan")
		return
	}

	database.DB.Model(&pengumuman).Update("is_deleted", true)

	helpers.Log(c, "delete", "pengumuman", "Menghapus pengumuman: "+pengumuman.Judul)

	go helpers.RevalidatePath("/pengumuman")

	helpers.OK(c, "Pengumuman berhasil dihapus", nil)
}

func ForceDeletePengumuman(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengumuman models.Pengumuman
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan di trash")
		return
	}

	if pengumuman.FileID != nil {
		var media models.Media
		if err := database.DB.First(&media, pengumuman.FileID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&pengumuman)

	helpers.Log(c, "force_delete", "pengumuman", "Menghapus permanen pengumuman: "+pengumuman.Judul)

	helpers.OK(c, "Pengumuman berhasil dihapus permanen", nil)
}

func RestorePengumuman(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengumuman models.Pengumuman
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan di trash")
		return
	}

	database.DB.Model(&pengumuman).Update("is_deleted", false)

	helpers.Log(c, "restore", "pengumuman", "Memulihkan pengumuman: "+pengumuman.Judul)

	helpers.OK(c, "Pengumuman berhasil dipulihkan", nil)
}

func PublishPengumuman(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pengumuman models.Pengumuman
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&pengumuman).Error; err != nil {
		helpers.NotFound(c, "Pengumuman tidak ditemukan")
		return
	}

	if pengumuman.Status == "published" {
		database.DB.Model(&pengumuman).Update("status", "draft")
		helpers.Log(c, "unpublish", "pengumuman", "Menyembunyikan pengumuman: "+pengumuman.Judul)
		go helpers.RevalidatePaths("/pengumuman", "/pengumuman/"+pengumuman.Slug)
		helpers.OK(c, "Pengumuman berhasil di-unpublish", gin.H{"status": "draft"})
		return
	}

	now := time.Now()
	database.DB.Model(&pengumuman).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	})

	helpers.Log(c, "publish", "pengumuman", "Mempublikasikan pengumuman: "+pengumuman.Judul)

	go helpers.RevalidatePaths("/pengumuman", "/pengumuman/"+pengumuman.Slug)
	helpers.OK(c, "Pengumuman berhasil dipublikasikan", gin.H{"status": "published"})
}

func GetAgendaList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	status := c.Query("status")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Agenda{}).
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var agenda models.Agenda
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", agenda)
}

func CreateAgenda(c *gin.Context) {
	var req structs.AgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	userID := c.GetUint("user_id")

	agenda := models.Agenda{
		Judul:          req.Judul,
		Slug:           helpers.UniqueSlug(helpers.GenerateSlug(req.Judul)),
		Deskripsi:      req.Deskripsi,
		Lokasi:         req.Lokasi,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
		Penyelenggara:  req.Penyelenggara,
		Status:         req.Status,
		IsPublished:    req.IsPublished,
		CreatedBy:      &userID,
	}

	if err := database.DB.Create(&agenda).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah agenda")
		return
	}

	helpers.Log(c, "create", "pengumuman", "Menambah agenda: "+req.Judul)

	if req.IsPublished {
		go helpers.RevalidatePaths("/agenda", "/agenda/"+agenda.Slug)
	}

	helpers.Created(c, "Agenda berhasil ditambahkan", agenda)
}

func UpdateAgenda(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var agenda models.Agenda
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan")
		return
	}

	var req structs.AgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Judul != agenda.Judul {
		agenda.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Judul))
	}

	agenda.Judul = req.Judul
	agenda.Deskripsi = req.Deskripsi
	agenda.Lokasi = req.Lokasi
	agenda.TanggalMulai = req.TanggalMulai
	agenda.TanggalSelesai = req.TanggalSelesai
	agenda.Penyelenggara = req.Penyelenggara
	agenda.Status = req.Status
	agenda.IsPublished = req.IsPublished

	if err := database.DB.Save(&agenda).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui agenda")
		return
	}

	helpers.Log(c, "update", "pengumuman", "Memperbarui agenda: "+req.Judul)

	go helpers.RevalidatePaths("/agenda", "/agenda/"+agenda.Slug)

	helpers.OK(c, "Agenda berhasil diperbarui", agenda)
}

func DeleteAgenda(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var agenda models.Agenda
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan")
		return
	}

	database.DB.Model(&agenda).Update("is_deleted", true)

	helpers.Log(c, "delete", "pengumuman", "Menghapus agenda: "+agenda.Judul)

	go helpers.RevalidatePath("/agenda")

	helpers.OK(c, "Agenda berhasil dihapus", nil)
}

func ForceDeleteAgenda(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var agenda models.Agenda
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan di trash")
		return
	}

	database.DB.Delete(&agenda)

	helpers.Log(c, "force_delete", "pengumuman", "Menghapus permanen agenda: "+agenda.Judul)

	helpers.OK(c, "Agenda berhasil dihapus permanen", nil)
}

func RestoreAgenda(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var agenda models.Agenda
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan di trash")
		return
	}

	database.DB.Model(&agenda).Update("is_deleted", false)

	helpers.Log(c, "restore", "pengumuman", "Memulihkan agenda: "+agenda.Judul)

	helpers.OK(c, "Agenda berhasil dipulihkan", nil)
}

func PublishAgenda(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var agenda models.Agenda
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&agenda).Error; err != nil {
		helpers.NotFound(c, "Agenda tidak ditemukan")
		return
	}

	newStatus := !agenda.IsPublished
	database.DB.Model(&agenda).Update("is_published", newStatus)

	go helpers.RevalidatePaths("/agenda", "/agenda/"+agenda.Slug)

	msg := "Agenda berhasil dipublikasikan"
	if !newStatus {
		msg = "Agenda berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "pengumuman", msg+": "+agenda.Judul)

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}
