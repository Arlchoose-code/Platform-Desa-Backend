// Platform Desa — Admin Regulasi Controller
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

func GetKategoriRegulasiList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")

	query := database.DB.Model(&models.KategoriRegulasi{}).Where("is_deleted = ?", showTrash)
	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}

	var kategori []models.KategoriRegulasi
	query.Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetKategoriRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriRegulasi
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", kategori)
}

func CreateKategoriRegulasi(c *gin.Context) {
	var req structs.KategoriRegulasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriRegulasi
	if err := database.DB.Where("nama = ? AND is_deleted = false", req.Nama).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	kategori := models.KategoriRegulasi{
		Nama: req.Nama,
		Slug: helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
	}

	if err := database.DB.Create(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah kategori")
		return
	}

	helpers.Log(c, "create", "regulasi", "Menambah kategori regulasi: "+req.Nama)

	helpers.Created(c, "Kategori berhasil ditambahkan", kategori)
}

func UpdateKategoriRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriRegulasi
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var req structs.KategoriRegulasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriRegulasi
	if err := database.DB.Where("nama = ? AND id != ? AND is_deleted = false", req.Nama, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	if req.Nama != kategori.Nama {
		kategori.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}
	kategori.Nama = req.Nama

	if err := database.DB.Save(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui kategori")
		return
	}

	helpers.Log(c, "update", "regulasi", "Memperbarui kategori regulasi: "+req.Nama)

	helpers.OK(c, "Kategori berhasil diperbarui", kategori)
}

func DeleteKategoriRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriRegulasi
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var count int64
	database.DB.Model(&models.Regulasi{}).Where("kategori_id = ? AND is_deleted = false", id).Count(&count)
	if count > 0 {
		helpers.BadRequest(c, "Kategori masih digunakan oleh regulasi aktif", nil)
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", true)

	helpers.Log(c, "delete", "regulasi", "Menghapus kategori regulasi: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus", nil)
}

func ForceDeleteKategoriRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriRegulasi
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&models.Regulasi{}).Where("kategori_id = ?", id).Updates(map[string]interface{}{
		"kategori_id":   nil,
		"kategori_nama": nil,
	})

	database.DB.Delete(&kategori)

	helpers.Log(c, "force_delete", "regulasi", "Menghapus permanen kategori regulasi: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus permanen", nil)
}

func RestoreKategoriRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriRegulasi
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", false)

	helpers.Log(c, "restore", "regulasi", "Memulihkan kategori regulasi: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dipulihkan", nil)
}

func GetRegulasiList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriID := c.Query("kategori_id")
	tahun := c.Query("tahun")
	isPublished := c.Query("is_published")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Regulasi{}).
		Preload("Kategori").
		Preload("File").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("judul LIKE ? OR nomor LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}
	if tahun != "" {
		query = query.Where("tahun = ?", tahun)
	}
	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var total int64
	query.Count(&total)

	var regulasi []models.Regulasi
	query.Scopes(helpers.Paginate(pg)).Order("tahun DESC, tanggal_terbit DESC").Find(&regulasi)

	helpers.OKPaginated(c, "Berhasil", regulasi, pg.Meta(total))
}

func GetRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var regulasi models.Regulasi
	if err := database.DB.Preload("Kategori").Preload("File").
		Where("id = ? AND is_deleted = false", id).
		First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", regulasi)
}

func CreateRegulasi(c *gin.Context) {
	var req structs.RegulasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	regulasi := models.Regulasi{
		Nomor:         req.Nomor,
		Judul:         req.Judul,
		Slug:          helpers.UniqueSlug(helpers.GenerateSlug(req.Judul)),
		Tentang:       req.Tentang,
		Tahun:         req.Tahun,
		TanggalTerbit: req.TanggalTerbit,
		FileID:        req.FileID,
		IsPublished:   req.IsPublished,
	}

	if req.KategoriID != nil {
		var kategori models.KategoriRegulasi
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			regulasi.KategoriID = req.KategoriID
			regulasi.KategoriNama = &kategori.Nama
		}
	}

	if err := database.DB.Create(&regulasi).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah regulasi")
		return
	}

	helpers.Log(c, "create", "regulasi", "Menambah regulasi: "+req.Judul)

	if req.IsPublished {
		go helpers.RevalidatePaths("/regulasi", "/regulasi/"+regulasi.Slug)
	}

	database.DB.Preload("Kategori").Preload("File").First(&regulasi, regulasi.ID)
	helpers.Created(c, "Regulasi berhasil ditambahkan", regulasi)
}

func UpdateRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var regulasi models.Regulasi
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan")
		return
	}

	var req structs.RegulasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Judul != regulasi.Judul {
		regulasi.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Judul))
	}

	regulasi.Nomor = req.Nomor
	regulasi.Judul = req.Judul
	regulasi.Tentang = req.Tentang
	regulasi.Tahun = req.Tahun
	regulasi.TanggalTerbit = req.TanggalTerbit
	regulasi.FileID = req.FileID
	regulasi.IsPublished = req.IsPublished

	if req.KategoriID != nil {
		var kategori models.KategoriRegulasi
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			regulasi.KategoriID = req.KategoriID
			regulasi.KategoriNama = &kategori.Nama
		}
	} else {
		regulasi.KategoriID = nil
		regulasi.KategoriNama = nil
	}

	if err := database.DB.Save(&regulasi).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui regulasi")
		return
	}

	helpers.Log(c, "update", "regulasi", "Memperbarui regulasi: "+req.Judul)

	go helpers.RevalidatePaths("/regulasi", "/regulasi/"+regulasi.Slug)

	database.DB.Preload("Kategori").Preload("File").First(&regulasi, regulasi.ID)
	helpers.OK(c, "Regulasi berhasil diperbarui", regulasi)
}

func DeleteRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var regulasi models.Regulasi
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan")
		return
	}

	database.DB.Model(&regulasi).Update("is_deleted", true)

	helpers.Log(c, "delete", "regulasi", "Menghapus regulasi: "+regulasi.Judul)

	go helpers.RevalidatePath("/regulasi")

	helpers.OK(c, "Regulasi berhasil dihapus", nil)
}

func ForceDeleteRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var regulasi models.Regulasi
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan di trash")
		return
	}

	if regulasi.FileID != nil {
		var file models.Media
		if err := database.DB.First(&file, regulasi.FileID).Error; err == nil {
			go helpers.DeleteFile(file.Path)
			database.DB.Delete(&file)
		}
	}

	database.DB.Delete(&regulasi)

	helpers.Log(c, "force_delete", "regulasi", "Menghapus permanen regulasi: "+regulasi.Judul)

	helpers.OK(c, "Regulasi berhasil dihapus permanen", nil)
}

func RestoreRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var regulasi models.Regulasi
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan di trash")
		return
	}

	database.DB.Model(&regulasi).Update("is_deleted", false)

	helpers.Log(c, "restore", "regulasi", "Memulihkan regulasi: "+regulasi.Judul)

	helpers.OK(c, "Regulasi berhasil dipulihkan", nil)
}

func PublishRegulasi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var regulasi models.Regulasi
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&regulasi).Error; err != nil {
		helpers.NotFound(c, "Regulasi tidak ditemukan")
		return
	}

	newStatus := !regulasi.IsPublished
	database.DB.Model(&regulasi).Update("is_published", newStatus)

	go helpers.RevalidatePaths("/regulasi", "/regulasi/"+regulasi.Slug)

	msg := "Regulasi berhasil dipublikasikan"
	if !newStatus {
		msg = "Regulasi berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "regulasi", msg+": "+regulasi.Judul)

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}
